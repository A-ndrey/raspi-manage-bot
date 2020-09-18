package board

import (
	"bytes"
	"fmt"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func GetTemperature() string {
	var result []string

	cpuMeasurement, err := GetCPUTemp()
	if err != nil {
		log.Println(err)
	} else {
		result = append(result, cpuMeasurement.String())
	}

	gpuMeasurement, err := GetGPUTemp()
	if err != nil {
		log.Println(err)
	} else {
		result = append(result, gpuMeasurement.String())
	}

	return strings.Join(result, "\n")
}

func GetGPUTemp() (db.Measurement, error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	var btsOut, btsErr bytes.Buffer
	cmd.Stdout = &btsOut
	cmd.Stderr = &btsErr
	err := cmd.Run()
	if err != nil {
		return db.Measurement{}, fmt.Errorf("can't get gpu temp: %s: %w", btsErr.String(), err)
	}

	submatch := gpuRegexp.FindStringSubmatch(btsOut.String())

	tempVal, err := strconv.ParseFloat(strings.TrimSpace(submatch[1]), 64)
	if err != nil {
		return db.Measurement{}, fmt.Errorf("can't parse gpu temp: %w", err)
	}

	gpuMeasurement := db.Measurement{
		Unit:        GPU_UNIT,
		Value:       tempVal,
		MeasureUnit: submatch[2],
		Timestamp:   time.Now(),
	}

	return gpuMeasurement, nil
}

func GetCPUTemp() (db.Measurement, error) {
	btsTemp, err := ioutil.ReadFile(cpuTempFilePath)
	if err != nil {
		return db.Measurement{}, fmt.Errorf("can't read %s: %w", cpuTempFilePath, err)
	}

	tempVal, err := strconv.ParseFloat(strings.TrimSpace(string(btsTemp)), 64)
	if err != nil {
		return db.Measurement{}, fmt.Errorf("can't parse cpu temp: %w", err)
	}

	cpuMeasurement := db.Measurement{
		Unit:        CPU_UNIT,
		Value:       tempVal / 1000,
		MeasureUnit: "C",
		Timestamp:   time.Now(),
	}

	return cpuMeasurement, nil
}
