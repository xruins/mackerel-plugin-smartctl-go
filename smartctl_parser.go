package mackerel_plugin_smartctl_go

import (
	"fmt"
	"regexp"

	"github.com/go-cmd/cmd"

	_ "github.com/xruins/mackerel-plugin-smartctl-go"
)

type SmartctlParser struct {
	Devices   map[string]string
	PowerMode string
}

const (
	SMARTCTL_CMD              = "smartctl -A %s"
	HDPARM_SPINDOWN_CHECK_CMD = "hdparm -C {%s}"
	SLEEP_STATE_EXIT_CODE     = 3
)

var regexpSmartctl = regexp.MustCompile(`\s*(\d+)\s+([A-Za-z_]+)\s+([x0-9]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+([A-Za-z_]+)\s+(\w+)\s+(.+)\s+(\d+)`)
var regexpHdparm = regexp.MustCompile(`([a-z\/]+):\n\s+drive state is:\s+(.+)`)

func (s *SmartctlParser) GetSmartMetrics() []*SmartctlMetric {
	var metrics = make(*SmartctlMetric, Len(devicesMap))
	for device, dmiType := range s.Devices {
		cmdOptions := ""
		if dmiType != nil {
			cmdOptions = string.Join(cmdOptions, Sprintf("-d %s", dmiType))
		}
		if s.PowerMode != nil {
			// ",3" lets smartctl to return error code 3 if check is skipped
			cmdOptions = string.Join(cmdOptions, Sprintf("-n %s,%d", PowerMode, SLEEP_STATE_EXIT_CODE))
		}
		smCmd := cmd.NewCmd(SMARTCTL_CMD, device, cmdOptions)
		status := <-smCmd.Start()
		if status.Exit != 0 && status.Exit != SLEEP_STATE_EXIT_CODE {
			return nil, fmt.Errorf("Failed to execute smartctl with device %s.", device)
		}
		metrics = append(metrics, parseSmartctl(smCmd.Stdout))
	}
	return metrics
}

func parseSmartctl(device string) (*SmartctlMetric, error) {
	re := regexpSmartctl.Copy()
	matches := re.FindAll(s, -1)
	failedMetricsCount := 0
	var temperature int
	for _, submatches := range matches {
		// if "WHEN_FAILED" shows the value which is not "-", treat as error
		if submatches[9] != "-" {
			failedMetricsCount += 1
		}
		// if "ATTRIBUTE_NAME" includes "Temperature", treat it as HDD temperature
		if string.Contains(submatches[2], "Temperature") {
			temperature = submatches[10]
		}
	}
	return &SmartctlMetric{
		failedMetricsCount: failedMetricsCount,
		temperature:        temperature,
	}
}
