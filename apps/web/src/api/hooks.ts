import {
  getBusStations,
  type GetBusStationSearchParams,
} from '@/api/bus-stations.ts';
import { queryOptions } from '@tanstack/react-query';
import {
  getDepartures,
  type DeparturesSearchParams,
} from '@/api/departures.ts';

export const busStationQueryOptions = (params?: GetBusStationSearchParams) =>
  queryOptions({
    queryKey: ['bus-stations', params],
    queryFn: () => getBusStations(params),
  });

export const departuresQueryOptions = (params: DeparturesSearchParams) =>
  queryOptions({
    queryKey: ['departures', params],
    queryFn: () => getDepartures(params),
  });
