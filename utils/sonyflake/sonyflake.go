package sonyflake

import (
	"errors"
	"sync"
	"time"
)

const (
	twepoch        = int64(1483228800000)             //开始时间截 (2017-01-01)
	workeridBits   = uint(10)                         //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	workeridMax    = int64(-1 ^ (-1 << workeridBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	workeridShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workeridBits      //时间戳左移位数
	workId         = int64(0)
)

// A Snowflake struct holds the basic information needed for a snowflake generator worker
type Snowflake struct {
	sync.Mutex
	timestamp int64
	workerid  int64
	sequence  int64
}

// 全局默认实例
var defaultSnowflake *Snowflake

// 初始化默认实例
func init() {
	sf, _ := NewSnowflake(workId)
	defaultSnowflake = sf
}

// NextID 便捷函数，直接生成单个ID
func NextID() int64 {
	return defaultSnowflake.Generate()
}

// BatchNextID 便捷函数，批量生成ID
func BatchNextID(count int) []int64 {
	return defaultSnowflake.BatchGenerate(count)
}

// NewSnowflake NewNode returns a new snowflake worker that can be used to generate snowflake IDs
func NewSnowflake(workerid int64) (*Snowflake, error) {
	if workerid < 0 || workerid > workeridMax {
		return nil, errors.New("workerid must be between 0 and 1023")
	}

	return &Snowflake{
		timestamp: 0,
		workerid:  workerid,
		sequence:  0,
	}, nil
}

// MustSnowflake 创建Snowflake实例，如果出错则panic
func MustSnowflake(workerid int64) *Snowflake {
	sf, err := NewSnowflake(workerid)
	if err != nil {
		panic(err)
	}
	return sf
}

// Generate creates and returns a unique snowflake ID
func (s *Snowflake) Generate() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1000000

	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	r := (now-twepoch)<<timestampShift | (s.workerid << workeridShift) | (s.sequence)
	return r
}

// BatchGenerate 批量生成唯一snowflake ID
func (s *Snowflake) BatchGenerate(count int) []int64 {
	if count <= 0 {
		return []int64{}
	}

	s.Lock()
	defer s.Unlock()

	ids := make([]int64, count)
	now := time.Now().UnixNano() / 1000000

	// 统一处理逻辑
	if s.timestamp != now {
		s.sequence = 0
		s.timestamp = now
	}

	for i := 0; i < count; i++ {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 序列号用完，等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
			s.timestamp = now
		}
		ids[i] = int64((now-twepoch)<<timestampShift | (s.workerid << workeridShift) | (s.sequence))
	}

	return ids
}
