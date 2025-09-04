package cachetype

type (
	cacheType int8
)

const (
	unknown cacheType = iota // invalid unknown
	in_memory
	redis
)
