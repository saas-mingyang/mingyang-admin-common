package convert

import "strconv"

// StringSliceToUint64Slice 将字符串切片转换为 uint64 切片
func StringSliceToUint64Slice(strSlice []string) []uint64 {
	result := make([]uint64, 0, len(strSlice))

	for _, str := range strSlice {
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			panic(err)
		}
		result = append(result, val)
	}

	return result
}

// Uint64SliceToStringSlice 将 uint64 切片转换为字符串切片
func Uint64SliceToStringSlice(uintSlice []uint64) []string {
	result := make([]string, len(uintSlice))

	for i, val := range uintSlice {
		result[i] = strconv.FormatUint(val, 10)
	}

	return result
}

// GetUint8FromProto 从 protobuf 中获取 uint8
func GetUint8FromProto(protoField *uint32) *uint8 {
	if protoField == nil {
		return nil
	}
	val := uint8(*protoField)
	return &val
}

// GetUint32FromProto 从 protobuf 中获取 uint32
func GetUint32FromProto(protoField *uint8) *uint32 {
	if protoField == nil {
		return nil
	}
	val := uint32(*protoField)
	return &val
}
