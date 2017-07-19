package service

// PerformanceIndex the application's performance
type PerformanceIndex struct {
	FailRatio     float64 `json:"failRatio"`
	FailAmount    int64   `json:"failAmount"`
	AccessAmount  int64   `json:"accessAmount"`
	MaxConcurrent int64   `json:"maxConcurrent"`
	MinLatency    int64   `json:"minLatency"`
	AvgLatency    int64   `json:"avgLatency"`
}

// Metric to monitor
type Metric struct {
	PerformanceIndex `json:"performance_index"`
	Host             string `json:"host"`
	AppName          string `json:"app_name"`
	Domain           string `json:"domain"`
}

func (m *Metric) copyMetric(metric Metric) {
	m.Host = metric.Host
	m.AppName = metric.AppName
	m.Domain = metric.Domain
	m.FailAmount = metric.FailAmount
	m.FailRatio = metric.FailRatio
	m.AccessAmount = metric.AccessAmount
	m.MaxConcurrent = metric.MaxConcurrent
	m.MinLatency = metric.MinLatency
	m.AvgLatency = metric.AvgLatency
}

type metricRepository interface {
	newMetric() (err error)
	getMetric() (metric Metric, err error)
	getAppMetric(appname string) (metric Metric, err error)
}

type metricRequest struct {
	AppName string `json:"appname"`
}

func (request metricRequest) isValid() (valid bool) {
	valid = true
	if request.AppName == "" {
		valid = false
	}
	return valid
}
