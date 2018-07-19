package mackerel_plugin_smartctl_go

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/xruins/mackerel-plugin-smart-go"
)

type SmartctlParser struct {
	Devices        map[string]string
	SkipOnSpindown bool
}

const (
	SMARTCTL_CMD              = "smartctl -A %s %s"
	HDPARM_SPINDOWN_CHECK_CMD = "hdparm -C {%s}"
)

var regexpSmartctl = regexp.MustCompile(`\s*(\d+)\s+([A-Za-z_]+)\s+([x0-9]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+([A-Za-z_]+)\s+(\w+)\s+(.+)\s+(\d+)`)
var regexpHdparm = regexp.MustCompile(`([a-z\/]+):\n\s+drive state is:\s+(.+)`)

func (s *SmartctlParser) Get() {
	var devicesStatus map[string]string

	// set the statistics of device sleep state.
	// if it is configured to check SMART even if device sleeps,
	// filled standy statuses with false.
	if s.SkipOnSpindown {
		// fetch whether devices sleeps or not.
		devices := make([]string, Len(s.Devices))
		for device, _ := range s.Devices {
			devices := append(devices, device)
		}
		devicesStatus = GetDevicesStatus(s.Devices)
	} else {

		for device, standby := range devicesStatus {
			devicesStatus[device] = false
		}
	}
}

func GetSmartMetrics(devicesMap map[string]string) []*SmartctlMetric {
	var metrics = make(*SmartctlMetric, Len(devicesMap))
	for device, dmiType := range devicesMap {
		var dmiTypeOption string
		if dmiType != nil {
			dmiTypeOption = fmt.Sprintf("-d %s", dmiType)
		} else {
			dmiTypeOption = ""
		}
		cmd := fmt.Sprintf(SMARTCTL_CMD, device, dmiTypeOption)
		result, err := os.Exec(cmd)
		metrics = append(metrics, parseSmartctl)
	}
}

func parseSmartctl(s string) (*SmartctlMetric, err) {
	re := regexpSmartctl.Copy()
	matches := re.FindAll(s, -1)
	for matches
}
func GetDevicesStatus(devices []string) (bool, err) {
	cmd := fmt.Sprintf(HDPARM_SPINDOWN_CHECK_CMD, strings.Join(devices, ","))
	result, err := os.Exec(cmd)
	if err != nil {
		return nil, errors.New(err)
	}
	re := regexpHdparm.Copy()
	// output of hdparm is such as`/dev/sda:
	// drive state is: standby`.
	// Then the first subsequence will be path to device,
	// the second one will be status of disk.
	matches := re.FindAll(result, -1)
	if matches == nil {
		return nil, fmt.Errorf("The output of `%s` does not match regexp.", cmd)
	}
	var res map[string]bool
	for _, submatches := range matches {
		device := submatches[1]
		if matches[2] == "standby" {
			standby := true
		} else {
			standby := false
		}
		res[device] = standby
	}

	return res, nil
}
