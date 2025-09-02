package migrationcommandtype

//go:generate goenums -f -c migration_command_type.go

type (
	migrationCommandType int8
)

const (
	unknown migrationCommandType = iota // invalid unknown
	up
	down
)
