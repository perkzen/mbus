export const ENDPOINTS = {
  BUS_STATIONS: '/bus-stations',
  DEPARTURES: '/departures',
} as const;

export const buildSearchParams = (params: Record<string, unknown>) => {
  return new URLSearchParams(
    Object.entries(params)
      .filter(([, value]) => value != null)
      .map(([key, value]) => [key, String(value)])
  );
};
