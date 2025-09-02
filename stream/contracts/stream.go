package contracts

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	streamclienttype "frisboo-bank/pkg/stream/contracts/enums/stream_client_type"
)

type (
	Stream interface {
		Type() streamclienttype.StreamClientType
		Logger() loggerContracts.Logger
	}

	StreamAdapter interface{}
)
