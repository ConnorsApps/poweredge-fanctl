package nvidia

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Stat struct {
	TemperatureCelsius float64
	PowerDrawWatts     float64
}

func getStats(ctx context.Context) (map[string]*Stat, error) {
	output, err := runCmd(ctx, "nvidia-smi", "stats", "-d", "pwrDraw,temp", "--count", "1")
	if err != nil {
		return nil, err
	}

	// Example output:
	// 0, pwrDraw , 1726958883095356, 6
	// 0, temp    , 1726958883696592, 31
	var devices = make(map[string]*Stat, 2)
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ", ")
		if len(fields) != 4 {
			return devices, fmt.Errorf("unexpected number of fields from nvidia-smi stats: %s", line)
		}
		var (
			device = strings.TrimSpace(fields[0])
			name   = strings.TrimSpace(fields[1])
		)
		value, err := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
		if err != nil {
			return devices, fmt.Errorf("error parsing number value from nvidia-smi stats: %s", fields[3])
		}
		if _, ok := devices[device]; !ok {
			devices[device] = &Stat{}
		}
		switch name {
		case "pwrDraw":
			devices[device].PowerDrawWatts = value
		case "temp":
			devices[device].TemperatureCelsius = value
		}
	}
	return devices, nil
}
