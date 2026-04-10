package tracking

import (
	"context"
	"errors"
	"testing"

	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestCreateEventRequiresAssignedDriver(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{}
	loads := stubLoadReader{
		load: &order.Load{
			ID:       "load-1",
			PosterID: "customer-1",
			Status:   order.StatusMatched,
		},
	}

	service := NewService(repo, loads, NewHub())
	_, err := service.CreateEvent(context.Background(), "driver-1", user.RoleDriver, "load-1", CreateEventRequest{
		Latitude:  -1.286389,
		Longitude: 36.817223,
	})
	if !errors.Is(err, ErrOnlyAssignedDriverCanTrack) {
		t.Fatalf("expected assigned driver error, got %v", err)
	}
}

func TestCreateEventPublishesToSubscribers(t *testing.T) {
	t.Parallel()

	assignedDriver := "driver-1"
	repo := &stubRepository{}
	loads := stubLoadReader{
		load: &order.Load{
			ID:               "load-1",
			PosterID:         "customer-1",
			Status:           order.StatusInTransit,
			AssignedDriverID: &assignedDriver,
		},
	}

	hub := NewHub()
	service := NewService(repo, loads, hub)
	stream := service.Subscribe("load-1")
	defer service.Unsubscribe("load-1", stream)

	event, err := service.CreateEvent(context.Background(), "driver-1", user.RoleDriver, "load-1", CreateEventRequest{
		Latitude:       -1.286389,
		Longitude:      36.817223,
		SpeedKPH:       62,
		HeadingDegrees: 130,
	})
	if err != nil {
		t.Fatalf("create event returned error: %v", err)
	}

	if repo.created == nil || repo.created.ID != event.ID {
		t.Fatalf("expected event persistence, got %+v", repo.created)
	}

	select {
	case broadcast := <-stream:
		if broadcast.ID != event.ID {
			t.Fatalf("expected broadcast %s, got %s", event.ID, broadcast.ID)
		}
	default:
		t.Fatal("expected websocket subscriber to receive event")
	}
}

type stubRepository struct {
	created *Event
}

func (r *stubRepository) Create(_ context.Context, event *Event) error {
	clone := *event
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
