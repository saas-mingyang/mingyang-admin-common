package convert

import "time"

// TimePtrFromUnix 辅助函数：时间戳转指针，为0返回nil
func TimePtrFromUnix(unix int64) *time.Time {
	if unix == 0 {
		return nil
	}
	t := time.Unix(unix, 0)
	return &t
}
