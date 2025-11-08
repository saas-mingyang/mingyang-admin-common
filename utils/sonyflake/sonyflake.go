package sonyflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"sync"
	"time"
)

type IDGenerator struct {
	sf *sonyflake.Sonyflake
	mu sync.Mutex
}

// NewIDGenerator 创建一个 Sonyflake ID 生成器
func NewIDGenerator() *IDGenerator {
	st := sonyflake.Settings{
		// 设置起始时间，使生成的数字偏小，利于首位1
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		// 固定机器ID，保证分布式安全
		MachineID: func() (uint16, error) {
			return 1, nil
		},
	}
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
	return &IDGenerator{
		sf: sf,
	}
}

func (g *IDGenerator) NextID() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	id, err := g.sf.NextID()
	if err != nil {
		panic(fmt.Sprintf("failed to generate ID: %v", err))
	}
	id18 := id % 1e18
	return fmt.Sprintf("1%018d", id18)
}

func (g *IDGenerator) NextBatch(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = g.NextID()
	}
	return ids
}
