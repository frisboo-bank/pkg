package echo

import (
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/syserrors"

	echoVendor "github.com/labstack/echo/v4"
)

type routeBuilder struct {
	echo   *echoVendor.Echo
	groups map[string]*echoVendor.Group
}

var _ contracts.RouteBuilder = (*routeBuilder)(nil)

func NewRouteBuilder(echo *echoVendor.Echo) contracts.RouteBuilder {
	return &routeBuilder{echo: echo}
}

func (r *routeBuilder) Build() any {
	return r.echo
}

func (r *routeBuilder) RegisterGroup(groupName string) contracts.RouteBuilder {
	if r.groups == nil {
		r.groups = make(map[string]*echoVendor.Group)
	}
	if _, exist := r.groups[groupName]; !exist {
		r.groups[groupName] = r.echo.Group(groupName)
	}
	return r
}

func (r *routeBuilder) RegisterGroupFunc(
	groupName string,
	builder func(group contracts.RouteGroup),
) contracts.RouteBuilder {
	r.RegisterGroup(groupName)
	builder(r.groups[groupName])
	return r
}

func (r *routeBuilder) RegisterRoutes(builder func(server any)) contracts.RouteBuilder {
	builder(r.echo)
	return r
}

func (r *routeBuilder) UseForGroup(groupName string, middlewares ...any) contracts.RouteBuilder {
	r.RegisterGroup(groupName)
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Newf("invalid middleware for group: `%s` `%v`", groupName, err))
	}
	r.groups[groupName].Use(ms...)
	return r
}
