package mackerel_plugin_smartctl_go

import (
	"flag"
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
				{Name: "temperature", Label: "Temperature", Stacked: false},
				{Name: "failed-metrics-count", Label: "Failed Metrics Count", Stacked: false},
			},
		},
	}
}

func (s SmartctlPlugin) FetchMetrics() (map[string]float64, error) {
	sm, err := SmartctlParser.Get()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch metrics: %s", err)
	}
	return map[string]float64{"seconds": ut}, nil
}

const (
	POWERMODE_DESCRIPTION = `
standby - [DEFAULT] check the disks are spinning. The disks are not spinning will be ignored.
sleep - check the disks is not in SLEEP mode.
never - always check disks.
`
	DMI_TYPE_DESCRIPTION = `
Specify DMI type if you want to let smartctl use non "ata" DMI type. (see -d option of "man smartctl" for detail)
If you execute with "mackerel-plugin-smartctl-go /dev/sda /dev/sdb --dmi-type ,nvme,scsi",
the commands "smartctl -A /dev/sda" and "smartctl -A /dev/sdb -d nvme", "smartctl -A /dev/sdc -d scsi" will be executed.
`
)

const (
	standby = iota
	sleep
	never
)

func main() {
	optPrefix := flag.String("metric-key-prefix", "uptime", "Metric key prefix")
	optPowerMode := flag.String("power-mode", "standby", POWERMODE_DESCRIPTION)
	optDmiType := flag.String("dmi-type", "", DMI_TYPE_DESCRIPTION)
	flag.Parse()

	u := UptimePlugin{
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
