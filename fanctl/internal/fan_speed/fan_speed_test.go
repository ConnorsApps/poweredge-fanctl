package fan_speed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFanSpeedForTemp(t *testing.T) {
	cpu := NewCPU()
	a := assert.New(t)
	cases := []struct {
		Temp          float64
		ExpectedSpeed float64
	}{
		{Temp: 0, ExpectedSpeed: 15},
		{Temp: 30, ExpectedSpeed: 15},
		{Temp: 31, ExpectedSpeed: 17.5},
		{Temp: 35, ExpectedSpeed: 27.5},
		{Temp: 40, ExpectedSpeed: 40},
	}

	for _, c := range cases {
		a.Equalf(c.ExpectedSpeed, cpu.FanSpeedForTemp(c.Temp), "temp %.1f", c.Temp)
	}
}
