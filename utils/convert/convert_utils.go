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

// Int32ToUint8Ptr Int32ToUint8PtrPool Int32ToUint8Ptr 辅助函数：int32转指针，超出范围返回nil
func Int32ToUint8Ptr(i uint32) *uint8 {
	if i < 0 || i > 255 {
		return nil
	}
	u := uint8(i)
	return &u
}

// Uint8PtrToUint32  Uint8PtrToUint32 辅助函数：指针转uint32，为nil返回0
func Uint8PtrToUint32(val *uint8) *uint32 {
	if val == nil {
		return nil
	}
	u := uint32(*val)
	return &u
}
