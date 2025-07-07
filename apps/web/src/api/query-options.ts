import {
  getBusStations,
  type BusStationSearchParams,
  getBusStationByCode,
} from '@/api/bus-stations.ts';
import { queryOptions } from '@tanstack/react-query';
import {
  getDepartures,
  type DeparturesSearchParams,
} from '@/api/departures.ts';

export const busStationQueryOptions = (params?: BusStationSearchParams) =>
  queryOptions({
    queryKey: ['bus-stations', params],
    queryFn: () => getBusStations(params),
  });

export const departuresQueryOptions = (params: DeparturesSearchParams) =>
  queryOptions({
    queryKey: ['departures', params],
    queryFn: () => getDepartures(params),
    enabled: !!params.from && !!params.to,
  });

export const busStationByCodeQueryOptions = (code: number) =>
  queryOptions({
    queryKey: ['bus-stations', code],
    queryFn: () => getBusStationByCode(code),
    enabled: !!code,
  });
