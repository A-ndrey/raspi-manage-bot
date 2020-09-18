package board

import (
	"context"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"log"
	"regexp"
	"time"
)

const (
	cpuTempFilePath = "/sys/class/thermal/thermal_zone0/temp"

	measuringPeriod = time.Minute

	GPU_UNIT = "GPU"
	CPU_UNIT = "CPU"
)

var gpuRegexp = regexp.MustCompile(`temp=(\d+\.\d+)'(.+)`)

func StartMeasuring(ctx context.Context, config configs.Config, monitoring chan<- db.Measurement) {
	log.Println("Start board measuring")
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(measuringPeriod):
			var measurements []db.Measurement

			if gpuMeasurement, err := GetGPUTemp(); err != nil {
				log.Println(err)
			} else {
				measurements = append(measurements, gpuMeasurement)
			}

			if cpuMeasurement, err := GetCPUTemp(); err != nil {
				log.Println(err)
			} else {
				measurements = append(measurements, cpuMeasurement)
			}

			for _, measurement := range measurements {
				if err := db.InsertMeasurement(measurement); err != nil {
					log.Println(err)
				}
				if isWarn(config.Monitoring, measurement) {
					monitoring <- measurement
				}
			}
		}
	}
}

func isWarn(limits []db.Measurement, actual db.Measurement) bool {
	for _, limit := range limits {
		if limit.Unit == actual.Unit && limit.MeasureUnit == actual.MeasureUnit && limit.Value <= actual.Value {
			return true
		}
	}

	return false
}
