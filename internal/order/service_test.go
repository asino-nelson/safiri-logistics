package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestCreateLoadAllowsCustomerAndAutoMatches(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{}
	matcher := &stubMatcher{matched: true}
	service := NewService(repo, nil, matcher)

	now := time.Date(2026, 4, 11, 8, 0, 0, 0, time.UTC)
	load, err := service.CreateLoad(context.Background(), "customer-1", user.RoleCustomer, CreateLoadRequest{
		Title:            "Excavator",
		Description:      "Heavy machinery heading to Nairobi",
		Origin:           "Mombasa",
		Destination:      "Nairobi",
		WeightKG:         12000,
		PickupLatitude:   -4.043477,
		PickupLongitude:  39.668206,
		DropoffLatitude:  -1.286389,
		DropoffLongitude: 36.817223,
		PickupAt:         &now,
		Priority:         5,
		EquipmentType:    "lowbed",
		CargoType:        "machinery",
	})
	if err != nil {
		t.Fatalf("create load returned error: %v", err)
	}

	if matcher.lastLoadID == "" {
		t.Fatal("expected auto matcher to run")
	}

	if repo.created == nil || repo.created.PosterID != "customer-1" {
		t.Fatal("expected repository to receive posted load")
	}

	if load.Priority != 5 || load.EquipmentType != "lowbed" {
		t.Fatalf("unexpected load payload: %+v", load)
	}
}

func TestCreateLoadRejectsDriver(t *testing.T) {
	t.Parallel()

	service := NewService(&stubRepository{}, nil, nil)

	_, err := service.CreateLoad(context.Background(), "driver-1", user.RoleDriver, CreateLoadRequest{
		Title:            "Steel bars",
		Description:      "Construction materials",
		Origin:           "Mombasa",
		Destination:      "Kisumu",
		WeightKG:         2200,
		PickupLatitude:   -4.043477,
		PickupLongitude:  39.668206,
		DropoffLatitude:  -0.091702,
		DropoffLongitude: 34.767956,
	})
	if !errors.Is(err, ErrOnlyCustomersCanPostLoads) {
		t.Fatalf("expected customer-only error, got %v", err)
	}
}

func TestListLoadsUsesRoleAwareFiltering(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{
		posterLoads: []Load{{ID: "own-load"}},
		allLoads:    []Load{{ID: "open-load"}, {ID: "own-load"}},
	}
	service := NewService(repo, nil, nil)

	customerLoads, err := service.ListLoads(context.Background(), "customer-1", user.RoleCustomer)
	if err != nil {
		t.Fatalf("list customer loads: %v", err)
	}

	if len(customerLoads) != 1 || customerLoads[0].ID != "own-load" {
		t.Fatalf("unexpected customer loads: %+v", customerLoads)
	}

	driverLoads, err := service.ListLoads(context.Background(), "driver-1", user.RoleDriver)
	if err != nil {
		t.Fatalf("list driver loads: %v", err)
	}

	if len(driverLoads) != 2 {
		t.Fatalf("expected all loads for driver, got %+v", driverLoads)
	}
}

func TestPickLoadRequiresApprovedKYC(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{
		load: &Load{
			ID:       "load-1",
			PosterID: "customer-1",
			Status:   StatusPosted,
		},
	}

	service := NewService(repo, stubKYCVerifier{approved: false}, nil)

	_, err := service.PickLoad(context.Background(), "driver-1", user.RoleDriver, "load-1")
	if !errors.Is(err, ErrDriverKYCRequired) {
		t.Fatalf("expected kyc requirement error, got %v", err)
	}
}

func TestPickMatchedLoadOnlyAllowsAssignedDriver(t *testing.T) {
	t.Parallel()

	assignedDriver := "driver-1"
	repo := &stubRepository{
		load: &Load{
			ID:               "load-1",
			PosterID:         "customer-1",
			Status:           StatusMatched,
			AssignedDriverID: &assignedDriver,
		},
	}

	service := NewService(repo, stubKYCVerifier{approved: true}, nil)

	_, err := service.PickLoad(context.Background(), "driver-2", user.RoleDriver, "load-1")
	if !errors.Is(err, ErrLoadNotAvailable) {
		t.Fatalf("expected load not available error, got %v", err)
	}
}

func TestPickLoadAndProgressDelivery(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{
		load: &Load{
			ID:       "load-1",
			PosterID: "customer-1",
			Status:   StatusPosted,
		},
	}

	service := NewService(repo, stubKYCVerifier{approved: true}, nil)

	picked, err := service.PickLoad(context.Background(), "driver-1", user.RoleDriver, "load-1")
	if err != nil {
		t.Fatalf("pick load returned error: %v", err)
	}

	if picked.Status != StatusPicked {
		t.Fatalf("expected picked status, got %s", picked.Status)
	}

	inTransit, err := service.UpdateLoadStatus(context.Background(), "driver-1", user.RoleDriver, "load-1", StatusInTransit)
	if err != nil {
		t.Fatalf("update to in_transit returned error: %v", err)
	}

	if inTransit.Status != StatusInTransit {
		t.Fatalf("expected in_transit status, got %s", inTransit.Status)
	}

	delivered, err := service.UpdateLoadStatus(context.Background(), "driver-1", user.RoleDriver, "load-1", StatusDelivered)
	if err != nil {
		t.Fatalf("update to delivered returned error: %v", err)
	}

	if delivered.DeliveredAt == nil {
		t.Fatal("expected delivered timestamp to be set")
	}
}

type stubRepository struct {
	created     *Load
	posterLoads []Load
	allLoads    []Load
	load        *Load
}

func (r *stubRepository) Create(_ context.Context, load *Load) error {
	clone := *load
	r.created = &clone
	if r.load == nil {
		r.load = &clone
	}
	return nil
}

func (r *stubRepository) ListByPosterID(_ context.Context, _ string) ([]Load, error) {
	return append([]Load(nil), r.posterLoads...), nil
}

func (r *stubRepository) ListAll(_ context.Context) ([]Load, error) {
	return append([]Load(nil), r.allLoads...), nil
}

func (r *stubRepository) ListUnassignedByPickupWindow(_ context.Context, _, _ time.Time) ([]Load, error) {
	return nil, nil
}

func (r *stubRepository) GetByID(_ context.Context, _ string) (*Load, error) {
	if r.load == nil {
		return nil, ErrLoadNotFound
	}

	clone := *r.load
	return &clone, nil
}

func (r *stubRepository) Update(_ context.Context, load *Load) error {
	clone := *load
	r.load = &clone
	return nil
}

type stubKYCVerifier struct {
	approved bool
}

func (s stubKYCVerifier) IsApproved(_ context.Context, _ string) (bool, error) {
	return s.approved, nil
}

type stubMatcher struct {
	matched    bool
	lastLoadID string
}

func (s *stubMatcher) MatchLoad(_ context.Context, loadID, _ string) (bool, error) {
	s.lastLoadID = loadID
	return s.matched, nil
}
