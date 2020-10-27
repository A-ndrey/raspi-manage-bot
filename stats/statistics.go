package stats

import (
	"github.com/A-ndrey/raspi-manage-bot/db"
	"sort"
)

func Median(measurements []db.Measurement) db.Measurement {
	size := len(measurements)
	if size == 0 {
		return db.Measurement{}
	}

	sort.Slice(measurements, func(i, j int) bool {
		return measurements[i].Value < measurements[j].Value
	})

	if size%2 == 0 {
		return Mean(measurements[size/2-1 : size/2+1])
	} else {
		return measurements[size/2]
	}
}

func Mean(measurements []db.Measurement) db.Measurement {
	size := len(measurements)
	if size == 0 {
		return db.Measurement{}
	}

	result := measurements[0]
	for i := 1; i < size; i++ {
		result.Value += measurements[i].Value
	}

	result.Value /= float64(size)

	return result
}

func GroupByUnit(measurements []db.Measurement) map[string][]db.Measurement {
	groups := make(map[string][]db.Measurement)
	for _, m := range measurements {
		groups[m.Unit] = append(groups[m.Unit], m)
	}

	return groups
}
