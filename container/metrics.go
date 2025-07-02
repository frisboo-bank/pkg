package container

import (
	"fmt"
	"strings"
	"sync"
)

type metricType = string

const (
	metricTypeModule    = metricType("module")
	metricTypeProvider  = metricType("provider")
	metricTypeHook      = metricType("hook")
	metricTypeDecorator = metricType("decorator")
	metricTypeInvoke    = metricType("invoke")
)

type moduleMetrics struct {
	modules    int
	providers  int
	hooks      int
	decorators int
	invokes    int
}

type ContainerMetrics struct {
	modules map[string]*moduleMetrics
	mu      sync.RWMutex
}

func NewDigMetrics() *ContainerMetrics {
	return &ContainerMetrics{
		modules: make(map[string]*moduleMetrics),
	}
}

func (m *ContainerMetrics) IncrementModules(moduleName string, count int) {
	m.increment(moduleName, metricTypeModule, count)
}

func (m *ContainerMetrics) IncrementProviders(moduleName string, count int) {
	m.increment(moduleName, metricTypeProvider, count)
}

func (m *ContainerMetrics) IncrementHooks(moduleName string, count int) {
	m.increment(moduleName, metricTypeHook, count)
}

func (m *ContainerMetrics) IncrementDecorators(moduleName string, count int) {
	m.increment(moduleName, metricTypeDecorator, count)
}

func (m *ContainerMetrics) IncrementInvokes(moduleName string, count int) {
	m.increment(moduleName, metricTypeInvoke, count)
}

func (m *ContainerMetrics) ToString() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.modules) <= 0 {
		return "no modules registered"
	}

	var buf strings.Builder
	for name, metrics := range m.modules {
		buf.WriteString(fmt.Sprintf("  • %s: ", name))
		buf.WriteString(fmt.Sprintf("sub-modules=%d, ", metrics.modules))
		buf.WriteString(fmt.Sprintf("providers=%d, ", metrics.providers))
		buf.WriteString(fmt.Sprintf("hooks=%d, ", metrics.hooks))
		buf.WriteString(fmt.Sprintf("decorators=%d, ", metrics.decorators))
		buf.WriteString(fmt.Sprintf("invokes=%d\n", metrics.invokes))
	}

	return buf.String()
}

func (m *ContainerMetrics) increment(moduleName string, metricType metricType, count int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.modules[moduleName]; !ok {
		m.modules[moduleName] = &moduleMetrics{}
	}
	metrics := m.modules[moduleName]

	switch metricType {
	case metricTypeModule:
		metrics.modules += count
	case metricTypeProvider:
		metrics.providers += count
	case metricTypeHook:
		metrics.hooks += count
	case metricTypeDecorator:
		metrics.decorators += count
	case metricTypeInvoke:
		metrics.invokes += count
	}
}
