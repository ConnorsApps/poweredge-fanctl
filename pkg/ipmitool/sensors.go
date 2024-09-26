package ipmitool

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Host     string
	User     string
	Password string
}

type Tool struct {
	config *Config
}

func New(c *Config) *Tool {
	return &Tool{config: c}
}

func (t *Tool) commandOptions(opts ...string) []string {
	return append(
		[]string{
			"-I", "lanplus",
			"-H", t.config.Host,
			"-U", t.config.User,
			"-P", t.config.Password,
		},
		opts...,
	)
}

type Sensors struct {
	AvgFanSpeedRPM        float64
	MaxTempCelsius        float64
	TempOk                bool
	ExhaustTempCelsius    float64
	ExhaustOk             bool
	PowerConsumptionWatts float64
	PowerConsumptionOk    bool
}

type reading struct {
	Sensor  string
	Reading string
	Units   string
	Status  string
}

func parseSensorData(readings []reading) (*Sensors, error) {
	var (
		fanCount      float64
		totalFanSpeed float64
		sensors       = &Sensors{
			TempOk: true,
		}
		err        error
		errs       = []error{}
		fanReading = regexp.MustCompile(`^Fan[0-9].*`)
	)

	for _, s := range readings {
		switch s.Sensor {
		case "Exhaust Temp":
			sensors.ExhaustTempCelsius, err = strconv.ParseFloat(s.Reading, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("error parsing exhaust temperature '%s': %w", s.Reading, err))
			}
			sensors.ExhaustOk = s.Status == "ok"
		case "Temp":
			temp, err := strconv.ParseFloat(s.Reading, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("error parsing temperature '%s': %w", s.Reading, err))
			} else if temp > sensors.MaxTempCelsius {
				sensors.MaxTempCelsius = temp
			}
			if s.Status != "ok" {
				sensors.TempOk = false
			}
		case "Pwr Consumption":
			sensors.PowerConsumptionWatts, err = strconv.ParseFloat(s.Reading, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("error parsing power consumption '%s': %w", s.Reading, err))
			}
			sensors.PowerConsumptionOk = s.Status == "ok"
		default:
			if strings.HasPrefix(s.Sensor, "Fan") && fanReading.MatchString(s.Sensor) {
				speed, err := strconv.ParseFloat(s.Reading, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("error parsing fan speed '%s': %w", s.Reading, err))
				} else {
					fanCount++
					totalFanSpeed += speed
				}
			}
		}
	}

	if sensors.MaxTempCelsius == 0 {
		errs = append(errs, fmt.Errorf("no temperature sensor found because max temperature is 0"))
	}

	sensors.AvgFanSpeedRPM = totalFanSpeed / fanCount

	return sensors, errors.Join(errs...)
}

// Sensors parses the output of 'ipmitool -c sensor'
func (t *Tool) Sensors() (*Sensors, error) {
	cmd := exec.Command("ipmitool", t.commandOptions("-c", "sensor")...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var (
		reader  = csv.NewReader(strings.NewReader(string(output)))
		sensors = make([]reading, 0, 40)
	)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("invalid sensor data: %v. Expected >= 4 columns", record)
		}

		sensors = append(sensors, reading{
			Sensor:  record[0],
			Reading: record[1],
			Units:   record[2],
			Status:  record[3],
		})
	}

	return parseSensorData(sensors)
}
