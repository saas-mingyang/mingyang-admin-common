package sonyflake

import (
	"fmt"
	"testing"
)

func TestGetSonyflakeID(t *testing.T) {
	/*ids := BatchNextID(100000)
	fmt.Println(ids)*/
	id := int64(1171789827910142449)
	uintId := uint64(id)

	// 验证转换前后数值相等
	if id != int64(uintId) {
		t.Errorf("转换后精度丢失: %d != %d", id, int64(uintId))
	}

	fmt.Printf("原始int64: %d\n", id)
	fmt.Printf("转换后uint64: %d\n", uintId)

}
