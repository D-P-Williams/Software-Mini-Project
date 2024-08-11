//nolint:mnd // File does multiple mathematical operations, ignore magic numbers in this file.
package transporthandler

import (
	"math"
	"time"
	"work-mini-project/pkg/configuration"
	customerhandler "work-mini-project/pkg/customerHandler"
)

type TransportHandler struct {
	config *configuration.Config
}

type TripDetails struct {
	Method   string
	Duration time.Duration
	Cost     float64
	Distance float64
}

func New(config *configuration.Config) *TransportHandler {
	return &TransportHandler{
		config: config,
	}
}

func (th *TransportHandler) CalculateCosts(customer customerhandler.Customer) []*TripDetails {
	transportMethods := []*TripDetails{}

	transportMethods = append(transportMethods, th.calculateLorry(customer))
	transportMethods = append(transportMethods, th.calculateCanalBoat(customer))
	transportMethods = append(transportMethods, th.calculateHelicopter(customer))

	return transportMethods
}

//nolint:nonamedreturns // Named returns for clarity with same type
func (th *TransportHandler) calculateXYDistances(customer customerhandler.Customer) (x float64, y float64) {
	diffX := math.Abs(float64(th.config.Company.GridX) - float64(customer.GridX))
	diffY := math.Abs(float64(th.config.Company.GridY) - float64(customer.GridY))

	return diffX, diffY
}

func (th *TransportHandler) calculateDirectDistance(customer customerhandler.Customer) float64 {
	diffX := float64(th.config.Company.GridX) - float64(customer.GridX)
	diffY := float64(th.config.Company.GridY) - float64(customer.GridY)

	return math.Sqrt(math.Pow(diffX, 2) + math.Pow(diffY, 2))
}

func (th *TransportHandler) calculateLorry(customer customerhandler.Customer) *TripDetails {
	diffX, diffY := th.calculateXYDistances(customer)

	totalDist := diffX + diffY

	// Calc Time
	speed := float64(th.config.Vehicles.Lorry.Speed)
	totalTimeHr := totalDist / speed

	trafficStops := math.Floor(float64(totalDist) / float64(th.config.Vehicles.Lorry.TrafficDelayFrequency))
	totalTime := totalTimeHr + (trafficStops * (float64(th.config.Vehicles.Lorry.TrafficDelayTime) / 60.0))
	totalTimeDuration := time.Duration(totalTime * float64(time.Hour))

	// Calc Cost
	cost := (1.0 / 12.0) * (math.Pow(totalDist, 2) - float64(95*totalDist) + 2880)

	return &TripDetails{
		Method:   "Lorry",
		Duration: totalTimeDuration,
		Cost:     cost,
		Distance: diffX + diffY,
	}
}

func (th *TransportHandler) calculateCanalBoat(customer customerhandler.Customer) *TripDetails {
	diffX, diffY := th.calculateXYDistances(customer)

	totalDist := diffX + diffY

	// Calc Time
	speed := float64(th.config.Vehicles.CanalBoat.Speed)
	totalTimeHr := totalDist / speed
	totalTimeDuration := time.Duration(totalTimeHr * float64(time.Hour))

	// Calc Cost
	cost := ((5 * totalDist) / 12.0) + (1280.0 / 12.0)

	return &TripDetails{
		Method:   "Canal boat",
		Duration: totalTimeDuration,
		Cost:     cost,
		Distance: diffX + diffY,
	}
}

func (th *TransportHandler) calculateHelicopter(customer customerhandler.Customer) *TripDetails {
	totalDist := th.calculateDirectDistance(customer)

	// Calc Time
	speed := float64(th.config.Vehicles.Helicopter.Speed)
	totalTimeHr := totalDist / speed
	totalTimeDuration := time.Duration(
		(totalTimeHr + float64(th.config.Vehicles.Helicopter.InitialDelay)/60.0) *
			float64(time.Hour),
	)

	// Calc Cost
	cost := (0.5 * totalDist) + 195

	return &TripDetails{
		Method:   "Helicopter",
		Duration: totalTimeDuration,
		Cost:     cost,
		Distance: totalDist,
	}
}
