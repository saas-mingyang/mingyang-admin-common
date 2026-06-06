package common

import "time"

const (
	// 文件服务 API 路径
	FileServicePresignedUploadURL   = "/presigned/upload_url"
	FileServicePresignedDownloadURL = "/presigned/download_url"

	FileServiceDefaultDeviceID    = "fileservice"
	FileServiceDefaultContentType = "application/octet-stream"

	// 超时
	FileServiceDialTimeout    = 15 * time.Second
	FileServiceRequestTimeout = 30 * time.Second
)
