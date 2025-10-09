package routing

import (
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/validation"
)

var _ contracts.RouteBuilder = (*routeBuilder)(nil)

type routeBuilder struct {
	engine contracts.RouterEngine
}

func NewBuilder(engine contracts.RouterEngine) contracts.RouteBuilder {
	validation.AssertNotNil("engine", engine)

	return &routeBuilder{engine: engine}
}

func (r *routeBuilder) Root() contracts.RouteGroup {
	return newRouteGroup(r.engine)
}
