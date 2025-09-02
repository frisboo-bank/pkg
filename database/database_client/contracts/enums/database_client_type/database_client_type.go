package databaseclienttype

//go:generate goenums -f -c database_client_type.go

type (
	databaseClientType int8
)

const (
	unknown databaseClientType = iota // invalid unknown
	postgres
)
