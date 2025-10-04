package migratortype

type (
	migratorType int8
)

const (
	unknown migratorType = iota // invalid
	goose
)
