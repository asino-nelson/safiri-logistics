package driver

import (
	"context"
	"errors"
	"testing"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestSubmitKYCRequiresDriverRole(t *testing.T) {
	t.Parallel()

	service := NewService(&stubRepository{})

	_, err := service.SubmitKYC(context.Background(), "customer-1", user.RoleCustomer, SubmitKYCRequest{
		NationalID:        "12345678",
		LicenseNumber:     "DL-23456",
		TruckRegistration: "KDA123A",
	})
	if !errors.Is(err, ErrOnlyDriversCanSubmitKYC) {
		t.Fatalf("expected driver-only error, got %v", err)
	}
}

func TestReviewKYCApprovesDriver(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{
		profile: &Profile{
			UserID:            "driver-1",
			NationalID:        "12345678",
			LicenseNumber:     "DL-23456",
			TruckRegistration: "KDA123A",
			Status:            KYCStatusPending,
		},
	}
	service := NewService(repo)

	profile, err := service.ReviewKYC(context.Background(), user.RoleAdmin, "driver-1", ReviewKYCRequest{
		Status: KYCStatusApproved,
	})
	if err != nil {
		t.Fatalf("review kyc returned error: %v", err)
	}

	if profile.Status != KYCStatusApproved {
		t.Fatalf("expected approved status, got %s", profile.Status)
	}

	approved, err := service.IsApproved(context.Background(), "driver-1")
	if err != nil {
		t.Fatalf("is approved returned error: %v", err)
	}

	if !approved {
		t.Fatal("expected approved driver profile")
	}
}

type stubRepository struct {
	profile *Profile
}

func (r *stubRepository) Upsert(_ context.Context, profile *Profile) error {
	clone := *profile
	r.profile = &clone
	return nil
}

func (r *stubRepository) GetByUserID(_ context.Context, _ string) (*Profile, error) {
	if r.profile == nil {
		return nil, ErrProfileNotFound
	}

	clone := *r.profile
	return &clone, nil
}
