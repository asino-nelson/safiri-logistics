package driver

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrOnlyDriversCanSubmitKYC = errors.New("only drivers can submit kyc")
	ErrOnlyAdminsCanReviewKYC  = errors.New("only admins can review driver kyc")
	ErrInvalidKYCStatus        = errors.New("invalid kyc status")
)

type Service struct {
	repo    Repository
	nowFunc func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:    repo,
		nowFunc: time.Now,
	}
}

func (s *Service) SubmitKYC(ctx context.Context, actorID, actorRole string, input SubmitKYCRequest) (*Profile, error) {
	if actorRole != user.RoleDriver {
		return nil, ErrOnlyDriversCanSubmitKYC
	}

	profile := &Profile{
		UserID:            actorID,
		NationalID:        strings.TrimSpace(input.NationalID),
		LicenseNumber:     strings.TrimSpace(input.LicenseNumber),
		TruckRegistration: strings.TrimSpace(input.TruckRegistration),
		Status:            KYCStatusPending,
		SubmittedAt:       s.nowFunc().UTC(),
	}

	if err := s.repo.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) ReviewKYC(ctx context.Context, actorRole, driverUserID string, input ReviewKYCRequest) (*Profile, error) {
	if actorRole != user.RoleAdmin {
		return nil, ErrOnlyAdminsCanReviewKYC
	}

	profile, err := s.repo.GetByUserID(ctx, driverUserID)
	if err != nil {
		return nil, err
	}

	now := s.nowFunc().UTC()
	status := strings.ToLower(strings.TrimSpace(input.Status))
	switch status {
	case KYCStatusApproved:
		profile.Status = KYCStatusApproved
		profile.ReviewedAt = &now
		profile.VerifiedAt = &now
		profile.RejectionReason = nil
	case KYCStatusRejected:
		reason := strings.TrimSpace(input.RejectionReason)
		profile.Status = KYCStatusRejected
		profile.ReviewedAt = &now
		profile.VerifiedAt = nil
		profile.RejectionReason = &reason
	default:
		return nil, ErrInvalidKYCStatus
	}

	if err := s.repo.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) IsApproved(ctx context.Context, userID string) (bool, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return false, nil
		}

		return false, err
	}

	return profile.Status == KYCStatusApproved, nil
}
