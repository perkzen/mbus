package departure

import (
	"fmt"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/ors"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/redis/go-redis/v9"
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
	//ctx := context.Background()
	//cacheKey := fmt.Sprintf("timetable_%d_%d_%s", fromID, toID, date)

	//if cached, ok := utils.TryGetFromCache[[]TimetableRow](ctx, s.cache, cacheKey); ok {
	//	return *cached, nil
	//}

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

	//locs := [][]float64{{fromStation.Lon, fromStation.Lat}, {toStation.Lon, toStation.Lat}}
	//matrix, err := s.orsApiClient.GetMatrix(locs)
	//if err != nil {
	//	return nil, err
	//}
	//distance := matrix.Distances[0][1]

	scheduleType := store.ScheduleTyp(date)

	var fromCode, toCode int
	var departures []store.Departure

	for _, fc := range fromStation.Codes {
		for _, tc := range toStation.Codes {
			d, err := s.departureStore.FindDepartures(fc, tc, scheduleType)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch departures from %d to %d: %w", fc, tc, err)
			}
			if len(d) > 0 {
				fromCode = fc
				toCode = tc
				departures = d
				// ❌ no break — keep checking for later matches
			}
		}
	}

	directions, err := s.departureStore.FindSharedDirectionsByCodes(fromCode, toCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find shared directions: %w", err)
	}

	toDeparturesMap := make(map[string][]store.Departure)
	for _, dir := range directions {
		toDepartures, err := s.departureStore.FindDeparturesByStationCodeAndDirection(toCode, dir, scheduleType)
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
			Distance:    0,
		})

	}

	utils.SortByDepartureAtAsc(rows)
	//utils.SaveToCache(ctx, s.cache, cacheKey, rows, 24*time.Hour)

	return rows, nil
}

func getArriveAt(fromDepTime, direction string, toDeps []store.Departure) (string, error) {
	for _, toDep := range toDeps {
		if toDep.Direction != direction {
			continue
		}

		println(fromDepTime, toDep.DepartureTime)

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
