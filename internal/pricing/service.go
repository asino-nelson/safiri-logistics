package pricing

import (
	"math"

	"github.com/asino-nelson/safiri-logistics/internal/dispatch"
	"github.com/asino-nelson/safiri-logistics/internal/order"
)

type Service struct{}

type Breakdown struct {
	BaseFeeKES            float64 `json:"base_fee_kes"`
	DistanceFeeKES        float64 `json:"distance_fee_kes"`
	WeightFeeKES          float64 `json:"weight_fee_kes"`
	PrioritySurchargeKES  float64 `json:"priority_surcharge_kes"`
	EquipmentSurchargeKES float64 `json:"equipment_surcharge_kes"`
	TotalKES              float64 `json:"total_kes"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Quote(load order.Load) float64 {
	return s.Breakdown(load).TotalKES
}

func (s *Service) Breakdown(load order.Load) Breakdown {
	distanceKM := dispatch.HaversineKM(load.PickupLatitude, load.PickupLongitude, load.DropoffLatitude, load.DropoffLongitude)
	baseFee := 15000.0
	distanceFee := distanceKM * 180.0
	weightFee := load.WeightKG * 0.45
	prioritySurcharge := float64(load.Priority-1) * 3500.0
	equipmentSurcharge := equipmentSurcharge(load.EquipmentType)

	total := baseFee + distanceFee + weightFee + prioritySurcharge + equipmentSurcharge

	return Breakdown{
		BaseFeeKES:            round2(baseFee),
		DistanceFeeKES:        round2(distanceFee),
		WeightFeeKES:          round2(weightFee),
		PrioritySurchargeKES:  round2(prioritySurcharge),
		EquipmentSurchargeKES: round2(equipmentSurcharge),
		TotalKES:              round2(total),
	}
}

func equipmentSurcharge(equipmentType string) float64 {
	switch equipmentType {
	case "lowbed":
		return 12000
	case "car_carrier":
		return 10000
	case "logging_trailer":
		return 9000
	default:
		return 5000
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
