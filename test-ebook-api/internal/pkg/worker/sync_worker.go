package worker

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg/queue"
	"test-ebook-api/internal/pkg/storage"
	"test-ebook-api/internal/repository"
	"time"
)

// SyncWorker 文件同步 Worker
// 从 uploadQueue 消费任务，将文件从 staging 同步到最终存储
type SyncWorker struct {
	uploadQueue  *queue.Queue
	ocrQueue     *queue.Queue
	staging      *storage.StagingStorage
	remote       storage.Storage
	repo         *repository.StandardRepository
	uploadRepo   *repository.UploadTaskRepository
	concurrency  int
	retryPolicy  RetryPolicy
	storageType  string
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// SyncWorkerConfig SyncWorker 配置
type SyncWorkerConfig struct {
	UploadQueue *queue.Queue
	OCRQueue    *queue.Queue
	Staging     *storage.StagingStorage
	Remote      storage.Storage
	Repo        *repository.StandardRepository
	UploadRepo  *repository.UploadTaskRepository
	Concurrency int
	RetryPolicy RetryPolicy
	StorageType string
}

// NewSyncWorker 创建同步 Worker
func NewSyncWorker(cfg SyncWorkerConfig) *SyncWorker {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 3
	}
	return &SyncWorker{
		uploadQueue: cfg.UploadQueue,
		ocrQueue:    cfg.OCRQueue,
		staging:     cfg.Staging,
		remote:      cfg.Remote,
		repo:        cfg.Repo,
		uploadRepo:  cfg.UploadRepo,
		concurrency: cfg.Concurrency,
		retryPolicy: cfg.RetryPolicy,
		storageType: cfg.StorageType,
	}
}

// Start 启动 Worker，先恢复未完成任务，再开始消费队列
func (w *SyncWorker) Start(ctx context.Context) {
	ctx, w.cancel = context.WithCancel(ctx)

	// 恢复未完成的任务
	w.recoverPendingTasks()

	// 启动并发消费者
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.consume(ctx, i)
	}
	log.Printf("[SyncWorker] 已启动 %d 个消费者", w.concurrency)
}

// Stop 优雅停止 Worker
func (w *SyncWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
	log.Println("[SyncWorker] 所有消费者已停止")
}

// recoverPendingTasks 进程启动时恢复未完成的同步任务
func (w *SyncWorker) recoverPendingTasks() {
	// 恢复 pending/uploading 状态的任务
	tasks, err := w.uploadRepo.GetPendingTasks()
	if err != nil {
		log.Printf("[SyncWorker] 恢复待处理任务失败: %v", err)
		return
	}

	// 恢复可重试的失败任务
	retryable, err := w.uploadRepo.GetRetryableTasks()
	if err != nil {
		log.Printf("[SyncWorker] 恢复可重试任务失败: %v", err)
	} else {
		tasks = append(tasks, retryable...)
	}

	if len(tasks) == 0 {
		log.Println("[SyncWorker] 无需恢复的任务")
		return
	}

	recovered := 0
	for _, t := range tasks {
		// 检查暂存文件是否还存在
		if !w.staging.Exists(t.LocalPath) {
			log.Printf("[SyncWorker] 暂存文件已丢失，标记任务 %d 为失败: %s", t.ID, t.LocalPath)
			t.Status = "failed"
			t.Error = "暂存文件已丢失，无法恢复"
			w.uploadRepo.Update(&t)

			// 更新关联文档状态
			if file, err := w.repo.FindFileByID(t.DocumentID); err == nil {
				file.SyncStatus = "sync_failed"
				w.repo.UpdateFile(file)
			}
			continue
		}

		task := queue.Task{
			Type:       queue.TaskUploadSync,
			DocumentID: t.DocumentID,
			UploadID:   t.ID,
			RetryCount: t.RetryCount,
		}
		if err := w.uploadQueue.Push(task); err != nil {
			log.Printf("[SyncWorker] 恢复任务 %d 推入队列失败: %v", t.ID, err)
			continue
		}
		recovered++
	}
	log.Printf("[SyncWorker] 已恢复 %d 个任务到队列", recovered)
}

// consume 单个消费者协程
func (w *SyncWorker) consume(ctx context.Context, workerID int) {
	defer w.wg.Done()
	log.Printf("[SyncWorker-%d] 消费者启动", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[SyncWorker-%d] 消费者收到停止信号", workerID)
			return
		case task, ok := <-w.uploadQueue.Pop():
			if !ok {
				log.Printf("[SyncWorker-%d] 队列已关闭，消费者退出", workerID)
				return
			}
			if task.Type == queue.TaskUploadSync {
				w.processUpload(workerID, task)
			}
		}
	}
}

// processUpload 处理单个同步任务
func (w *SyncWorker) processUpload(workerID int, task queue.Task) {
	uploadTask, err := w.uploadRepo.GetByID(task.UploadID)
	if err != nil {
		log.Printf("[SyncWorker-%d] 获取上传任务 %d 失败: %v", workerID, task.UploadID, err)
		return
	}

	// 更新状态为 uploading
	uploadTask.Status = "uploading"
	w.uploadRepo.Update(uploadTask)

	// 更新文档同步状态
	file, err := w.repo.FindFileByID(uploadTask.DocumentID)
	if err != nil {
		log.Printf("[SyncWorker-%d] 获取文档 %d 失败: %v", workerID, uploadTask.DocumentID, err)
		w.handleFailure(workerID, uploadTask, task, "关联文档不存在")
		return
	}
	file.SyncStatus = "syncing"
	w.repo.UpdateFile(file)

	log.Printf("[SyncWorker-%d] 开始同步文件: 任务=%d, 文档=%d, 暂存=%s",
		workerID, uploadTask.ID, uploadTask.DocumentID, uploadTask.LocalPath)

	// 从 staging 读取文件
	reader, err := w.staging.Get(uploadTask.LocalPath)
	if err != nil {
		w.handleFailure(workerID, uploadTask, task, "读取暂存文件失败: "+err.Error())
		return
	}
	defer reader.Close()

	// 同步到最终存储
	remotePath, err := w.syncToFinalStorage(uploadTask, reader)
	if err != nil {
		w.handleFailure(workerID, uploadTask, task, "同步到最终存储失败: "+err.Error())
		return
	}

	// 同步成功
	uploadTask.Status = "synced"
	uploadTask.RemotePath = remotePath
	uploadTask.Error = ""
	w.uploadRepo.Update(uploadTask)

	// 更新文档的 FilePath 和同步状态
	file.FilePath = remotePath
	file.SyncStatus = "synced"
	w.repo.UpdateFile(file)

	// 删除 staging 文件
	if err := w.staging.Remove(uploadTask.LocalPath); err != nil {
		log.Printf("[SyncWorker-%d] 删除暂存文件失败（不影响同步结果）: %v", workerID, err)
	}

	log.Printf("[SyncWorker-%d] 同步成功: 任务=%d, 文档=%d, 远程路径=%s",
		workerID, uploadTask.ID, uploadTask.DocumentID, remotePath)

	// 推送 OCR 任务（如果有 OCR 队列且文件需要 OCR）
	if w.ocrQueue != nil && file.Status == 0 {
		ocrTask := queue.Task{
			Type:       queue.TaskOCR,
			DocumentID: file.ID,
		}
		if err := w.ocrQueue.Push(ocrTask); err != nil {
			log.Printf("[SyncWorker-%d] 推送 OCR 任务失败: %v", workerID, err)
		}
	}
}

// syncToFinalStorage 同步文件到最终存储
func (w *SyncWorker) syncToFinalStorage(uploadTask *model.UploadTask, reader io.ReadCloser) (string, error) {
	fileName := filepath.Base(uploadTask.LocalPath)

	if w.storageType == "local" {
		// 本地存储：移动文件到 uploads 目录
		return w.moveToLocal(uploadTask.LocalPath, fileName)
	}

	// 云存储：上传到远端
	remotePath, err := w.remote.Save(fileName, reader)
	if err != nil {
		return "", err
	}
	return remotePath, nil
}

// moveToLocal 本地存储模式：从 staging 移动到 uploads 目录
func (w *SyncWorker) moveToLocal(stagingPath string, fileName string) (string, error) {
	src, err := os.Open(stagingPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	savePath, err := w.remote.Save(fileName, src)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

// handleFailure 处理同步失败
func (w *SyncWorker) handleFailure(workerID int, uploadTask *model.UploadTask, task queue.Task, errMsg string) {
	uploadTask.RetryCount++
	uploadTask.Error = errMsg

	if w.retryPolicy.ShouldRetry(uploadTask.RetryCount) {
		// 还可以重试
		nextRetry := w.retryPolicy.NextRetryTime(uploadTask.RetryCount)
		uploadTask.NextRetryAt = &nextRetry
		uploadTask.Status = "failed" // 暂标为 failed，等到时间后重新入队
		w.uploadRepo.Update(uploadTask)

		// 延迟后重新入队
		delay := w.retryPolicy.NextDelay(uploadTask.RetryCount - 1)
		log.Printf("[SyncWorker-%d] 任务 %d 将在 %v 后重试 (第 %d/%d 次)",
			workerID, uploadTask.ID, delay, uploadTask.RetryCount, uploadTask.MaxRetry)

		go func(t queue.Task, d time.Duration) {
			time.Sleep(d)
			t.RetryCount = uploadTask.RetryCount
			if err := w.uploadQueue.Push(t); err != nil {
				log.Printf("[SyncWorker] 重试任务 %d 推入队列失败: %v", t.UploadID, err)
			}
		}(task, delay)
	} else {
		// 已达最大重试次数
		uploadTask.Status = "failed"
		w.uploadRepo.Update(uploadTask)
		log.Printf("[SyncWorker-%d] 任务 %d 已达最大重试次数 (%d)，标记为最终失败",
			workerID, uploadTask.ID, uploadTask.MaxRetry)
	}

	// 更新文档同步状态
	if file, err := w.repo.FindFileByID(uploadTask.DocumentID); err == nil {
		file.SyncStatus = "sync_failed"
		w.repo.UpdateFile(file)
	}
}
