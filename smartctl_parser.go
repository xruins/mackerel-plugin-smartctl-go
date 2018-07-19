package mackerel_plugin_smartctl_go

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type SmartctlParser struct {
	Devices        map[string]string
	SkipOnSpindown bool
}

const (
	SMARTCTL_CMD              = "smartctl -a %s"
	HDPARM_SPINDOWN_CHECK_CMD = "hdparm -C {%s}"
)

var regexpSmartctl = regexp.MustCompile(`\s*(\d+)\s+([A-Za-z_]+)\s+([x0-9]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+([A-Za-z_]+)\s+(\w+)\s+(.+)\s+(\d+)`)
var regexpHdparm = regexp.MustCompile(`([a-z\/]+):\n\s+drive state is:\s+(.+)`)

func (s *SmartctlParser) Get() {
	devices := make([]string, Len(s.Devices))
	for device, _ := range s.Devices {
		devices := append(devices, device)
	}
	devicesStatus := GetDevicesStatus(s.Devices)
}

func GetDevicesStatus(devices []string) (bool, err) {
	cmd := fmt.Sprintf(SMARTCTL_CMD, strings.Join(devices, ","))
	result, err := os.Exec(cmd)
	if err != nil {
		return nil, errors.New(err)
	}
	re := HDPARM_SPINDOWN_CHECK_CMD.Copy()
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
