package user

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}
