package rpcservertype

type (
	rpcServerType int8
)

const (
	unknown rpcServerType = iota // unknown
	invalid rpcServerType = iota // invalid
	grpc
)
