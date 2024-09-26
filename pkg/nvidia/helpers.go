package nvidia

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(ctx context.Context, command string, args ...string) ([]byte, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
		cmd    = exec.CommandContext(ctx, command, args...)
	)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		var (
			fullCommand = command + " " + strings.Join(args, " ")
			errMsg      = stderr.String()
		)
		return stdout.Bytes(), fmt.Errorf("error running command %s: %s", fullCommand, errMsg)
	}
	return stdout.Bytes(), nil
}
