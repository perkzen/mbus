import { useQuery } from '@tanstack/react-query';
import { busStationsQueryOptions } from '@/api/query-options';
import Combobox from '@/components/ui/combobox';
import type { FC } from 'react';

type Props = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  className?: string;
};

// TODO: start using this instead of API Combobox
const BusStationCombobox: FC<Props> = ({
  value,
  onChange,
  placeholder = 'Izberi postajo',
  searchPlaceholder = 'Poišči postajo',
  className,
}) => {
  const { data = [] } = useQuery(busStationsQueryOptions());

  const options = data.map((station) => ({
    label: station.name,
    value: station.id.toString(),
  }));

  return (
    <Combobox
      options={options}
      value={value}
      onChange={onChange}
      placeholder={placeholder}
      searchPlaceholder={searchPlaceholder}
      className={className}
    />
  );
};

export default BusStationCombobox;
