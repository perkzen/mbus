package departure

import (
	"context"
	"fmt"
	"github.com/perkzen/mbus/apps/bus-service/internal/errs"
	"github.com/perkzen/mbus/apps/bus-service/internal/provider/openrouteservice"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type Service struct {
	orsApiClient    *openrouteservice.APIClient
	cache           *redis.Client
	busStationStore store.BusStationStore
	departureStore  store.DepartureStore
	busLineStore    store.BusLineStore
	directionStore  store.DirectionStore
	enableCache     bool
}

type Option func(*Service)

func WithCache(enabled bool) Option {
	return func(s *Service) {
		s.enableCache = enabled
	}
}

func NewService(
	orsApiClient *openrouteservice.APIClient,
	cache *redis.Client,
	busStationStore store.BusStationStore,
	departureStore store.DepartureStore,
	busLineStore store.BusLineStore,
	directionStore store.DirectionStore,
	opts ...Option,

) *Service {
	s := &Service{
		orsApiClient:    orsApiClient,
		cache:           cache,
		busStationStore: busStationStore,
		departureStore:  departureStore,
		busLineStore:    busLineStore,
		directionStore:  directionStore,
		enableCache:     true,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) GenerateTimetable(fromID, toID int, date string) ([]TimetableRow, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("timetable_%d_%d_%s", fromID, toID, date)

	loader := func() ([]TimetableRow, error) {
		return s.buildDeparturesTimetable(fromID, toID, date)
	}

	if s.enableCache {
		return utils.WithCache(ctx, s.cache, cacheKey, 24*time.Hour, loader)
	}

	return loader()
}

func (s *Service) buildDeparturesTimetable(fromID, toID int, date string) ([]TimetableRow, error) {
	fromStation, err := s.busStationStore.FindBusStationByID(fromID)
	if err != nil {
		return nil, errs.BusStationNotFoundError(fromID)
	}

	toStation, err := s.busStationStore.FindBusStationByID(toID)
	if err != nil {
		return nil, errs.BusStationNotFoundError(toID)
	}

	locs := [][]float64{{fromStation.Lon, fromStation.Lat}, {toStation.Lon, toStation.Lat}}
	matrix, err := s.orsApiClient.GetMatrix(locs)
	if err != nil {
		return nil, err
	}

	distance := matrix.Distances[0][1]
	duration := matrix.Durations[0][1]
	schedule := store.ScheduleTyp(date)

	fromCode, toCode, departures, err := s.findValidDeparturePair(fromStation, toStation, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find valid departure pair: %w", err)
	}

	directions, err := s.directionStore.FindSharedDirectionsByCodes(fromCode, toCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find shared directions: %w", err)
	}

	toDeparturesMap, err := s.buildToDeparturesMap(toCode, directions, toStation, departures, schedule)
	if err != nil {
		return nil, err
	}

	rows := make([]TimetableRow, 0, len(departures))
	for _, dep := range departures {
		toDepTimes := toDeparturesMap[dep.Direction]

		arriveAt, err := getArriveAt(dep.DepartureTime, dep.Direction, toDepTimes, duration)
		if err != nil {
			continue
		}

		start, _ := utils.ParseClock(dep.DepartureTime)
		end, _ := utils.ParseClock(arriveAt)

		rows = append(rows, TimetableRow{
			ID:          dep.ID,
			Direction:   dep.Direction,
			Line:        dep.Line.Name,
			FromStation: Station{Name: fromStation.Name, ID: fromStation.ID},
			ToStation:   Station{Name: toStation.Name, ID: toStation.ID},
			DepartureAt: dep.DepartureTime,
			ArriveAt:    arriveAt,
			Duration:    utils.FormatDuration(start, end),
			Distance:    distance,
		})
	}

	utils.SortByDepartureAtAsc(rows)
	return rows, nil
}

func (s *Service) findValidDeparturePair(fromStation, toStation *store.BusStation, date string) (int, int, []store.Departure, error) {
	schedule := store.ScheduleTyp(date)

	// Try to find departures where toStation is final stop
	for _, fromCode := range fromStation.Codes {
		for _, toCode := range toStation.Codes {
			if departures, err := s.findDeparturesViaDirection(fromCode, toStation.SanitizedName(), schedule); err == nil && len(departures) > 0 {
				return fromCode, toCode, departures, nil
			}
		}
	}

	// Fallback to direct departure lookup
	for _, fromCode := range fromStation.Codes {
		for _, toCode := range toStation.Codes {
			departures, err := s.departureStore.FindDepartures(fromCode, toCode, schedule)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("failed to fetch departures from %d to %d: %w", fromCode, toCode, err)
			}
			if len(departures) > 0 {
				return fromCode, toCode, departures, nil
			}
		}
	}

	return 0, 0, nil, fmt.Errorf("no departures found between given station codes")
}

func (s *Service) findDeparturesViaDirection(fromCode int, sanitizedToName string, schedule store.ScheduleType) ([]store.Departure, error) {
	directions, err := s.directionStore.FindDirectionsByStationCode(fromCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find directions by station code %d: %w", fromCode, err)
	}

	for _, dir := range directions {
		if strings.HasSuffix(dir.Name, sanitizedToName) {
			return s.departureStore.FindDeparturesByStationCodeAndDirection(fromCode, dir.Name, schedule)
		}
	}
	return nil, nil
}

func (s *Service) buildToDeparturesMap(
	toCode int,
	directions []string,
	toStation *store.BusStation,
	departures []store.Departure,
	schedule store.ScheduleType,
) (map[string][]store.Departure, error) {
	toDeparturesMap := make(map[string][]store.Departure)

	for _, dir := range directions {
		toDepartures, err := s.departureStore.FindDeparturesByStationCodeAndDirection(toCode, dir, schedule)
		if err != nil {
			return nil, fmt.Errorf("failed to find departures for direction %s: %w", dir, err)
		}
		toDeparturesMap[dir] = toDepartures
	}

	ensureEmptyDirectionSlotIfFinalStop(toStation, departures, toDeparturesMap)

	return toDeparturesMap, nil
}

func ensureEmptyDirectionSlotIfFinalStop(toStation *store.BusStation, departures []store.Departure, toDeparturesMap map[string][]store.Departure) {
	for _, dep := range departures {
		if strings.HasSuffix(dep.Direction, toStation.SanitizedName()) && toDeparturesMap[dep.Direction] == nil {
			toDeparturesMap[dep.Direction] = []store.Departure{}
		}
	}
}

func getArriveAt(fromDepTime, direction string, toDeps []store.Departure, durationMin float64) (string, error) {
	fromTime, err := utils.ParseClock(fromDepTime)
	if err != nil {
		return "", fmt.Errorf("invalid from departure time: %w", err)
	}

	for _, toDep := range toDeps {
		if toDep.Direction != direction {
			continue
		}
		toTime, err := utils.ParseClock(toDep.DepartureTime)
		if err != nil {
			continue
		}
		if toTime.After(fromTime) {
			return toDep.DepartureTime, nil
		}
	}

	// Fallback: add travel duration to departure time
	arrival := fromTime.Add(time.Duration(durationMin * float64(time.Minute)))
	return arrival.Format("15:04"), nil
}
