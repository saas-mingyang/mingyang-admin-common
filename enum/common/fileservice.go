package common

import "time"

const (

	// FileServiceURL is the base URL for the file service
	FileServicePresignedUploadURL = "/presigned/upload_url"
	// FileServicePresignedDownloadURL is the base URL for the file service
	FileServicePresignedDownloadURL = "/presigned/download_url"

	FileServiceDefaultDeviceID    = "fileservice"
	FileServiceDefaultContentType = "application/octet-stream"

	FileServiceDialTimeout    = 15 * time.Second
	FileServiceRequestTimeout = 30 * time.Second
)
