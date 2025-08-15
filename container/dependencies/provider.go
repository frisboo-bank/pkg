package dependencies

import (
	"errors"
	"strings"

	"frisboo-bank/pkg/options"
)

type Provider interface {
	Dependency
	Constructor() any
	Options() *options.OptionBuilder[ProviderOptions]
}

var _ Provider = (*provider)(nil)

type ProviderOptions struct {
	Name  string
	Group string
}

func ProvideWithOptions() *options.OptionBuilder[ProviderOptions] {
	return options.Apply(&ProviderOptions{})
}

func Name(name string) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		name = strings.TrimSpace(name)

		if name == "" {
			return errors.New("provider options: name can't be empty")
		}

		opts.Name = name
		return nil
	})
}

func Group(group string) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		group = strings.TrimSpace(group)

		if group == "" {
			return errors.New("provider options: group can't be empty")
		}

		opts.Group = group
		return nil
	})
}

type provider struct {
	constructor any
	options     *options.OptionBuilder[ProviderOptions]
}

func Provide(constructor any, opts ...*options.OptionBuilder[ProviderOptions]) Provider {
	var opt *options.OptionBuilder[ProviderOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &provider{constructor, opt}
}

func (i *provider) Constructor() any { return i.constructor }

func (i *provider) Options() *options.OptionBuilder[ProviderOptions] { return i.options }

func (i *provider) IsDependency() {}
