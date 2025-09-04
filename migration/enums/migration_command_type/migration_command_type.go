package migrationcommandtype

type (
	migrationCommandType int8
)

const (
	unknown migrationCommandType = iota // invalid
	up
	down
)
