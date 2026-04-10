package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestInitiateCheckoutRequiresLoadOwner(t *testing.T) {
	t.Parallel()

	service := NewService(&stubRepository{}, stubLoadReader{
		load: &order.Load{
			ID:             "load-1",
			PosterID:       "customer-1",
			QuotedPriceKES: 150000,
		},
	}, NewSandboxMPESAClient())

	_, err := service.InitiateCheckout(context.Background(), "customer-2", user.RoleCustomer, "load-1", CheckoutRequest{
		PhoneNumber: "254700000001",
	})
	if !errors.Is(err, ErrOnlyLoadOwnerCanPay) {
		t.Fatalf("expected load-owner error, got %v", err)
	}
}

func TestInitiateCheckoutPersistsPayment(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{}
	service := NewService(repo, stubLoadReader{
		load: &order.Load{
			ID:             "load-1",
			PosterID:       "customer-1",
			QuotedPriceKES: 285000,
		},
	}, NewSandboxMPESAClient())

	payment, err := service.InitiateCheckout(context.Background(), "customer-1", user.RoleCustomer, "load-1", CheckoutRequest{
		PhoneNumber: "254700000001",
	})
	if err != nil {
		t.Fatalf("initiate checkout returned error: %v", err)
	}

	if payment.Provider != "mpesa" || payment.Status != StatusInitiated {
		t.Fatalf("unexpected payment response: %+v", payment)
	}

	if repo.created == nil || repo.created.ProviderReference == "" {
		t.Fatalf("expected payment persistence, got %+v", repo.created)
	}
}

type stubRepository struct {
	created *Payment
}

func (r *stubRepository) Create(_ context.Context, payment *Payment) error {
	clone := *payment
	r.created = &clone
	return nil
}

type stubLoadReader struct {
	load *order.Load
}

func (s stubLoadReader) GetByID(_ context.Context, _ string) (*order.Load, error) {
	if s.load == nil {
		return nil, order.ErrLoadNotFound
	}

	clone := *s.load
	return &clone, nil
}
