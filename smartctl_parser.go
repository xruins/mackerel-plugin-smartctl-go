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
	for device, device_type := range s {
		if isSpindown(device) {
			continue
		}
		parseSmartctl(result)
	}
}

func isSpindown(devices []string) (bool, err) {
	cmd := fmt.Sprintf(SMARTCTL_CMD, strings.Join(devices, ","))
	result, err := os.Exec(cmd)
	if err != nil {
		return nil, errors.New(err)
	}
	re := HDPARM_SPINDOWN_CHECK_CMD.Copy()
	isSleep := re.Find(result)
	if isSleep != "standby" {
		return true
	} else {
		return false
	}
}
