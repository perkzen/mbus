package departure

import (
	"context"
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/bus-service/internal/integrations/ors"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"github.com/perkzen/mbus/bus-service/internal/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

type Service struct {
	marpromClient   *marprom.APIClient
	orsApiClient    *ors.APIClient
	cache           *redis.Client
	busStationStore store.BusStationStore
}

func NewService(marpromClient *marprom.APIClient, orsApiClient *ors.APIClient, cache *redis.Client, busStationStore store.BusStationStore) *Service {
	return &Service{
		marpromClient:   marpromClient,
		orsApiClient:    orsApiClient,
		cache:           cache,
		busStationStore: busStationStore,
	}
}

func (s *Service) GenerateTimetable(fromCode, toCode int, date string) ([]Departure, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("timetable_%d_%d_%s", fromCode, toCode, date)

	if cached, ok := utils.TryGetFromCache[[]Departure](ctx, s.cache, cacheKey); ok {
		return *cached, nil
	}

	fromStation, err := s.busStationStore.FindBusStationByCode(fromCode)
	if err != nil {
		return nil, err
	}

	toStation, err := s.busStationStore.FindBusStationByCode(toCode)
	if err != nil {
		return nil, err
	}

	fromDetails, err := s.marpromClient.GetBusStationDetails(fromCode, date)
	if err != nil {
		return nil, err
	}

	toDetails, err := s.marpromClient.GetBusStationDetails(toCode, date)
	if err != nil {
		return nil, err
	}

	departures, err := s.marpromClient.GetDeparturesFromStationToStation(fromCode, toCode, date)
	if err != nil {
		return nil, err
	}

	locs := [][]float64{{fromStation.Lon, fromStation.Lat}, {toStation.Lon, toStation.Lat}}
	matrix, err := s.orsApiClient.GetMatrix(locs)
	if err != nil {
		return nil, err
	}
	distance := matrix.Distances[0][1]

	rows := make([]Departure, 0, len(departures))
	for _, dep := range departures {
		for _, depTime := range dep.Times {
			arriveAt, _ := s.getArriveAt(depTime, dep.Line, dep.Direction, toDetails.Departures)

			start, _ := utils.ParseClock(depTime)
			end, _ := utils.ParseClock(arriveAt)
			formattedDuration := utils.FormatDuration(start, end)

			rows = append(rows, Departure{
				Direction:   dep.Direction,
				Line:        dep.Line,
				FromStation: Station{Name: fromDetails.Name, Code: fromCode},
				ToStation:   Station{Name: toDetails.Name, Code: toCode},
				DepartureAt: depTime,
				ArriveAt:    arriveAt,
				Duration:    formattedDuration,
				Distance:    distance,
			})
		}
	}

	utils.SaveToCache(ctx, s.cache, cacheKey, rows, 24*time.Hour)
	return rows, nil
}

func (s *Service) getArriveAt(fromDepTime, line, direction string, toDepTimes []marprom.Departure) (string, error) {
	for _, toDep := range toDepTimes {

		if toDep.Direction != direction {
			continue
		}

		for _, toDepTime := range toDep.Times {
			fromTime, _ := utils.ParseClock(fromDepTime)
			toTime, _ := utils.ParseClock(toDepTime)
			if toTime.After(fromTime) {
				return toDepTime, nil
			}

		}
	}

	return "", fmt.Errorf("station is not in route or no valid arrival time found for line %s and direction %s", line, direction)
}
