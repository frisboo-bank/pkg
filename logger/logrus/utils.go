package logrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = (*prefixHook)(nil)

type prefixHook struct {
	prefix string
}

func (p *prefixHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] %s", p.prefix, entry.Message)
	return nil
}

func (p *prefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
