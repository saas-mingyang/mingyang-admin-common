package sonyflake

import (
	"fmt"
	"testing"
)

func TestGetSonyflakeID(t *testing.T) {
	generator := NewIDGenerator()
	batch := generator.NextBatch(100)
	fmt.Printf("sonyflake id: %s", batch)
}
