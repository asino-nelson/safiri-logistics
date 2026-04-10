package dispatch

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/asino-nelson/safiri-logistics/internal/driver"
	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrNoEligibleDrivers     = errors.New("no eligible drivers available")
	ErrOnlyAdminsCanOptimize = errors.New("only admins can optimize dispatch")
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, payload map[string]any) error
}

type noopPublisher struct{}

func (noopPublisher) Publish(context.Context, string, map[string]any) error {
	return nil
}

type LoadRepository interface {
	GetByID(ctx context.Context, loadID string) (*order.Load, error)
	Update(ctx context.Context, load *order.Load) error
	ListUnassignedByPickupWindow(ctx context.Context, from, to time.Time) ([]order.Load, error)
}

type DriverDirectory interface {
	ListOperationalDrivers(ctx context.Context) ([]driver.Profile, error)
}

type Service struct {
	loads     LoadRepository
	drivers   DriverDirectory
	publisher EventPublisher
	nowFunc   func() time.Time
}

type OptimizeRequest struct {
	PickupDate time.Time `json:"pickup_date" binding:"required"`
}

type OptimizeResult struct {
	TotalLoads int            `json:"total_loads"`
	Matched    int            `json:"matched"`
	Skipped    int            `json:"skipped"`
	Results    []MatchOutcome `json:"results"`
}

type MatchOutcome struct {
	LoadID   string   `json:"load_id"`
	DriverID *string  `json:"driver_id,omitempty"`
	Score    *float64 `json:"score,omitempty"`
	Reason   string   `json:"reason"`
}

func NewService(loads LoadRepository, drivers DriverDirectory, publisher EventPublisher) *Service {
	if publisher == nil {
		publisher = noopPublisher{}
	}

	return &Service{
		loads:     loads,
		drivers:   drivers,
		publisher: publisher,
		nowFunc:   time.Now,
	}
}

func (s *Service) MatchLoad(ctx context.Context, loadID, source string) (bool, error) {
	load, err := s.loads.GetByID(ctx, loadID)
	if err != nil {
		return false, err
	}

	if load.AssignedDriverID != nil || load.Status != order.StatusPosted {
		return false, nil
	}

	drivers, err := s.drivers.ListOperationalDrivers(ctx)
	if err != nil {
		return false, err
	}

	candidate, score, ok := chooseBestCandidate(*load, drivers, nil)
	if !ok {
		return false, ErrNoEligibleDrivers
	}

	now := s.nowFunc().UTC()
	load.AssignedDriverID = &candidate.UserID
	load.Status = order.StatusMatched
	load.MatchedAt = &now
	load.MatchingScore = &score
	load.AssignmentSource = source
	load.UpdatedAt = now

	if err := s.loads.Update(ctx, load); err != nil {
		return false, err
	}

	_ = s.publisher.Publish(ctx, "load.matched", map[string]any{
		"load_id":           load.ID,
		"driver_id":         candidate.UserID,
		"assignment_source": source,
		"matching_score":    score,
	})

	return true, nil
}

func (s *Service) OptimizeAssignments(ctx context.Context, actorRole string, request OptimizeRequest) (*OptimizeResult, error) {
	if actorRole != user.RoleAdmin {
		return nil, ErrOnlyAdminsCanOptimize
	}

	from := request.PickupDate.UTC().Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)

	loads, err := s.loads.ListUnassignedByPickupWindow(ctx, from, to)
	if err != nil {
		return nil, err
	}

	drivers, err := s.drivers.ListOperationalDrivers(ctx)
	if err != nil {
		return nil, err
	}

	result := &OptimizeResult{
		TotalLoads: len(loads),
		Results:    make([]MatchOutcome, 0, len(loads)),
	}

	reservedDrivers := make(map[string]struct{})
	for _, load := range loads {
		candidate, score, ok := chooseBestCandidate(load, drivers, reservedDrivers)
		if !ok {
			result.Skipped++
			result.Results = append(result.Results, MatchOutcome{
				LoadID: load.ID,
				Reason: "no eligible drivers",
			})
			continue
		}

		now := s.nowFunc().UTC()
		loadCopy := load
		loadCopy.AssignedDriverID = &candidate.UserID
		loadCopy.Status = order.StatusMatched
		loadCopy.MatchedAt = &now
		loadCopy.MatchingScore = &score
		loadCopy.AssignmentSource = "admin_dispatch"
		loadCopy.UpdatedAt = now

		if err := s.loads.Update(ctx, &loadCopy); err != nil {
			return nil, err
		}

		reservedDrivers[candidate.UserID] = struct{}{}
		result.Matched++
		driverID := candidate.UserID
		scoreCopy := score
		result.Results = append(result.Results, MatchOutcome{
			LoadID:   load.ID,
			DriverID: &driverID,
			Score:    &scoreCopy,
			Reason:   "matched by priority-distance heuristic",
		})

		_ = s.publisher.Publish(ctx, "dispatch.optimized", map[string]any{
			"load_id":        load.ID,
			"driver_id":      candidate.UserID,
			"matching_score": score,
		})
	}

	return result, nil
}

func chooseBestCandidate(load order.Load, drivers []driver.Profile, reserved map[string]struct{}) (driver.Profile, float64, bool) {
	type scoredDriver struct {
		driver driver.Profile
		score  float64
	}

	candidates := make([]scoredDriver, 0, len(drivers))
	for _, candidate := range drivers {
		if reserved != nil {
			if _, exists := reserved[candidate.UserID]; exists {
				continue
			}
		}

		if !isEligible(load, candidate) {
			continue
		}

		candidates = append(candidates, scoredDriver{
			driver: candidate,
			score:  scoreCandidate(load, candidate),
		})
	}

	if len(candidates) == 0 {
		return driver.Profile{}, 0, false
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].driver.YearsExperience > candidates[j].driver.YearsExperience
		}

		return candidates[i].score < candidates[j].score
	})

	return candidates[0].driver, candidates[0].score, true
}

func isEligible(load order.Load, candidate driver.Profile) bool {
	if candidate.Status != driver.KYCStatusApproved {
		return false
	}

	if !candidate.IsOnline || !candidate.IsAvailable {
		return false
	}

	if candidate.CurrentLatitude == nil || candidate.CurrentLongitude == nil {
		return false
	}

	if candidate.MaxLoadKG > 0 && candidate.MaxLoadKG < load.WeightKG {
		return false
	}

	if candidate.EquipmentType != "" && load.EquipmentType != "" && candidate.EquipmentType != load.EquipmentType {
		return false
	}

	return true
}

func scoreCandidate(load order.Load, candidate driver.Profile) float64 {
	distanceKM := haversineKM(load.PickupLatitude, load.PickupLongitude, *candidate.CurrentLatitude, *candidate.CurrentLongitude)
	experiencePenalty := 10.0 / float64(candidate.YearsExperience+1)
	priorityPenalty := float64(6 - load.Priority)

	return distanceKM*0.65 + experiencePenalty*0.20 + priorityPenalty*0.15
}

func haversineKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKM = 6371.0

	dLat := radians(lat2 - lat1)
	dLon := radians(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(radians(lat1))*math.Cos(radians(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKM * c
}

func radians(value float64) float64 {
	return value * math.Pi / 180
}
