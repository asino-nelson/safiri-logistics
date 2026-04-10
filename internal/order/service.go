package order

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var ErrOnlyCustomersCanPostLoads = errors.New("only customers can post loads")

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
