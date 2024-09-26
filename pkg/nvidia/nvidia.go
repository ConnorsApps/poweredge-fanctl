package nvidia

import (
	"context"
	"fmt"
)

type Nvidia struct {
	gpus map[string]*GPU
}

func (n *Nvidia) updateGPUs(ctx context.Context) error {
	gpus, err := listGPUs(ctx)
	if err != nil {
		return err
	}

	n.gpus = make(map[string]*GPU, len(gpus))
	for _, gpu := range gpus {
		n.gpus[gpu.Index] = gpu
	}
	return nil
}

func New(ctx context.Context) (*Nvidia, error) {
	n := &Nvidia{}
	if err := n.updateGPUs(ctx); err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Nvidia) GetGPUs() []*GPU {
	var (
		gpus = make([]*GPU, 0, len(n.gpus))
		i    = 0
	)
	for _, gpu := range n.gpus {
		gpus[i] = gpu
		i++
	}
	return gpus
}

type GPUStat struct {
	*Stat
	*GPU
}

type Stats []GPUStat

func (stats Stats) MaxTemp() float64 {
	var max float64
	for _, stat := range stats {
		if stat.TemperatureCelsius > max {
			max = stat.TemperatureCelsius
		}
	}
	return max
}

func (n *Nvidia) FetchStats(ctx context.Context) (Stats, error) {
	stats, err := getStats(ctx)
	if err != nil {
		return nil, err
	}
	var (
		i        = 0
		gpuStats = make([]GPUStat, len(stats))
	)

	for index, stat := range stats {
		if n.gpus[index] == nil {
			// If we don't have record of this GPU then update the list
			if err := n.updateGPUs(ctx); err != nil {
				return nil, err
			}
			// If we still don't have record of this GPU then return an error
			if n.gpus[index] == nil {
				return nil, fmt.Errorf("GPU with index %s not found from nvidia-smi -L", index)
			}
		}
		gpuStats[i] = GPUStat{
			Stat: stat,
			GPU:  n.gpus[index],
		}
		i++
	}
	return gpuStats, nil
}
