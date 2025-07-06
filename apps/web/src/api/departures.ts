import { buildSearchParams, ENDPOINTS } from '@/api/common.ts';
import { api } from '@/lib/axios.ts';

type Station = {
  name: string;
  code: number;
};

export type Departure = {
  direction: string;
  line: string;
  fromStation: Station;
  toStation: Station;
  duration: string;
  distance: number;
  departureAt: string;
  arriveAt: string;
};

export type DeparturesSearchParams = {
  from: string;
  to: string;
  date?: string;
};

export const getDepartures = async (params: DeparturesSearchParams) => {
  const searchParams = buildSearchParams(params);

  const url = `${ENDPOINTS.DEPARTURES}?${searchParams.toString()}`;
  const response = await api.get<Departure[]>(url);
  return response.data;
};
