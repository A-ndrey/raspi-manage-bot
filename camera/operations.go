package camera

import (
	"bytes"
	"fmt"
	"os/exec"
)

func TakePhoto() ([]byte, error) {
	cmd := exec.Command("raspistill", "-t", "2", "-o", "-")

	var bts bytes.Buffer
	var btsErr bytes.Buffer

	cmd.Stderr = &btsErr
	cmd.Stdout = &bts

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("can't take photo: %s", btsErr.String())
	}

	return bts.Bytes(), nil
}
