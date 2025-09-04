package databaseclienttype

type (
	databaseClientType int8
)

const (
	unknown databaseClientType = iota // invalid
	postgres
)
