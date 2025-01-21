package fan_speed

// Set fan control in manual mode
// If errors are occuring, set fans to 30% just to be safe
// Make logic interface

// MaxTemp <= 30 set fans to 15%
// MaxTemp >= 40 set fans to 40%
// MaxTemp >= 45 set fans to 50%
// MaxTemp >= 50 set fans to 70%

// GPU temp <= 40 set fans to 15%
// GPU temp >= 45 set fans to 40%
// GPU temp >= 50 set fans to 70%

// Make a smooth linar transiton between two points
func linearIncraseInPercent(baseTemp, maxTemp, baseSpeed, maxSpeed, currentTemp float64) float64 {
	var (
		changeInTemp = currentTemp - baseTemp
		percentDiff  = maxSpeed - baseSpeed
		tempDiff     = maxTemp - baseTemp
	)

	return (changeInTemp / tempDiff) * (percentDiff)
}

func tempSpeedRamp(minSpeed, maxSpeed float64, tempPerSpeed [][]float64, temp float64) float64 {
	for i, tempSpeed := range tempPerSpeed {
		if temp <= tempSpeed[0] {
			if i == 0 {
				return minSpeed
			}
			return minSpeed + linearIncraseInPercent(tempPerSpeed[i-1][0], tempSpeed[0], minSpeed, tempSpeed[1], temp)
		}
	}
	return maxSpeed
}

type CPU struct {
	minSpeed     float64
	maxSpeed     float64
	tempPerSpeed [][]float64
}

func NewCPU() *CPU {
	return &CPU{
		minSpeed: 10,
		maxSpeed: 80,
		tempPerSpeed: [][]float64{
			{45, 15},
			{50, 30},
			{55, 80},
		},
	}
}

func (c *CPU) FanSpeedForTemp(temp float64) float64 {
	return tempSpeedRamp(c.minSpeed, c.maxSpeed, c.tempPerSpeed, temp)
}

type GPU struct {
	minSpeed     float64
	maxSpeed     float64
	tempPerSpeed [][]float64
}

func NewGPU() *GPU {
	return &GPU{
		minSpeed: 10,
		maxSpeed: 80,
		tempPerSpeed: [][]float64{
			{45, 15},
			{50, 40},
			{55, 70},
			{60, 100},
		},
	}
}

func (g *GPU) FanSpeedForTemp(temp float64) float64 {
	return tempSpeedRamp(g.minSpeed, g.maxSpeed, g.tempPerSpeed, temp)
}
