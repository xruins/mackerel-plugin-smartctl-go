package mackerel_plugin_smartctl_go

import (
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type SmartctlPlugin struct {
	Prefix string
}

type SmartctlMetric struct {
	Device             string
	Temperature        string
	FailedMetricsCount int
}

func (s SmartctlPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(s.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"": {
			Label: labelPrefix,
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "temperature", Label: "Temperature"},
				{Name: "failed-metrics-count", Label: "Failed Metrics Count"},
			},
		},
	}
}

func (s SmartctlPlugin) FetchMetrics() (map[string]float64, error) {
	ut, err := uptime.Get()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch  metrics: %s", err)
	}
	return map[string]float64{"seconds": ut}, nil
}
