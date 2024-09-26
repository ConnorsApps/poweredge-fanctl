package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/ConnorsApps/poweredge-fanctl/fanctl/internal/config"
	"github.com/ConnorsApps/poweredge-fanctl/fanctl/internal/fan_speed"
	"github.com/ConnorsApps/poweredge-fanctl/pkg/ipmitool"
	"github.com/ConnorsApps/poweredge-fanctl/pkg/nvidia"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SpeedForTempature interface {
	FanSpeedForTemp(temp float64) float64
}

func main() {
	var (
		ctx   = context.Background()
		c     = config.MustRead(os.Getenv("CONFIG_PATH"))
		ipmit = ipmitool.New(c.IDRAC)

		lastSetFanSpeed float64
		fanSpeedForCPU  SpeedForTempature = fan_speed.NewCPU()
		fanSpeedForGPU  SpeedForTempature = fan_speed.NewGPU()
	)

	nvidiaGPUs, err := nvidia.New(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error with checking NVIDIA GPUs")
	}

	if err := ipmit.SetFanControlToManualMode(); err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to set fan control into manual mode. Make sure the iDRAC version is less then 3.30.30.30")
	}

	for {
		gpuStats, statErr := nvidiaGPUs.FetchStats(ctx)
		if statErr != nil {
			log.Error().Err(statErr).Msg("Error getting NVIDIA GPUs")
		}

		idracStats, sensorsErr := ipmit.Sensors()
		if sensorsErr != nil {
			log.Error().Err(sensorsErr).Msg("Error getting sensor data")
		}
		if statErr != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		// Pick whichever wants higher fan speeds
		var (
			gpuTemp      = gpuStats.MaxTemp()
			cpuSuggested = fanSpeedForCPU.FanSpeedForTemp(idracStats.MaxTempCelsius)
			gpuSuggested = fanSpeedForGPU.FanSpeedForTemp(gpuTemp)
			newFanSpeed  = math.Max(cpuSuggested, gpuSuggested)
		)

		for _, gpu := range gpuStats {
			log.Debug().
				Str("gpu", gpu.Name).
				Str("temp", fmt.Sprintf("%.0fC", gpu.TemperatureCelsius)).
				Msg("GPU Stats")
		}

		if !idracStats.ExhaustOk {
			log.Warn().
				Str("exhaustTemp", fmt.Sprintf("%.0fC", idracStats.ExhaustTempCelsius)).
				Msg("Exhaust Temp is not ok")
		}

		if !idracStats.PowerConsumptionOk {
			log.Warn().
				Str("power", fmt.Sprintf("%.0fW", idracStats.PowerConsumptionWatts)).
				Msg("Power Consumption is not ok")
		}

		if !idracStats.TempOk {
			log.Warn().Msg("Tempature is not ok")
		}

		log.Debug().
			Str("exhaustTemp", fmt.Sprintf("%.0fC", idracStats.ExhaustTempCelsius)).
			Str("temp", fmt.Sprintf("%.0fC", idracStats.MaxTempCelsius)).
			Msg("iDRAC Stats")

		log.Debug().
			Str("newFanSpeed", fmt.Sprintf("%.0f", newFanSpeed)+"%").
			Str("cpuSuggested", fmt.Sprintf("%.0f", cpuSuggested)+"%").
			Str("gpuSuggested", fmt.Sprintf("%.0f", gpuSuggested)+"%").
			Msg("Speed Suggested")

		percentChange := math.Abs((newFanSpeed - lastSetFanSpeed) / lastSetFanSpeed * 100)

		// Don't set a fan speed if the speed hasn't changed
		// by more than 3 percent
		if percentChange > 3 {
			log.Info().
				Str("temp", fmt.Sprintf("%.0fC", idracStats.MaxTempCelsius)).
				Msgf("New Fan speed %.0f%%", newFanSpeed)

			if err := ipmit.SetFanSpeed(uint8(newFanSpeed)); err != nil {
				log.Error().Err(err).Msg("Unable to set fan speed")
			} else {
				lastSetFanSpeed = newFanSpeed
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func init() {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.SetLogLevel(level)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
