package contracts

type Endpoint interface {
	MapEndpoint()
}

type RouteGroup interface{}

type RouteBuilder interface {
	RegisterRoutes(builder func(server any)) RouteBuilder
	RegisterGroup(groupName string) RouteBuilder
	RegisterGroupFunc(groupName string, builder func(group RouteGroup)) RouteBuilder
	UseForGroup(groupName string, middleware ...any) RouteBuilder
	Build() any
}
