package cleanuppolicy

type (
	cleanupPolicy uint8
)

const (
	unknown cleanupPolicy = iota // invalid
	lru
	lfu
	tinyLfu
	fifi
)
