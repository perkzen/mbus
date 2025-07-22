package departure

import (
	"context"
	"fmt"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/openrouteservice"
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
}

func NewService(
	orsApiClient *openrouteservice.APIClient,
	cache *redis.Client,
	busStationStore store.BusStationStore,
	departureStore store.DepartureStore,
	busLineStore store.BusLineStore,
	directionStore store.DirectionStore,
) *Service {
	return &Service{
		orsApiClient:    orsApiClient,
		cache:           cache,
		busStationStore: busStationStore,
		departureStore:  departureStore,
		busLineStore:    busLineStore,
		directionStore:  directionStore,
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

	toDeparturesMap := make(map[string][]store.Departure)
	for _, dir := range directions {
		toDepartures, err := s.departureStore.FindDeparturesByStationCodeAndDirection(toCode, dir, schedule)
		if err != nil {
			return nil, fmt.Errorf("failed to find departures for direction %s: %w", dir, err)
		}
		toDeparturesMap[dir] = toDepartures
	}

	// Handle edge case where toStation is a final stop
	for _, dep := range departures {
		if strings.HasSuffix(dep.Direction, toStation.Name) {
			if _, exists := toDeparturesMap[dep.Direction]; !exists {
				// Empty slice, it means we need to find all departures that end at this station
				toDeparturesMap[dep.Direction] = []store.Departure{}
			}
		}
	}

	rows := make([]TimetableRow, 0)
	for _, dep := range departures {
		toDepTimes, ok := toDeparturesMap[dep.Direction]
		if !ok {
			continue
		}

		arriveAt, err := getArriveAt(dep.DepartureTime, dep.Direction, toDepTimes, duration)
		if err != nil {
			return []TimetableRow{}, nil
		}

		start, _ := utils.ParseClock(dep.DepartureTime)
		end, _ := utils.ParseClock(arriveAt)
		formattedDuration := utils.FormatDuration(start, end)

		rows = append(rows, TimetableRow{
			ID:          dep.ID,
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

func (s *Service) findValidDeparturePair(fromStation, toStation *store.BusStation, date string) (int, int, []store.Departure, error) {
	schedule := store.ScheduleTyp(date)

	// Station codes lengths can range from 1 to 3, stations with 1 or 3 codes are edge cases.
	// Usually those stations tend to be final stops or major hubs and need to be handled differently.
	if len(toStation.Codes) != 2 {
		for _, fromCode := range fromStation.Codes {
			directions, err := s.directionStore.FindDirectionsByStationCode(fromCode)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("failed to find directions by station code %d: %w", fromCode, err)
			}

			allDepartures := make([]store.Departure, 0)
			for _, dir := range directions {
				// If the direction name ends with the toStation name, it means it's a valid direction for this station. (final stop)
				if strings.HasSuffix(dir.Name, toStation.Name) {
					departures, err := s.departureStore.FindDeparturesByStationCodeAndDirection(fromCode, dir.Name, schedule)
					if err != nil {
						return 0, 0, nil, fmt.Errorf("failed to find departures for direction %s: %w", dir.Name, err)
					}
					if len(departures) > 0 {
						allDepartures = append(allDepartures, departures...)
					}
				}
			}

			if len(allDepartures) > 0 {
				return fromCode, -1, allDepartures, nil
			}
		}
	}

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

func getArriveAt(fromDepTime, direction string, toDeps []store.Departure, durationMin float64) (string, error) {
	fromTime, err := utils.ParseClock(fromDepTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse from departure time: %w", err)
	}

	if len(toDeps) == 0 {
		arrivalTime := fromTime.Add(time.Duration(durationMin * float64(time.Minute)))
		return arrivalTime.Format("15:04"), nil
	}

	for _, toDep := range toDeps {
		if toDep.Direction != direction {
			continue
		}

		toTime, err := utils.ParseClock(toDep.DepartureTime)
		if err != nil {
			return "", fmt.Errorf("failed to parse departure time: %w", err)
		}

		fromTime, err := utils.ParseClock(fromDepTime)
		if err != nil {
			return "", fmt.Errorf("failed to parse from departure time: %w", err)
		}

		if toTime.After(fromTime) {
			return toDep.DepartureTime, nil
		}
	}

	return "", fmt.Errorf("no valid arrival time found for direction %s", direction)
}
