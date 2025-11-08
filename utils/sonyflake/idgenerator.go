package sonyflake

import (
	"github.com/emicklei/go-restful/v3/log"
	"github.com/sony/sonyflake"
	"sync"
)

var (
	sf     *sonyflake.Sonyflake
	sfOnce sync.Once
)

func init() {
	sfOnce.Do(func() {
		var st sonyflake.Settings
		sf = sonyflake.NewSonyflake(st)
		if sf == nil {
			panic("sonyflake not working in this machine")
		}
	})
}

// GenerateID 生成一个新的雪花 ID
func GenerateID() (uint64, error) {
	id, err := sf.NextID()
	if err != nil {
		log.Printf("Error generating sonyflake ID: %v", err)
		return 0, err
	}
	log.Printf("Generated Sonyflake ID: %d", id)
	return id, nil
}
