package options

type Services map[string]string

func (s *Services) Decode(value string) error {
	return nil
}
