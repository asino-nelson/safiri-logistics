package tracking

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrOnlyAssignedDriverCanTrack = errors.New("only assigned driver can post tracking updates")
)

type LoadReader interface {
	GetByID(ctx context.Context, loadID string) (*order.Load, error)
}

type Service struct {
	repo    Repository
	loads   LoadReader
	hub     *Hub
	nowFunc func() time.Time
}

func NewService(repo Repository, loads LoadReader, hub *Hub) *Service {
	if hub == nil {
		hub = NewHub()
	}

	return &Service{
		repo:    repo,
		loads:   loads,
		hub:     hub,
		nowFunc: time.Now,
	}
}

func (s *Service) CreateEvent(ctx context.Context, actorID, actorRole, loadID string, request CreateEventRequest) (*Event, error) {
	if actorRole != user.RoleDriver {
		return nil, ErrOnlyAssignedDriverCanTrack
	}

	load, err := s.loads.GetByID(ctx, loadID)
	if err != nil {
		return nil, err
	}

	if load.AssignedDriverID == nil || *load.AssignedDriverID != actorID {
		return nil, ErrOnlyAssignedDriverCanTrack
	}

	event := &Event{
		ID:             uuid.NewString(),
		LoadID:         loadID,
		DriverID:       actorID,
		Latitude:       request.Latitude,
		Longitude:      request.Longitude,
		SpeedKPH:       request.SpeedKPH,
		HeadingDegrees: request.HeadingDegrees,
		RecordedAt:     s.nowFunc().UTC(),
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}

	s.hub.Publish(loadID, *event)
	return event, nil
}

func (s *Service) Subscribe(loadID string) chan Event {
	return s.hub.Subscribe(loadID)
}

func (s *Service) Unsubscribe(loadID string, ch chan Event) {
	s.hub.Unsubscribe(loadID, ch)
}
