package dispatch

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asino-nelson/safiri-logistics/internal/driver"
	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestMatchLoadPicksNearestEligibleDriver(t *testing.T) {
	t.Parallel()

	loadRepo := &stubLoadRepository{
		load: &order.Load{
			ID:               "load-1",
			Status:           order.StatusPosted,
			WeightKG:         12000,
			Priority:         4,
			EquipmentType:    "lowbed",
			PickupLatitude:   -1.286389,
			PickupLongitude:  36.817223,
			AssignmentSource: "manual",
		},
	}

	nearLat := -1.286500
	nearLng := 36.817000
	farLat := -4.043477
	farLng := 39.668206

	driverDirectory := stubDriverDirectory{
		drivers: []driver.Profile{
			{
				UserID:           "driver-near",
				Status:           driver.KYCStatusApproved,
				YearsExperience:  7,
				MaxLoadKG:        20000,
				EquipmentType:    "lowbed",
				IsOnline:         true,
				IsAvailable:      true,
				CurrentLatitude:  &nearLat,
				CurrentLongitude: &nearLng,
			},
			{
				UserID:           "driver-far",
				Status:           driver.KYCStatusApproved,
				YearsExperience:  12,
				MaxLoadKG:        40000,
				EquipmentType:    "lowbed",
				IsOnline:         true,
				IsAvailable:      true,
				CurrentLatitude:  &farLat,
				CurrentLongitude: &farLng,
			},
		},
	}

	service := NewService(loadRepo, driverDirectory, nil)

	matched, err := service.MatchLoad(context.Background(), "load-1", "auto")
	if err != nil {
		t.Fatalf("match load returned error: %v", err)
	}

	if !matched {
		t.Fatal("expected load to match")
	}

	if loadRepo.updated == nil || loadRepo.updated.AssignedDriverID == nil || *loadRepo.updated.AssignedDriverID != "driver-near" {
		t.Fatalf("expected nearest driver assignment, got %+v", loadRepo.updated)
	}

	if loadRepo.updated.Status != order.StatusMatched {
		t.Fatalf("expected matched status, got %s", loadRepo.updated.Status)
	}
}

func TestOptimizeAssignmentsHonorsPriorityAndExperience(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 11, 8, 0, 0, 0, time.UTC)
	latNairobi := -1.286389
	lngNairobi := 36.817223
	latThika := -1.033200
	lngThika := 37.069200

	loadRepo := &stubLoadRepository{
		loads: []order.Load{
			{
				ID:               "high-priority",
				Status:           order.StatusPosted,
				WeightKG:         30000,
				Priority:         5,
				EquipmentType:    "flatbed",
				PickupLatitude:   latNairobi,
				PickupLongitude:  lngNairobi,
				PickupAt:         now,
				AssignmentSource: "manual",
			},
			{
				ID:               "standard-priority",
				Status:           order.StatusPosted,
				WeightKG:         15000,
				Priority:         2,
				EquipmentType:    "flatbed",
				PickupLatitude:   latThika,
				PickupLongitude:  lngThika,
				PickupAt:         now.Add(2 * time.Hour),
				AssignmentSource: "manual",
			},
		},
	}

	driverDirectory := stubDriverDirectory{
		drivers: []driver.Profile{
			{
				UserID:           "experienced-driver",
				Status:           driver.KYCStatusApproved,
				YearsExperience:  12,
				MaxLoadKG:        36000,
				EquipmentType:    "flatbed",
				IsOnline:         true,
				IsAvailable:      true,
				CurrentLatitude:  &latNairobi,
				CurrentLongitude: &lngNairobi,
			},
			{
				UserID:           "junior-driver",
				Status:           driver.KYCStatusApproved,
				YearsExperience:  2,
				MaxLoadKG:        18000,
				EquipmentType:    "flatbed",
				IsOnline:         true,
				IsAvailable:      true,
				CurrentLatitude:  &latThika,
				CurrentLongitude: &lngThika,
			},
		},
	}

	service := NewService(loadRepo, driverDirectory, nil)
	result, err := service.OptimizeAssignments(context.Background(), user.RoleAdmin, OptimizeRequest{
		PickupDate: now,
	})
	if err != nil {
		t.Fatalf("optimize assignments returned error: %v", err)
	}

	if result.Matched != 2 {
		t.Fatalf("expected both loads matched, got %+v", result)
	}

	if len(loadRepo.updates) != 2 {
		t.Fatalf("expected two persisted updates, got %d", len(loadRepo.updates))
	}

	if *loadRepo.updates[0].AssignedDriverID != "experienced-driver" {
		t.Fatalf("expected highest priority load to receive experienced driver, got %+v", loadRepo.updates[0])
	}
}

func TestOptimizeAssignmentsRejectsNonAdmins(t *testing.T) {
	t.Parallel()

	service := NewService(&stubLoadRepository{}, stubDriverDirectory{}, nil)

	_, err := service.OptimizeAssignments(context.Background(), user.RoleCustomer, OptimizeRequest{
		PickupDate: time.Now(),
	})
	if !errors.Is(err, ErrOnlyAdminsCanOptimize) {
		t.Fatalf("expected admin-only error, got %v", err)
	}
}

type stubLoadRepository struct {
	load    *order.Load
	loads   []order.Load
	updated *order.Load
	updates []*order.Load
}

func (r *stubLoadRepository) GetByID(_ context.Context, _ string) (*order.Load, error) {
	if r.load == nil {
		return nil, order.ErrLoadNotFound
	}

	clone := *r.load
	return &clone, nil
}

func (r *stubLoadRepository) Update(_ context.Context, load *order.Load) error {
	clone := *load
	r.updated = &clone
	r.updates = append(r.updates, &clone)
	if r.load != nil && r.load.ID == clone.ID {
		r.load = &clone
	}
	return nil
}

func (r *stubLoadRepository) ListUnassignedByPickupWindow(_ context.Context, _, _ time.Time) ([]order.Load, error) {
	return append([]order.Load(nil), r.loads...), nil
}

type stubDriverDirectory struct {
	drivers []driver.Profile
}

func (s stubDriverDirectory) ListOperationalDrivers(_ context.Context) ([]driver.Profile, error) {
	return append([]driver.Profile(nil), s.drivers...), nil
}
