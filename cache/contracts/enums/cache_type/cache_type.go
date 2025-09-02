package cachetype

//go:generate goenums -f -c cache_type.go

type (
	cacheType int8
)

const (
	unknown cacheType = iota // invalid unknown
	in_memory
	redis
)
