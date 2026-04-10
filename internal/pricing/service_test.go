package pricing

import (
	"testing"

	"github.com/asino-nelson/safiri-logistics/internal/order"
)

func TestQuoteRewardsShorterDistanceAndLowerComplexity(t *testing.T) {
	t.Parallel()

	service := NewService()

	shortHaul := order.Load{
		PickupLatitude:   -1.286389,
		PickupLongitude:  36.817223,
		DropoffLatitude:  -1.292066,
		DropoffLongitude: 36.821945,
		WeightKG:         5000,
		Priority:         2,
		EquipmentType:    "flatbed",
	}

	heavyHaul := order.Load{
		PickupLatitude:   -4.043477,
		PickupLongitude:  39.668206,
		DropoffLatitude:  -1.286389,
		DropoffLongitude: 36.817223,
		WeightKG:         30000,
		Priority:         5,
		EquipmentType:    "lowbed",
	}

	shortPrice := service.Quote(shortHaul)
	heavyPrice := service.Quote(heavyHaul)

	if heavyPrice <= shortPrice {
		t.Fatalf("expected heavy haul to cost more, got short=%.2f heavy=%.2f", shortPrice, heavyPrice)
	}
}
