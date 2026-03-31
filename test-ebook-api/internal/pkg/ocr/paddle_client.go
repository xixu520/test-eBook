package ocr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"test-ebook-api/internal/config"
	"time"
)

type PaddleClient struct {
	token           string
	model           string
	url             string
	optionalPayload map[string]interface{}
	client          *http.Client // 用于轮询等短请求
	uploadClient    *http.Client // 用于提交/下载大文件，超时更长
}

func NewPaddleClient() *PaddleClient {
	cfg := config.GlobalConfig.OCR
	return NewPaddleClientWithParams(cfg.PaddleToken, cfg.PaddleJobURL)
}

func NewPaddleClientWithParams(token, url string) *PaddleClient {
	cfg := config.GlobalConfig.OCR
	return &PaddleClient{
		token: token,
		model: cfg.PaddleModel,
		url:   url,
		optionalPayload: map[string]interface{}{
			"useDocOrientationClassify": cfg.UseDocOrientationClassify,
			"useDocUnwarping":           cfg.UseDocUnwarping,
			"useChartRecognition":       cfg.UseChartRecognition,
		},
		client:       &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second},
		uploadClient: &http.Client{Timeout: 10 * time.Minute}, // 上传/下载大文件用更长超时
	}
}

func (p *PaddleClient) SetModel(model string) {
	p.model = model
}

type V2JobResponse struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Data      struct {
		JobID string `json:"jobId"`
	} `json:"data"`
}

type V2ResultResponse struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Data      struct {
		State     string `json:"state"` // pending, running, done, failed
		ResultURL struct {
			JSONURL string `json:"jsonUrl"`
		} `json:"resultUrl"`
		ErrorMsg string `json:"errorMsg"`
	} `json:"data"`
}

func (p *PaddleClient) SubmitTask(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	_ = writer.WriteField("model", p.model)
	optJSON, _ := json.Marshal(p.optionalPayload)
	_ = writer.WriteField("optionalPayload", string(optJSON))

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "bearer "+p.token)

	resp, err := p.uploadClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res V2JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.ErrorCode != 0 {
		return "", errors.New(res.ErrorMsg)
	}

	return res.Data.JobID, nil
}

func (p *PaddleClient) GetResult(jobID string) (string, string, error) {
	url := fmt.Sprintf("%s/%s", p.url, jobID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var res V2ResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	switch res.Data.State {
	case "done":
		content, err := p.downloadResult(res.Data.ResultURL.JSONURL)
		if err != nil {
			return "", "failed", err
		}
		return content, "success", nil
	case "failed":
		return "", "failed", errors.New(res.Data.ErrorMsg)
	case "pending", "running":
		return "", "processing", nil
	default:
		return "", "processing", nil
	}
}

func (p *PaddleClient) downloadResult(jsonlURL string) (string, error) {
	resp, err := p.uploadClient.Get(jsonlURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var pageResult struct {
			Result struct {
				LayoutParsingResults []struct {
					Markdown struct {
						Text string `json:"text"`
					} `json:"markdown"`
				} `json:"layoutParsingResults"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(line), &pageResult); err != nil {
			continue
		}
		for _, res := range pageResult.Result.LayoutParsingResults {
			sb.WriteString(res.Markdown.Text)
			sb.WriteString("\n\n")
		}
	}
	return sb.String(), nil
}

func (p *PaddleClient) TestConnection() error {
	// 使用 fileUrl 模式发送一个最小化的 POST 请求来验证 Token
	payload := map[string]interface{}{
		"fileUrl": "https://invalid-test.example.com/test.pdf",
		"model":   p.model,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", p.url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("构建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("网络连接失败: %v", err)
	}
	defer resp.Body.Close()

	// 根据 PaddleOCR 官方 API 文档及实际测试，按 HTTP 状态码精确匹配错误
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case 400:
			if strings.Contains(string(bodyBytes), "10005") {
				return nil // ErrorCode/Code 10005 表示“文件内容无法解析”，这意味着 Token 已通过且模型合法
			}
			return fmt.Errorf("请求体参数错误 (HTTP 400): %s", string(bodyBytes))
		case 401, 403:
			return errors.New("Token 错误: 请检查 Token 是否正确，或 URL 是否与 Token 匹配")
		case 413:
			return errors.New("请求体过大: 请减少文件的页数或大小")
		case 422:
			return errors.New("参数无效: 请检查请求参数格式")
		case 429:
			return errors.New("已超出单日解析最大页数 (3000页): 请使用其他模型或稍后再试")
		case 500:
			return errors.New("服务器内部错误: 请稍后再试或联系 PaddleOCR 官方")
		case 503:
			return errors.New("当前请求过多: 请稍后再试")
		case 504:
			return errors.New("网关超时: 请稍后再试")
		default:
			return fmt.Errorf("服务器响应异常 (HTTP %d), body: %s", resp.StatusCode, string(bodyBytes))
		}
	}

	// 解析 JSON 响应体
	var res V2JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("服务端响应格式错误: %v", err)
	}

	// 严格校验业务错误码 — 只有 errorCode 为 0 才视为 Token 有效
	if res.ErrorCode != 0 {
		return fmt.Errorf("连接测试未通过: %s (错误码: %d)", res.ErrorMsg, res.ErrorCode)
	}

	return nil // Token 有效，任务已成功创建
}
