package rpcservertype

//go:generate goenums -f -c rpc_server_type.go

type (
	rpcServerType int8
)

const (
	unknown rpcServerType = iota // invalid unknown
	grpc
)
