import { api } from '@/lib/axios.ts';
import { buildSearchParams, ENDPOINTS } from '@/api/common.ts';

export type GetBusStationSearchParams = {
  line?: string;
  name?: string;
  offset?: number;
  limit?: number;
};

export type GetBusStationsResponse = { name: string; code: string }[];

const DEFAULT_LIMIT = 10;
const DEFAULT_OFFSET = 0;

export const getBusStations = async (
  params: GetBusStationSearchParams = {}
) => {
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
