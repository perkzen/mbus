import {
  ApiComboBox,
  type ApiComboboxProps,
  type ComboBoxItem,
} from '@/components/ui/api-combobox.tsx';
import { busStationsQueryOptions } from '@/api/query-options.ts';
import type {
  BusStation,
  BusStationSearchParams,
  GetBusStationsResponse,
} from '@/api/bus-stations.ts';

const mapBusStationToOption = ({
  name,
  id,
}: {
  name: string;
  id: number;
}): ComboBoxItem => ({
  label: name,
  value: id.toString(),
});

export function toComboBoxItem(station?: BusStation): ComboBoxItem {
  if (!station) {
    return { value: '', label: '' };
  }

  return {
    value: station.id.toString(),
    label: station.name,
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
      queryOptionsFactory={(search) =>
        busStationsQueryOptions({ name: search })
      }
      mapDataToItems={(data) => data.map(mapBusStationToOption)}
      className={className}
    />
  );
};
