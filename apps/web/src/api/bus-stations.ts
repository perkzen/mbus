import { api } from '@/lib/axios.ts';
import { buildSearchParams, ENDPOINTS } from '@/api/common.ts';

export type BusStationSearchParams = {
  line?: string;
  name?: string;
  offset?: number;
  limit?: number;
};

export type BusStation = {
  name: string;
  code: number;
};

export type GetBusStationsResponse = BusStation[];

const DEFAULT_LIMIT = 10;
const DEFAULT_OFFSET = 0;

export const getBusStations = async (params: BusStationSearchParams = {}) => {
  const mergedParams = {
    limit: DEFAULT_LIMIT,
    offset: DEFAULT_OFFSET,
    ...params,
  };

  const searchParams = buildSearchParams(mergedParams);

  const url = `${ENDPOINTS.BUS_STATIONS}?${searchParams.toString()}`;

  const response = await api.get<GetBusStationsResponse>(url);
  return response.data;
};

export const getBusStationByCode = async (code: number) => {
  const url = `${ENDPOINTS.BUS_STATIONS}/${code}`;

  const response = await api.get<BusStation>(url);
  return response.data;
};
