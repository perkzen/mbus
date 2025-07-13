package departure

import (
	"context"
	"fmt"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/ors"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

type Service struct {
	orsApiClient    *ors.APIClient
	cache           *redis.Client
	busStationStore store.BusStationStore
	departureStore  store.DepartureStore
}

func NewService(orsApiClient *ors.APIClient, cache *redis.Client, busStationStore store.BusStationStore, departureStore store.DepartureStore) *Service {
	return &Service{
		orsApiClient:    orsApiClient,
		cache:           cache,
		busStationStore: busStationStore,
		departureStore:  departureStore,
	}
}

func (s *Service) GenerateTimetable(fromID, toID int, date string) ([]TimetableRow, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("timetable_%d_%d_%s", fromID, toID, date)

	if cached, ok := utils.TryGetFromCache[[]TimetableRow](ctx, s.cache, cacheKey); ok {
		return *cached, nil
	}

	fromStation, err := s.busStationStore.FindBusStationByID(fromID)
	if err != nil {
		return nil, fmt.Errorf("failed to find from station: %w", err)
	}

	toStation, err := s.busStationStore.FindBusStationByID(toID)
	if err != nil {
		return nil, fmt.Errorf("failed to find to station: %w", err)
	}

	if fromStation == nil || toStation == nil {
		return nil, fmt.Errorf("one of the bus stations not found: fromID=%d, toID=%d", fromID, toID)
	}

	locs := [][]float64{{fromStation.Lon, fromStation.Lat}, {toStation.Lon, toStation.Lat}}
	matrix, err := s.orsApiClient.GetMatrix(locs)
	if err != nil {
		return nil, err
	}
	distance := matrix.Distances[0][1]

	departures, err := s.departureStore.FindDeparturesFromStationToStation(fromID, toID, store.ScheduleTyp(date))
	if err != nil {
		return nil, fmt.Errorf("failed to find departures: %w", err)
	}

	directions, err := s.departureStore.FindSharedDirections(fromID, toID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shared directions: %w", err)
	}

	toDeparturesMap := make(map[string][]store.Departure)
	for _, dir := range directions {
		toDepartures, err := s.departureStore.FindDeparturesByStationIDAndDirection(toID, dir, store.ScheduleTyp(date))
		if err != nil {
			return nil, fmt.Errorf("failed to find departures for direction %s: %w", dir, err)
		}
		toDeparturesMap[dir] = toDepartures
	}

	rows := make([]TimetableRow, 0, len(departures))
	for _, dep := range departures {
		toDepTimes, ok := toDeparturesMap[dep.Direction]
		if !ok {
			continue
		}

		arriveAt, err := getArriveAt(dep.DepartureTime, dep.Direction, toDepTimes)
		if err != nil {
			return []TimetableRow{}, nil
		}

		start, _ := utils.ParseClock(dep.DepartureTime)
		end, _ := utils.ParseClock(arriveAt)
		formattedDuration := utils.FormatDuration(start, end)

		rows = append(rows, TimetableRow{
			Direction:   dep.Direction,
			Line:        dep.Line.Name,
			FromStation: Station{Name: fromStation.Name, ID: fromStation.ID},
			ToStation:   Station{Name: toStation.Name, ID: toStation.ID},
			DepartureAt: dep.DepartureTime,
			ArriveAt:    arriveAt,
			Duration:    formattedDuration,
			Distance:    distance,
		})

	}

	utils.SortByDepartureAtAsc(rows)
	utils.SaveToCache(ctx, s.cache, cacheKey, rows, 24*time.Hour)

	return rows, nil
}

func getArriveAt(fromDeTime, direction string, toDeps []store.Departure) (string, error) {
	for _, toDep := range toDeps {
		if toDep.Direction != direction {
			continue
		}

		toTime, err := utils.ParseClock(toDep.DepartureTime)
		if err != nil {
			return "", fmt.Errorf("failed to parse departure time: %w", err)
		}

		fromTime, err := utils.ParseClock(fromDeTime)
		if err != nil {
			return "", fmt.Errorf("failed to parse from departure time: %w", err)
		}

		if toTime.After(fromTime) {
			return toDep.DepartureTime, nil
		}
	}

	return "", fmt.Errorf("no valid arrival time found for direction %s", direction)

}
