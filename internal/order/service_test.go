package order

import (
	"context"
	"errors"
	"testing"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestCreateLoadAllowsCustomer(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{}
	service := NewService(repo, nil)

	load, err := service.CreateLoad(context.Background(), "customer-1", user.RoleCustomer, CreateLoadRequest{
		Title:       "Fresh produce",
		Description: "Vegetables headed to Nairobi",
		Origin:      "Eldoret",
		Destination: "Nairobi",
		WeightKG:    1200,
	})
	if err != nil {
		t.Fatalf("create load returned error: %v", err)
	}

	if load.Status != StatusPosted {
		t.Fatalf("expected posted status, got %s", load.Status)
	}

	if repo.created == nil || repo.created.PosterID != "customer-1" {
		t.Fatal("expected repository to receive posted load")
	}
}

func TestCreateLoadRejectsDriver(t *testing.T) {
	t.Parallel()

	service := NewService(&stubRepository{}, nil)

	_, err := service.CreateLoad(context.Background(), "driver-1", user.RoleDriver, CreateLoadRequest{
		Title:       "Steel bars",
		Description: "Construction materials",
		Origin:      "Mombasa",
		Destination: "Kisumu",
		WeightKG:    2200,
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
	service := NewService(repo, nil)
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

	service := NewService(repo, stubKYCVerifier{approved: false})

	_, err := service.PickLoad(context.Background(), "driver-1", user.RoleDriver, "load-1")
	if !errors.Is(err, ErrDriverKYCRequired) {
		t.Fatalf("expected kyc requirement error, got %v", err)
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

	service := NewService(repo, stubKYCVerifier{approved: true})

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
	return nil
}

func (r *stubRepository) ListByPosterID(_ context.Context, _ string) ([]Load, error) {
	return append([]Load(nil), r.posterLoads...), nil
}

func (r *stubRepository) ListAll(_ context.Context) ([]Load, error) {
	return append([]Load(nil), r.allLoads...), nil
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
