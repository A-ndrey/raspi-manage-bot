package board

import (
	"bytes"
	"context"
	"fmt"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	cpuTempFilePath = "/sys/class/thermal/thermal_zone0/temp"

	measuringPeriod = time.Minute

	GPU_UNIT = "GPU"
	CPU_UNIT = "CPU"
)

var gpuRegexp = regexp.MustCompile(`temp=(\d+\.\d+)'(.+)`)

func StartMeasuring(ctx context.Context) {
	log.Println("Start board measuring")
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(measuringPeriod):
			gpuMeasurement, err := getGPUTemp()
			if err != nil {
				log.Println(err)
			} else {
				err = db.InsertMeasurement(gpuMeasurement)
				if err != nil {
					log.Println(err)
				}
			}

			cpuMeasurement, err := getCPUTemp()
			if err != nil {
				log.Println(err)
			} else {
				err := db.InsertMeasurement(cpuMeasurement)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func GetTemperature() string {
	var result []string

	cpuMeasurement, err := getCPUTemp()
	if err != nil {
		log.Println(err)
	} else {
		result = append(result, cpuMeasurement.String())
	}

	gpuMeasurement, err := getGPUTemp()
	if err != nil {
		log.Println(err)
	} else {
		result = append(result, gpuMeasurement.String())
	}

	return strings.Join(result, "\n")
}

func getGPUTemp() (db.Measurement, error) {
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

func getCPUTemp() (db.Measurement, error) {
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
