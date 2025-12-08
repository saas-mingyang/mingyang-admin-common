package convert

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

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

func convertAnyToValue(v interface{}) (*structpb.Value, error) {
	switch val := v.(type) {
	case nil:
		return structpb.NewNullValue(), nil
	case bool:
		return structpb.NewBoolValue(val), nil
	case int:
		return structpb.NewNumberValue(float64(val)), nil
	case int32:
		return structpb.NewNumberValue(float64(val)), nil
	case int64:
		return structpb.NewNumberValue(float64(val)), nil
	case float32:
		return structpb.NewNumberValue(float64(val)), nil
	case float64:
		return structpb.NewNumberValue(val), nil
	case string:
		return structpb.NewStringValue(val), nil
	case []interface{}:
		// 处理数组
		list, err := convertSliceToValue(val)
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
		return convertAnyToValue(items)
	case map[string]interface{}:
		// 递归处理嵌套map
		s, err := convertMapToStructPB(val)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(s), nil
	default:
		// 尝试JSON序列化再解析
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var decoded interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			return nil, err
		}
		return convertAnyToValue(decoded)
	}
}

// convertMapToStructPB 将 map[string]interface{} 转换为 *structpb.Struct
func convertMapToStructPB(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}

	// 使用 structpb 的标准转换方法
	return structpb.NewStruct(m)
}

func convertSliceToValue(items []interface{}) (*structpb.ListValue, error) {
	values := make([]*structpb.Value, 0, len(items))
	for _, item := range items {
		value, err := convertAnyToValue(item)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return &structpb.ListValue{Values: values}, nil
}

func convertMapToStructPBAdvanced(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}

	fields := make(map[string]*structpb.Value)
	for k, v := range m {
		value, err := convertAnyToValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", k, err)
		}
		fields[k] = value
	}
	return &structpb.Struct{Fields: fields}, nil
}
