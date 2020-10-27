package sensors

import (
	"errors"
	"github.com/A-ndrey/raspi-manage-bot/db"
	"github.com/zack-wang/go-ccs811"
	"sync"
	"time"
)

const (
	dev  = "/dev/i2c-1"
	addr = 0x5a

	delay = 5 * time.Minute

	CO2Unit  = "CO2"
	TVOCUnit = "TVOC"
)

var isAvailable = false
var mutex sync.Mutex

func init() {
	mutex.Lock()
	defer mutex.Unlock()
	isAvailable = ccs811.Begin(dev, addr)
	if !isAvailable {
		go retry()
	}
}

func GetAirQuality() ([]db.Measurement, error) {
	mutex.Lock()
	eco2, tvoc, isValid := ccs811.ReadData(dev, addr)
	mutex.Unlock()
	if !isValid {
		return nil, errors.New("can't read air quality")
	}

	eco2Measurement := db.Measurement{
		Unit:        CO2Unit,
		Value:       float64(eco2),
		MeasureUnit: "ppm",
		Timestamp:   time.Now(),
	}

	tvocMeasurement := db.Measurement{
		Unit:        TVOCUnit,
		Value:       float64(tvoc),
		MeasureUnit: "ppm",
		Timestamp:   time.Now(),
	}

	return []db.Measurement{eco2Measurement, tvocMeasurement}, nil
}

func retry() {
	for {
		time.Sleep(delay)
		mutex.Lock()
		isAvailable = ccs811.Begin(dev, addr)
		mutex.Unlock()
		if isAvailable {
			return
		}
	}
}
