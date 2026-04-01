package worker

import (
	"context"
	"fmt"
	"log"
	"test-ebook-api/internal/pkg/storage"
	"test-ebook-api/internal/repository"
	"time"
)

// OrphanCleaner 孤儿文件清理器
type OrphanCleaner struct {
	staging      *storage.StagingStorage
	remote       storage.Storage
	repo         *repository.StandardRepository
	uploadRepo   *repository.UploadTaskRepository
	storageType  string
	interval     time.Duration // 定时清理间隔
	staleTimeout time.Duration // 暂存文件过期时间
	cancel       context.CancelFunc
}

// OrphanCleanerConfig 清理器配置
type OrphanCleanerConfig struct {
	Staging      *storage.StagingStorage
	Remote       storage.Storage
	Repo         *repository.StandardRepository
	UploadRepo   *repository.UploadTaskRepository
	StorageType  string
	IntervalH    int // 清理间隔（小时）
	StaleHours   int // 暂存文件过期时间（小时）
}

// OrphanScanResult 孤儿扫描结果
type OrphanScanResult struct {
	StagingCleaned    int      `json:"staging_cleaned"`    // 清理的暂存文件数量
	ExpiredCleaned    int      `json:"expired_cleaned"`    // 清理的过期软删除记录数量
	CloudOrphans      []string `json:"cloud_orphans"`      // 云端孤儿文件列表
	Errors            []string `json:"errors"`             // 错误信息
}

// NewOrphanCleaner 创建孤儿清理器
func NewOrphanCleaner(cfg OrphanCleanerConfig) *OrphanCleaner {
	interval := time.Duration(cfg.IntervalH) * time.Hour
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	staleTimeout := time.Duration(cfg.StaleHours) * time.Hour
	if staleTimeout <= 0 {
		staleTimeout = 24 * time.Hour
	}
	return &OrphanCleaner{
		staging:      cfg.Staging,
		remote:       cfg.Remote,
		repo:         cfg.Repo,
		uploadRepo:   cfg.UploadRepo,
		storageType:  cfg.StorageType,
		interval:     interval,
		staleTimeout: staleTimeout,
	}
}

// Start 启动定时清理（Level 1 + Level 2）
func (c *OrphanCleaner) Start(ctx context.Context) {
	ctx, c.cancel = context.WithCancel(ctx)
	go func() {
		log.Printf("[OrphanCleaner] 已启动，间隔 %v，暂存过期 %v", c.interval, c.staleTimeout)

		// 首次启动延迟 5 分钟再执行（等系统稳定）
		select {
		case <-time.After(5 * time.Minute):
		case <-ctx.Done():
			return
		}

		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		// 首次执行
		c.runAutoClean()

		for {
			select {
			case <-ctx.Done():
				log.Println("[OrphanCleaner] 收到停止信号，清理器退出")
				return
			case <-ticker.C:
				c.runAutoClean()
			}
		}
	}()
}

// Stop 停止清理器
func (c *OrphanCleaner) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

// runAutoClean 自动清理（Level 1 + Level 2）
func (c *OrphanCleaner) runAutoClean() {
	log.Println("[OrphanCleaner] 开始自动清理...")
	result := &OrphanScanResult{}

	c.cleanStaleStagingFiles(result)
	c.cleanExpiredSoftDeletes(result)

	log.Printf("[OrphanCleaner] 自动清理完成: 暂存清理=%d, 过期清理=%d, 错误=%d",
		result.StagingCleaned, result.ExpiredCleaned, len(result.Errors))
}

// RunManualScan 手动触发全量扫描（Level 1 + Level 2 + Level 3）
func (c *OrphanCleaner) RunManualScan() *OrphanScanResult {
	log.Println("[OrphanCleaner] 管理员触发手动扫描...")
	result := &OrphanScanResult{}

	c.cleanStaleStagingFiles(result)
	c.cleanExpiredSoftDeletes(result)
	c.scanCloudOrphans(result)

	log.Printf("[OrphanCleaner] 手动扫描完成: 暂存清理=%d, 过期清理=%d, 云端孤儿=%d, 错误=%d",
		result.StagingCleaned, result.ExpiredCleaned, len(result.CloudOrphans), len(result.Errors))

	return result
}

// cleanStaleStagingFiles Level 1: 清理残留暂存文件
func (c *OrphanCleaner) cleanStaleStagingFiles(result *OrphanScanResult) {
	staleFiles, err := c.staging.ListStale(c.staleTimeout)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("扫描暂存目录失败: %v", err))
		return
	}

	for _, filePath := range staleFiles {
		// 检查是否有对应的 upload_task 记录
		task, err := c.uploadRepo.GetTaskByLocalPath(filePath)
		if err != nil {
			// 无记录 → 孤儿文件，安全删除
			if removeErr := c.staging.Remove(filePath); removeErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("删除孤儿暂存文件失败 %s: %v", filePath, removeErr))
			} else {
				result.StagingCleaned++
				log.Printf("[OrphanCleaner] 删除孤儿暂存文件: %s", filePath)
			}
			continue
		}

		switch task.Status {
		case "synced":
			// 已同步成功，安全删除暂存文件
			if removeErr := c.staging.Remove(filePath); removeErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("删除已同步暂存文件失败 %s: %v", filePath, removeErr))
			} else {
				result.StagingCleaned++
				log.Printf("[OrphanCleaner] 删除已同步暂存文件: %s", filePath)
			}
		case "failed":
			if task.RetryCount >= task.MaxRetry {
				// 最终失败，保留文件但记录日志
				log.Printf("[OrphanCleaner] 保留最终失败任务的暂存文件: %s (任务 %d)", filePath, task.ID)
			}
			// 否则可能还在等待重试，不处理
		default:
			// pending/uploading 状态，不处理（可能正在同步中）
		}
	}
}

// cleanExpiredSoftDeletes Level 2: 清理过期的软删除记录
func (c *OrphanCleaner) cleanExpiredSoftDeletes(result *OrphanScanResult) {
	// 查询 30 天前软删除的文档
	cutoff := time.Now().AddDate(0, 0, -30)
	files, err := c.repo.GetExpiredSoftDeletedFiles(cutoff)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("查询过期软删除记录失败: %v", err))
		return
	}

	for _, f := range files {
		// 删除物理文件
		if f.FilePath != "" {
			if delErr := c.remote.Delete(f.FilePath); delErr != nil {
				log.Printf("[OrphanCleaner] 删除物理文件失败 %s: %v", f.FilePath, delErr)
				result.Errors = append(result.Errors, fmt.Sprintf("删除物理文件失败 %s: %v", f.FilePath, delErr))
			}
		}

		// 硬删除 DB 记录
		if delErr := c.repo.HardDeleteFileByID(f.ID); delErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("硬删除DB记录失败 ID=%d: %v", f.ID, delErr))
		} else {
			result.ExpiredCleaned++
			log.Printf("[OrphanCleaner] 硬删除过期文档: ID=%d, 路径=%s", f.ID, f.FilePath)
		}
	}
}

// scanCloudOrphans Level 3: 云端孤儿扫描（仅手动触发）
func (c *OrphanCleaner) scanCloudOrphans(result *OrphanScanResult) {
	if c.storageType == "local" {
		// 本地存储不需要云端扫描
		return
	}

	// 使用 ListObjects 接口列举云端文件（需要 Storage 接口支持）
	lister, ok := c.remote.(storage.ObjectLister)
	if !ok {
		result.Errors = append(result.Errors, "当前存储后端不支持对象列举，无法进行云端孤儿扫描")
		return
	}

	objects, err := lister.ListObjects()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("列举云端对象失败: %v", err))
		return
	}

	// 获取所有文档的 FilePath 集合
	allPaths, err := c.repo.GetAllFilePaths()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("获取文档路径列表失败: %v", err))
		return
	}

	pathSet := make(map[string]bool, len(allPaths))
	for _, p := range allPaths {
		pathSet[p] = true
	}

	// 比对
	for _, obj := range objects {
		if !pathSet[obj] {
			result.CloudOrphans = append(result.CloudOrphans, obj)
			log.Printf("[OrphanCleaner] 发现云端孤儿文件: %s", obj)
		}
	}
}
