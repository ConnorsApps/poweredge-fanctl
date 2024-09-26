package nvidia

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type GPU struct {
	Index string
	Name  string
}

func listGPUs(ctx context.Context) ([]*GPU, error) {
	output, err := runCmd(ctx, "nvidia-smi", "-L")
	if err != nil {
		return nil, err
	}

	var (
		gpuIndexRegex = regexp.MustCompile(`GPU \b([0-9]|[1-9][0-9]|100)\b:`)
		gpuNameRegex  = regexp.MustCompile(`: .* \(UUID`)
	)
	// Example output:
	// GPU 0: GeForce GTX 1080 Ti (UUID: GPU-3d8b3e4c-2e0c-5c3e-1b8c-5b8b6b7e2c7f)
	// GPU 1: GeForce GTX 1080 Ti (UUID: GPU-3d8b3e4c-2e0c-5c3e-1b8c-5b8b6b7e2c7f)
	var gpus []*GPU
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		var (
			// : GeForce GTX 1080 Ti (UUID
			gpuName = strings.TrimSuffix(strings.TrimPrefix(gpuNameRegex.FindString(line), ": "), " (UUID")
			// GPU 0:
			gpuIndex = strings.TrimSuffix(strings.TrimPrefix(gpuIndexRegex.FindString(line), "GPU "), ":")
		)

		if len(gpuName) == 0 || len(gpuIndex) == 0 {
			return gpus, fmt.Errorf("unexpected result from parsing nvidia-smi -L: %s", line)
		}

		gpus = append(gpus, &GPU{
			Index: gpuIndex,
			Name:  gpuName,
		})
	}
	return gpus, nil
}
