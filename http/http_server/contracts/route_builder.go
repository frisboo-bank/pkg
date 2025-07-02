package contracts

type RouteBuilder interface {
	RegisterRoutes(builder func(server any)) RouteBuilder
	RegisterGroup(groupName string) RouteBuilder
	RegisterGroupFunc(groupName string, builder func(group any)) RouteBuilder
	UseForGroup(groupName string, middleware ...any) RouteBuilder
	Build() any
}
