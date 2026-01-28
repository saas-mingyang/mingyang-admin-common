package convert

// TimeRange 时间范围
type TimeRange struct {
	StartTime *int64 `json:"startTime,omitempty"`
	EndTime   *int64 `json:"endTime,omitempty"`
}

// NewTimeRange 根据数组初始化 TimeRange
func NewTimeRange(arr []*int64) (startTime, endTime *int64) {
	tr := &TimeRange{}
	if len(arr) >= 1 {
		tr.StartTime = arr[0]
	}
	if len(arr) >= 2 {
		tr.EndTime = arr[1]
	}
	return tr.StartTime, tr.EndTime
}
