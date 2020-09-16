package board

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Reboot() error {
	cmd := exec.Command("/bin/bash", "-c", "sudo reboot")
	var bts bytes.Buffer
	cmd.Stderr = &bts
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't reboot: %s: %w", bts.String(), err)
	}

	return nil
}
