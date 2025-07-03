package loglevel

//go:generate goenums -f log_level.go

type (
	logLevel int8
)

const (
	unknown_level logLevel = iota // invalid unknown
	debug_level                   // debug
	info_level                    // info
	warn_level                    // warn
	error_level                   // error
	panic_level                   // panic
	fatal_level                   // fatal
)
