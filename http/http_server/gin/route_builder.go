package gin

import (
	"fmt"

	"frisboo-bank/pkg/http/http_server/contracts"

	"github.com/gin-gonic/gin"
)

type routeBuilder struct {
	engine *gin.Engine
	groups map[string]*gin.RouterGroup
}

var _ contracts.RouteBuilder = (*routeBuilder)(nil)

func NewRouteBuilder(engine *gin.Engine) contracts.RouteBuilder {
	return &routeBuilder{
		engine: engine,
	}
}

func (r *routeBuilder) Build() any {
	return r.engine
}

func (r *routeBuilder) RegisterGroup(groupName string) contracts.RouteBuilder {
	if _, exists := r.groups[groupName]; !exists {
		r.groups[groupName] = r.engine.Group(groupName)
	}

	return r
}

func (r *routeBuilder) RegisterGroupFunc(groupName string, builder func(group any)) contracts.RouteBuilder {
	r.RegisterGroup(groupName)
	builder(r.groups[groupName])

	return r
}

func (r *routeBuilder) RegisterRoutes(builder func(server any)) contracts.RouteBuilder {
	builder(r.engine)

	return r
}

func (r *routeBuilder) UseForGroup(groupName string, middlewares ...any) contracts.RouteBuilder {
	r.RegisterGroup(groupName)

	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(fmt.Errorf("gin-server: invalid middleware for group: `%s` `%v`", groupName, err))
	}

	r.groups[groupName].Use(ms...)

	return r
}
