package filex

import (
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"mingyang.com/admin-simple-admin-file/internal/enum/filetype"
)

// ConvertFileTypeToUint8 converts file type string to uint8.
func ConvertFileTypeToUint8(fileType string) uint8 {
	switch fileType {
	case "other":
		return filetype.Other
	case "image":
		return filetype.Image
	case "video":
		return filetype.Video
	case "audio":
		return filetype.Audio
	default:
		return filetype.Other
	}
}

func ConvertUrlStringToFileUint64(urlStr string) (uint64, error) {
	urlData, err := url.Parse(urlStr)
	if err != nil {
		logx.Error("failed to parse url", logx.Field("details", err), logx.Field("data", urlStr))
		return 0, err
	}

	fileId := filepath.Base(urlData.Path)

	if len(fileId) >= 36 {
		fileId = fileId[:36]
	} else if len(fileId) < 36 {
		return 0, errorx.NewApiBadRequestError("wrong file path")
	}

	id, err := ParseSnowflakeID(fileId)
	if err != nil {
		logx.Error("failed to parse snowflake id", logx.Field("details", err), logx.Field("data", fileId))
		return 0, errorx.NewApiBadRequestError("wrong file path")
	}
	return id, nil
}

// ParseSnowflakeID 字符串雪花ID转uint64
func ParseSnowflakeID(snowflakeStr string) (uint64, error) {
	return strconv.ParseUint(snowflakeStr, 10, 64)
}
