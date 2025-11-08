package sonyflake

import (
	"fmt"
	"testing"
)

func TestGetSonyflakeID(t *testing.T) {
	ids := BatchNextID(100000)
	fmt.Println(ids)
}
