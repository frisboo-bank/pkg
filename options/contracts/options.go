package contracts

type Options interface {
	SetDefaults()
	Validate() error
	Clone() Options
}
