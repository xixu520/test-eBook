package ocr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"test-ebook-api/internal/config"
	"time"
)

type PaddleClient struct {
	token  string
	model  string
	url    string
	client *http.Client
}

func NewPaddleClient() *PaddleClient {
	cfg := config.GlobalConfig.OCR
	return &PaddleClient{
		token:  cfg.PaddleToken,
		model:  cfg.PaddleModel,
		url:    cfg.PaddleJobURL,
		client: &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second},
	}
}

type JobResponse struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Data      struct {
		JobID string `json:"jobId"`
	} `json:"data"`
}

type ResultResponse struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Data      struct {
		Status int    `json:"status"` // 1: 成功, 其他: 处理中或失败
		Result string `json:"result"`
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
	
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	io.Copy(part, file)
	
	writer.WriteField("model", p.model)
	
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jobResp JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return "", err
	}

	if jobResp.ErrorCode != 0 {
		return "", errors.New(jobResp.ErrorMsg)
	}

	return jobResp.Data.JobID, nil
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

	var resResp ResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&resResp); err != nil {
		return "", "", err
	}
	fmt.Printf("[PaddleAPI] JobID: %s, Status: %d, Message: %s\n", jobID, resResp.Data.Status, resResp.ErrorMsg)

	if resResp.ErrorCode != 0 {
		return "", "failed", errors.New(resResp.ErrorMsg)
	}

	if resResp.Data.Status == 1 {
		return resResp.Data.Result, "success", nil
	} else if resResp.Data.Status == 2 {
		return "", "failed", nil
	}
	return "", "processing", nil
}
