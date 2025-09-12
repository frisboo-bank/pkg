package container

type constructorWithParams struct {
	fn any
}

func (c constructorWithParams) underlying() any { return c.fn }

func ConstructorWithParams(fn any) any {
	return constructorWithParams{fn: fn}
}
