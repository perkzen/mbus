import {
  ApiComboBox,
  type ComboBoxItem,
} from '@/components/ui/api-combobox.tsx';
import { busStationQueryOptions } from '@/api/hooks.ts';
import type {
  GetBusStationSearchParams,
  GetBusStationsResponse,
} from '@/api/bus-stations.ts';

const mapBusStationToOption = ({
  name,
  code,
}: {
  name: string;
  code: string;
}): ComboBoxItem => ({
  label: `${name} (${code})`,
  value: code,
});

type BusStationSelectProps = {
  selectedItem?: ComboBoxItem;
  onSelect: (item: ComboBoxItem) => void;
  className?: string;
};

export const BusStationSelect = ({
  selectedItem,
  onSelect,
  className,
}: BusStationSelectProps) => {
  return (
    <ApiComboBox<
      GetBusStationsResponse,
      Error,
      GetBusStationsResponse,
      (string | GetBusStationSearchParams | undefined)[]
    >
      selectedItem={selectedItem}
      onSelect={onSelect}
      queryOptionsFactory={(search) => busStationQueryOptions({ name: search })}
      mapDataToItems={(data) => data.map(mapBusStationToOption)}
      className={className}
    />
  );
};
