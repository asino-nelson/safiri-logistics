package order

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

const (
	defaultPriority      = 3
	defaultCargoType     = "general"
	defaultEquipmentType = "flatbed"
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

type AutoMatcher interface {
	MatchLoad(ctx context.Context, loadID, source string) (bool, error)
}

type Service struct {
	repo        Repository
	kycVerifier KYCVerifier
	autoMatcher AutoMatcher
	nowFunc     func() time.Time
}

func NewService(repo Repository, kycVerifier KYCVerifier, autoMatcher AutoMatcher) *Service {
	return &Service{
		repo:        repo,
		kycVerifier: kycVerifier,
		autoMatcher: autoMatcher,
		nowFunc:     time.Now,
	}
}

func (s *Service) CreateLoad(ctx context.Context, actorID, actorRole string, input CreateLoadRequest) (*Load, error) {
	if actorRole != user.RoleCustomer {
		return nil, ErrOnlyCustomersCanPostLoads
	}

	now := s.nowFunc().UTC()
	pickupAt := now
	if input.PickupAt != nil {
		pickupAt = input.PickupAt.UTC()
	}

	priority := input.Priority
	if priority == 0 {
		priority = defaultPriority
	}

	load := &Load{
		ID:               uuid.NewString(),
		PosterID:         actorID,
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		Origin:           strings.TrimSpace(input.Origin),
		Destination:      strings.TrimSpace(input.Destination),
		WeightKG:         input.WeightKG,
		Status:           StatusPosted,
		PickupLatitude:   input.PickupLatitude,
		PickupLongitude:  input.PickupLongitude,
		DropoffLatitude:  input.DropoffLatitude,
		DropoffLongitude: input.DropoffLongitude,
		PickupAt:         pickupAt,
		Priority:         priority,
		CargoType:        normalizeCargoType(input.CargoType),
		EquipmentType:    normalizeEquipmentType(input.EquipmentType),
		AssignmentSource: "manual",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, load); err != nil {
		return nil, err
	}

	if s.autoMatcher != nil {
		matched, err := s.autoMatcher.MatchLoad(ctx, load.ID, "auto")
		if err == nil && matched {
			reloaded, reloadErr := s.repo.GetByID(ctx, load.ID)
			if reloadErr == nil {
				return reloaded, nil
			}
		}
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

	if !canDriverPick(load, actorID) {
		return nil, ErrLoadNotAvailable
	}

	now := s.nowFunc().UTC()
	if load.AssignedDriverID == nil {
		load.AssignedDriverID = &actorID
	}
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

func canDriverPick(load *Load, actorID string) bool {
	switch load.Status {
	case StatusPosted:
		return load.AssignedDriverID == nil
	case StatusMatched:
		return load.AssignedDriverID != nil && *load.AssignedDriverID == actorID
	default:
		return false
	}
}

func isValidTransition(current, next string) bool {
	switch current {
	case StatusPicked, StatusMatched:
		return next == StatusInTransit
	case StatusInTransit:
		return next == StatusDelivered
	default:
		return false
	}
}

func normalizeCargoType(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return defaultCargoType
	}

	return normalized
}

func normalizeEquipmentType(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return defaultEquipmentType
	}

	return normalized
}
