package cleanuppolicy

//go:generate goenums -f -c cleanup_policy.go

type (
	cleanupPolicy int8
)

const (
	unknown cleanupPolicy = iota // invalid unknown
	lru
	lfu
	tinyLfu
	fifi
)
