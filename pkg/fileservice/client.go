package fileservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/enum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

type uploadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url string `json:"url"`
		Key string `json:"key"`
	} `json:"data"`
}

type downloadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
}

// Upload 上传文件到文件服务，返回 object key 和预签名下载地址。
// baseURL 是文件服务地址（如 http://localhost:9102），data 是文件内容，fileName 是文件名。
// contentType 和 deviceId 可传空字符串，使用默认值。
func Upload(ctx context.Context, baseURL string, data []byte, fileName, contentType, deviceId string, expiresIn int64) (objectKey, downloadURL string, err error) {
	if contentType == common.EmptyString {
		contentType = common.FileServiceDefaultContentType
	}
	if deviceId == common.EmptyString {
		deviceId = common.FileServiceDefaultDeviceID
	}

	presignedURL, objectKey, err := getPresignedUploadURL(ctx, baseURL, fileName, int64(len(data)), contentType, deviceId, expiresIn)

	logx.Infof("presignedURL=%s", presignedURL)

	if err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	if err := putRaw(ctx, presignedURL, data, contentType); err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	downloadURL, err = getDownloadURL(ctx, baseURL, objectKey, expiresIn)
	return objectKey, downloadURL, err
}

// GetDownloadURL 获取预签名下载地址。
// baseURL 是文件服务地址，key 是文件在对象存储中的 key。
func GetDownloadURL(ctx context.Context, baseURL, key string, expiresIn int64) (downloadURL string, err error) {
	return getDownloadURL(ctx, baseURL, key, expiresIn)
}

// getPresignedUploadURL 向文件服务申请预签名上传地址。
func getPresignedUploadURL(ctx context.Context, baseURL string, fileName string, fileSize int64, contentType, deviceId string, expiresIn int64) (presignedURL, key string, err error) {

	//打baseUrl
	fmt.Printf("baseURL=%s", baseURL)

	body, _ := json.Marshal(map[string]interface{}{
		"fileName":    fileName,
		"fileSize":    fileSize,
		"contentType": contentType,
		"deviceId":    deviceId,
		"expiresIn":   expiresIn,
	})

	ctx, cancel := context.WithTimeout(ctx, common.FileServiceRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+common.FileServicePresignedUploadURL, bytes.NewReader(body))
	if err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	req.Header.Set(common.ContentType, "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {

		return common.EmptyString, common.EmptyString, fmt.Errorf("presigned upload API status %d: %s", resp.StatusCode, string(respBody))
	}

	var result uploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return common.EmptyString, common.EmptyString, err
	}
	logx.Infof(
		"upload url=%s key=%s",
		result.Data.Url,
		result.Data.Key,
	)

	if result.Code != common.Zero {
		return common.EmptyString, common.EmptyString, fmt.Errorf("presigned upload API error: %s", result.Msg)
	}
	return result.Data.Url, result.Data.Key, nil
}

// getDownloadURL 获取预签名下载地址。
func getDownloadURL(ctx context.Context, baseURL, key string, expiresIn int64) (downloadURL string, err error) {

	body, _ := json.Marshal(map[string]interface{}{
		"key":       key,
		"expiresIn": expiresIn,
	})

	ctx, cancel := context.WithTimeout(ctx, common.FileServiceRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+common.FileServicePresignedDownloadURL, bytes.NewReader(body))
	if err != nil {
		return common.EmptyString, err
	}
	req.Header.Set(common.ContentType, "application/json")

	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("resp: %v\n", resp)

	if err != nil {
		return common.EmptyString, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return common.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return common.EmptyString, fmt.Errorf("presigned download API status %d: %s", resp.StatusCode, string(respBody))
	}

	var result downloadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return common.EmptyString, err
	}
	if result.Code != common.Zero {
		return common.EmptyString, fmt.Errorf("presigned download API error: %s", result.Msg)
	}
	return result.Data.Url, nil
}

// putRaw 通过原始 TCP/TLS 连接发送 PUT 请求，避免 Go http.Client 自动添加的非签名头。
func putRaw(ctx context.Context, rawURL string, data []byte, contentType string) error {
	ctx, cancel := context.WithTimeout(ctx, common.FileServiceRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		rawURL,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set(common.ContentType, contentType)

	client := &http.Client{
		Timeout: common.FileServiceRequestTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed, status=%d, body=%s", resp.StatusCode, string(body))
	}

	return nil
}
