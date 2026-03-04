package convert

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"strconv"
	"strings"
	"time"
)

// Converter 提供 Protobuf 结构体转换功能
type Converter struct{}

// NewConverter 创建新的转换器实例
func NewConverter() *Converter {
	return &Converter{}
}

// TimestampUnit 定义时间戳单位常量
type TimestampUnit string

const (
	UnitSecond      TimestampUnit = "second"
	UnitMillisecond TimestampUnit = "millisecond"
	UnitMicrosecond TimestampUnit = "microsecond"
)

// TimestampToTime 将时间戳转换为time.Time
// unit: 时间戳单位（秒、毫秒、微秒）
func TimestampToTime(timestamp int64, unit TimestampUnit) *time.Time {
	if timestamp == 0 {
		return nil
	}

	var t time.Time

	switch unit {
	case UnitMillisecond:
		// 处理毫秒，兼容负数
		sec := timestamp / 1000
		nsec := (timestamp % 1000) * int64(time.Millisecond)
		t = time.Unix(sec, nsec)
	case UnitMicrosecond:
		// 处理微秒，兼容负数
		sec := timestamp / 1e6
		nsec := (timestamp % 1e6) * int64(time.Microsecond)
		t = time.Unix(sec, nsec)
	case UnitSecond:
		fallthrough
	default:
		t = time.Unix(timestamp, 0)
	}

	return &t
}
func TimeFromFrontendMillis(millis int64, timeUnit TimestampUnit) *time.Time {
	return TimestampToTime(millis, timeUnit)
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

// MapToStruct 将 map[string]interface{} 转换为 *structpb.Struct
// 使用 structpb.NewStruct 的简单版本
func (c *Converter) MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}
	return structpb.NewStruct(m)
}

// MapToStructAdvanced 将 map[string]interface{} 转换为 *structpb.Struct
// 支持更多数据类型的高级版本
func (c *Converter) MapToStructAdvanced(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}

	fields := make(map[string]*structpb.Value)
	for k, v := range m {
		value, err := c.anyToValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", k, err)
		}
		fields[k] = value
	}
	return &structpb.Struct{Fields: fields}, nil
}

// anyToValue 将任意类型转换为 *structpb.Value
func (c *Converter) anyToValue(v interface{}) (*structpb.Value, error) {
	switch val := v.(type) {
	case nil:
		return structpb.NewNullValue(), nil
	case bool:
		return structpb.NewBoolValue(val), nil
	case int:
		return structpb.NewNumberValue(float64(val)), nil
	case int8:
		return structpb.NewNumberValue(float64(val)), nil
	case int16:
		return structpb.NewNumberValue(float64(val)), nil
	case int32:
		return structpb.NewNumberValue(float64(val)), nil
	case int64:
		return structpb.NewNumberValue(float64(val)), nil
	case uint:
		return structpb.NewNumberValue(float64(val)), nil
	case uint8:
		return structpb.NewNumberValue(float64(val)), nil
	case uint16:
		return structpb.NewNumberValue(float64(val)), nil
	case uint32:
		return structpb.NewNumberValue(float64(val)), nil
	case uint64:
		return structpb.NewNumberValue(float64(val)), nil
	case float32:
		return structpb.NewNumberValue(float64(val)), nil
	case float64:
		return structpb.NewNumberValue(val), nil
	case string:
		return structpb.NewStringValue(val), nil
	case []interface{}:
		// 处理数组
		list, err := c.sliceToValue(val)
		if err != nil {
			return nil, err
		}
		return structpb.NewListValue(list), nil
	case []map[string]interface{}:
		// 处理map数组
		var items []interface{}
		for _, m := range val {
			items = append(items, m)
		}
		return c.anyToValue(items)
	case map[string]interface{}:
		// 递归处理嵌套map
		s, err := c.MapToStructAdvanced(val)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(s), nil
	default:
		// 尝试JSON序列化再解析
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
		var decoded interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			return nil, err
		}
		return c.anyToValue(decoded)
	}
}

// ValueToAny 将 *structpb.Value 转换为 interface{}
func (c *Converter) ValueToAny(v *structpb.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.GetKind().(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_StructValue:
		return c.StructToMap(v.GetStructValue())
	case *structpb.Value_ListValue:
		list := v.GetListValue()
		result := make([]interface{}, len(list.GetValues()))
		for i, item := range list.GetValues() {
			result[i] = c.ValueToAny(item)
		}
		return result
	default:
		return nil
	}
}

// IsEmptyStruct 检查 structpb.Struct 是否为空
func (c *Converter) IsEmptyStruct(s *structpb.Struct) bool {
	return s == nil || len(s.Fields) == 0
}

// MergeStructs 合并多个 structpb.Struct
func (c *Converter) MergeStructs(structs ...*structpb.Struct) (*structpb.Struct, error) {
	merged := make(map[string]*structpb.Value)

	for _, s := range structs {
		if s == nil {
			continue
		}
		for k, v := range s.Fields {
			merged[k] = v
		}
	}

	return &structpb.Struct{Fields: merged}, nil
}

// ConvertWithDefault 转换 map 到 structpb，提供默认值
func (c *Converter) ConvertWithDefault(m map[string]interface{}, defaultValue *structpb.Struct) *structpb.Struct {
	if m == nil || len(m) == 0 {
		return defaultValue
	}

	result, err := c.MapToStructAdvanced(m)
	if err != nil {
		// 如果转换失败，返回默认值
		return defaultValue
	}

	return result
}

// sliceToValue 将 []interface{} 转换为 *structpb.ListValue
func (c *Converter) sliceToValue(items []interface{}) (*structpb.ListValue, error) {
	values := make([]*structpb.Value, 0, len(items))
	for _, item := range items {
		value, err := c.anyToValue(item)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return &structpb.ListValue{Values: values}, nil
}

// StructToMap 将 *structpb.Struct 转换为 map[string]interface{}
func (c *Converter) StructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// StringSliceToUint64Slice 将字符串切片转换为 uint64 切片
func StringSliceToUint64Slice(strSlice []string) []uint64 {
	fmt.Printf("StringSliceToUint64Slice: %v\n", strSlice)
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
	fmt.Printf("Uint64SliceToStringSlice: %v\n", uintSlice)
	result := make([]string, len(uintSlice))

	for i, val := range uintSlice {
		result[i] = strconv.FormatUint(val, 10)
	}

	return result
}

// =============================
// uint64 -> []string
// =============================

// Uint64ToStrings 将多个 uint64 转换为 []string
func Uint64ToStrings(ids ...uint64) []string {
	if len(ids) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, strconv.FormatUint(id, 10))
	}
	return result
}

// Uint64SliceToStrings 支持直接传 []uint64
func Uint64SliceToStrings(ids []uint64) []string {
	return Uint64ToStrings(ids...)
}

// =============================
// []string -> []uint64
// =============================

// StringsToUint64 将多个 string 转换为 []uint64
func StringsToUint64(strs ...string) ([]uint64, error) {
	if len(strs) == 0 {
		return []uint64{}, nil
	}

	result := make([]uint64, 0, len(strs))
	for _, s := range strs {
		if s == "" {
			continue
		}

		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

// StringSliceToUint64 支持直接传 []string
func StringSliceToUint64(strs []string) ([]uint64, error) {
	return StringsToUint64(strs...)
}

// =============================
// 逗号字符串互转
// =============================

// Uint64ToJoinString 1,2,3
func Uint64ToJoinString(ids ...uint64) string {
	return strings.Join(Uint64ToStrings(ids...), ",")
}

// JoinStringToUint64 "1,2,3"
func JoinStringToUint64(s string) ([]uint64, error) {
	if s == "" {
		return []uint64{}, nil
	}
	parts := strings.Split(s, ",")
	return StringSliceToUint64(parts)
}

// MapToStruct 将 map[string]interface{} 解码为指定结构体
func MapToStruct(cfg map[string]interface{}, target interface{}) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

// StructToMap 将结构体转换为 map[string]interface{}
func StructToMap(v interface{}) (map[string]interface{}, error) {
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// 定义一个 map 用来存储反序列化后的结果
	var resultMap map[string]interface{}

	// 将 JSON 字符串反序列化为 map
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

// ToStringMap 将 map[string]interface{} 转换为 map[string]string
func ToStringMap(original map[string]interface{}) (map[string]string, error) {
	converted := make(map[string]string)
	for key, value := range original {
		// 如果值是 map 类型，则序列化为 JSON 字符串
		switch v := value.(type) {
		case map[string]interface{}:
			// 序列化为 JSON 字符串
			jsonValue, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal value: %v", err)
			}
			converted[key] = string(jsonValue)
		default:
			// 否则，直接转换为字符串
			converted[key] = fmt.Sprintf("%v", value)
		}
	}
	return converted, nil
}
