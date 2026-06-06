package fileservice

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/saas-mingyang/mingyang-admin-common/enum/common"
	"io"
	"net"
	"net/http"
	"net/url"
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
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	var conn net.Conn
	d := &net.Dialer{Timeout: common.FileServiceDialTimeout}

	ctx, cancel := context.WithTimeout(ctx, common.FileServiceRequestTimeout)
	defer cancel()

	if u.Scheme == "https" {
		conn, err = tls.DialWithDialer(d, "tcp", u.Host, &tls.Config{InsecureSkipVerify: false})
	} else {
		conn, err = d.DialContext(ctx, "tcp", u.Host)
	}
	if err != nil {
		return err
	}
	defer conn.Close()

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("PUT %s HTTP/1.1\r\n", u.RequestURI()))
	buf.WriteString(fmt.Sprintf("Host: %s\r\n", u.Host))
	buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data)))
	buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	buf.WriteString("\r\n")
	buf.Write(data)

	if _, err = conn.Write(buf.Bytes()); err != nil {
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("presigned URL PUT status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
