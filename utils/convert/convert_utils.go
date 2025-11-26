package convert

import "strconv"

// StringSliceToUint64Slice 将字符串切片转换为 uint64 切片
func StringSliceToUint64Slice(strSlice []string) []uint64 {
	result := make([]uint64, 0, len(strSlice))

	for _, str := range strSlice {
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return result
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
