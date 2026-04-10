package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrOnlyLoadOwnerCanPay = errors.New("only the load owner can initiate payment")
	ErrInvalidQuotedAmount = errors.New("load has no payable quoted amount")
)

type LoadReader interface {
	GetByID(ctx context.Context, loadID string) (*order.Load, error)
}

type Service struct {
	repo     Repository
	loads    LoadReader
	provider Provider
	nowFunc  func() time.Time
}

func NewService(repo Repository, loads LoadReader, provider Provider) *Service {
	return &Service{
		repo:     repo,
		loads:    loads,
		provider: provider,
		nowFunc:  time.Now,
	}
}

func (s *Service) InitiateCheckout(ctx context.Context, actorID, actorRole, loadID string, request CheckoutRequest) (*Payment, error) {
	if actorRole != user.RoleCustomer {
		return nil, ErrOnlyLoadOwnerCanPay
	}

	load, err := s.loads.GetByID(ctx, loadID)
	if err != nil {
		return nil, err
	}

	if load.PosterID != actorID {
		return nil, ErrOnlyLoadOwnerCanPay
	}

	if load.QuotedPriceKES <= 0 {
		return nil, ErrInvalidQuotedAmount
	}

	response, err := s.provider.InitiateCheckout(ctx, ProviderRequest{
		PhoneNumber: request.PhoneNumber,
		AmountKES:   load.QuotedPriceKES,
		AccountRef:  load.ID,
		Description: fmt.Sprintf("Safiri load %s checkout", load.ID),
	})
	if err != nil {
		return nil, err
	}

	now := s.nowFunc().UTC()
	payment := &Payment{
		ID:                uuid.NewString(),
		LoadID:            load.ID,
		CustomerID:        actorID,
		Provider:          response.Provider,
		ProviderReference: response.ProviderReference,
		PhoneNumber:       request.PhoneNumber,
		AmountKES:         load.QuotedPriceKES,
		Currency:          "KES",
		Status:            response.Status,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
