package contracts

type RouterEngine interface {
	Group(path string, middlewares ...any) RouterEngine
	Handle(method string, path string, handler any, middlewares ...any)
	Static(prefix string, root string)
}

type RouteBuilder interface {
	Root() RouteGroup
}

type RouteGroup interface {
	// Nested group
	Group(path string, middlewares ...any) RouteGroup

	DELETE(path string, handler any, middlewares ...any) RouteGroup
	GET(path string, handler any, middlewares ...any) RouteGroup
	HEAD(path string, handler any, middlewares ...any) RouteGroup
	OPTIONS(path string, handler any, middlewares ...any) RouteGroup
	PATCH(path string, handler any, middlewares ...any) RouteGroup
	POST(path string, handler any, middlewares ...any) RouteGroup
	PUT(path string, handler any, middlewares ...any) RouteGroup
	TRACE(path string, handler any, middlewares ...any) RouteGroup

	Handle(method string, path string, handler any, middlewares ...any) RouteGroup
	Static(prefix string, root string) RouteGroup
}
