package order

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrOnlyCustomersCanPostLoads   = errors.New("only customers can post loads")
	ErrOnlyDriversCanPickLoads     = errors.New("only drivers can pick loads")
	ErrDriverKYCRequired           = errors.New("driver must have approved kyc before picking loads")
	ErrLoadNotAvailable            = errors.New("load is not available for pickup")
	ErrOnlyAssignedDriverCanUpdate = errors.New("only the assigned driver can update load status")
	ErrInvalidLoadStatusTransition = errors.New("invalid load status transition")
)

type KYCVerifier interface {
	IsApproved(ctx context.Context, userID string) (bool, error)
}

type Service struct {
	repo        Repository
	kycVerifier KYCVerifier
	nowFunc     func() time.Time
}

func NewService(repo Repository, kycVerifier KYCVerifier) *Service {
	return &Service{
		repo:        repo,
		kycVerifier: kycVerifier,
		nowFunc:     time.Now,
	}
}

func (s *Service) CreateLoad(ctx context.Context, actorID, actorRole string, input CreateLoadRequest) (*Load, error) {
	if actorRole != user.RoleCustomer {
		return nil, ErrOnlyCustomersCanPostLoads
	}

	now := s.nowFunc().UTC()
	load := &Load{
		ID:          uuid.NewString(),
		PosterID:    actorID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Origin:      strings.TrimSpace(input.Origin),
		Destination: strings.TrimSpace(input.Destination),
		WeightKG:    input.WeightKG,
		Status:      StatusPosted,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, load); err != nil {
		return nil, err
	}

	return load, nil
}

func (s *Service) ListLoads(ctx context.Context, actorID, actorRole string) ([]Load, error) {
	if actorRole == user.RoleCustomer {
		return s.repo.ListByPosterID(ctx, actorID)
	}

	return s.repo.ListAll(ctx)
}

func (s *Service) PickLoad(ctx context.Context, actorID, actorRole, loadID string) (*Load, error) {
	if actorRole != user.RoleDriver {
		return nil, ErrOnlyDriversCanPickLoads
	}

	if s.kycVerifier == nil {
		return nil, ErrDriverKYCRequired
	}

	approved, err := s.kycVerifier.IsApproved(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if !approved {
		return nil, ErrDriverKYCRequired
	}

	load, err := s.repo.GetByID(ctx, loadID)
	if err != nil {
		return nil, err
	}

	if load.Status != StatusPosted || load.AssignedDriverID != nil {
		return nil, ErrLoadNotAvailable
	}

	now := s.nowFunc().UTC()
	load.AssignedDriverID = &actorID
	load.Status = StatusPicked
	load.PickedAt = &now
	load.UpdatedAt = now

	if err := s.repo.Update(ctx, load); err != nil {
		return nil, err
	}

	return load, nil
}

func (s *Service) UpdateLoadStatus(ctx context.Context, actorID, actorRole, loadID, newStatus string) (*Load, error) {
	if actorRole != user.RoleDriver {
		return nil, ErrOnlyAssignedDriverCanUpdate
	}

	load, err := s.repo.GetByID(ctx, loadID)
	if err != nil {
		return nil, err
	}

	if load.AssignedDriverID == nil || *load.AssignedDriverID != actorID {
		return nil, ErrOnlyAssignedDriverCanUpdate
	}

	nextStatus := strings.ToLower(strings.TrimSpace(newStatus))
	if !isValidTransition(load.Status, nextStatus) {
		return nil, ErrInvalidLoadStatusTransition
	}

	now := s.nowFunc().UTC()
	load.Status = nextStatus
	load.UpdatedAt = now
	if nextStatus == StatusDelivered {
		load.DeliveredAt = &now
	}

	if err := s.repo.Update(ctx, load); err != nil {
		return nil, err
	}

	return load, nil
}

func isValidTransition(current, next string) bool {
	switch current {
	case StatusPicked:
		return next == StatusInTransit
	case StatusInTransit:
		return next == StatusDelivered
	default:
		return false
	}
}
