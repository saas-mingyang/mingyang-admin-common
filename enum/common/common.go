package common

const (
	StatusNormal uint8 = 1
	StatusBanned uint8 = 2
)

const (
	DefaultParentId uint64 = 1000000
	// DefaultInvalidRoleId used to judge whether the token belongs to core, if it is DefaultInvalidRoleId, it belongs to mms
	DefaultInvalidRoleId uint64 = 1000000
	DateFormat                  = "2006-01-02"
	DateTimeFormat              = "2006-01-02 15:04:05"
	TimeFormat                  = "15:04:05"
)
