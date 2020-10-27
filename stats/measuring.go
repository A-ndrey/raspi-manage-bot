package stats

import (
	"context"
	"github.com/A-ndrey/raspi-manage-bot/board"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"github.com/A-ndrey/raspi-manage-bot/sensors"
	"log"
	"time"
)

const measuringPeriod = 10 * time.Second

func StartMeasuring(ctx context.Context) {
	log.Println("Start measuring")
	defer log.Println("End measuring")

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(measuringPeriod):
			var measurements []db.Measurement

			if gpuMeasurement, err := board.GetGPUTemp(); err != nil {
				log.Println(err)
			} else {
				measurements = append(measurements, gpuMeasurement)
			}

			if cpuMeasurement, err := board.GetCPUTemp(); err != nil {
				log.Println(err)
			} else {
				measurements = append(measurements, cpuMeasurement)
			}

			if airQualityMeasurement, err := sensors.GetAirQuality(); err != nil {
				log.Println(err)
			} else {
				measurements = append(measurements, airQualityMeasurement...)
			}

			for _, measurement := range measurements {
				if err := db.InsertMeasurement(measurement); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
