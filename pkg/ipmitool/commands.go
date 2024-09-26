package ipmitool

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func setFanSpeedCommand(speedPercent uint8) []string {
	return []string{
		"raw",
		"0x30",
		"0x30",
		"0x02",
		"0xff",
		fmt.Sprintf("0x%x", speedPercent),
	}
}

func (t *Tool) SetFanSpeed(speedPercent uint8) error {
	output, err := exec.Command(
		"ipmitool",
		t.commandOptions(
			setFanSpeedCommand(speedPercent)...,
		)...,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to set fan speed: %s err=%w", string(output), err)
	}
	return nil
}

func (t *Tool) setFanControlToAutoMode() ([]byte, error) {
	return exec.Command(
		"ipmitool",
		t.commandOptions("raw", "0x30", "0x30", "0x01", "0x00")...,
	).CombinedOutput()
}

func (t *Tool) SetFanControlToManualMode() error {
	output, err := t.setFanControlToAutoMode()
	if err == nil {
		return nil
	}

	if !strings.Contains(err.Error(), "insufficient resources for session") {
		return fmt.Errorf("unable to set fan speed: %s err=%w", string(output), err)
	}

	// Wait before retrying
	time.Sleep(5 * time.Second)
	output, err = t.setFanControlToAutoMode()
	if err != nil {
		return fmt.Errorf("unable to set fan speed: %s err=%w", string(output), err)
	}

	return nil
}
