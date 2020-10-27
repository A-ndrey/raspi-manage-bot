package stats

import (
	"context"
	"fmt"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"log"
	"time"
)

const analyzingPeriod = time.Minute

var lastAnalyzedResult = make(map[string]db.Measurement)

func StartAnalyzing(ctx context.Context, config configs.Config, monitoring chan<- db.Measurement) {
	log.Println("Start analyzing")
	defer log.Println("End analyzing")

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(analyzingPeriod):
			measurements, err := db.GetMeasurementsInTimeInterval(time.Now().Add(-analyzingPeriod), time.Now())
			if err != nil {
				log.Println(err)
			}
			for k, v := range GroupByUnit(measurements) {
				measurement := Median(v)
				lastAnalyzedResult[k] = measurement
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

func GetMeasurementByUnit(unit string) (db.Measurement, error) {
	measurement, ok := lastAnalyzedResult[unit]
	if !ok {
		return db.Measurement{}, fmt.Errorf("measurement %s still unavailable", unit)
	}

	return measurement, nil
}
