import {
  ApiComboBox,
  type ApiComboboxProps,
  type ComboBoxItem,
} from '@/components/ui/api-combobox.tsx';
import { busStationQueryOptions } from '@/api/query-options.ts';
import type {
  BusStation,
  BusStationSearchParams,
  GetBusStationsResponse,
} from '@/api/bus-stations.ts';

const mapBusStationToOption = ({
  name,
  code,
}: {
  name: string;
  code: number;
}): ComboBoxItem => ({
  label: `${name} (${code})`,
  value: code.toString(),
});

export function toComboBoxItem(station?: BusStation): ComboBoxItem {
  if (!station) {
    return { value: '', label: '' };
  }

  return {
    value: station.code.toString(),
    label: station.name
      ? `${station.name} (${station.code})`
      : `${station.code}`,
  };
}

type BusStationSelectProps = {
  className?: string;
} & Pick<
  ApiComboboxProps<
    GetBusStationsResponse,
    Error,
    GetBusStationsResponse,
    (string | BusStationSearchParams | undefined)[]
  >,
  'selectedItem' | 'onSelect' | 'searchPlaceholder' | 'selectPlaceholder'
>;

export const BusStationSelect = ({
  className,
  ...props
}: BusStationSelectProps) => {
  return (
    <ApiComboBox<
      GetBusStationsResponse,
      Error,
      GetBusStationsResponse,
      (string | BusStationSearchParams | undefined)[]
    >
      {...props}
      queryOptionsFactory={(search) => busStationQueryOptions({ name: search })}
      mapDataToItems={(data) => data.map(mapBusStationToOption)}
      className={className}
    />
  );
};
