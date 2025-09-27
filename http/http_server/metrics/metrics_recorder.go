package metrics

var _ MetricsRecorder = (*metricsRecorder)(nil)

type MetricsRecorder any

type metricsRecorder struct{}

func NewMetricsRecorder() MetricsRecorder {
	return &metricsRecorder{}
}
