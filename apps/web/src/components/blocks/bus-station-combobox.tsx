import { useQuery } from '@tanstack/react-query';
import { busStationsQueryOptions } from '@/api/query-options';
import type { FC } from 'react';
import { Combobox } from '@/components/ui/combobox.tsx';

type Props = {
  value?: string;
  onChange: (value: string) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyPlaceholder?: string;
  className?: string;
};

const BusStationCombobox: FC<Props> = ({
  value,
  onChange,
  placeholder = 'Izberi postajo',
  searchPlaceholder = 'Poišči postajo',
  emptyPlaceholder = 'Ni rezultatov',
  className,
}) => {
  const { data } = useQuery(busStationsQueryOptions());

  const options =
    data?.map((station) => ({
      label: station.name,
      value: station.id.toString(),
    })) || [];

  return (
    <Combobox
      options={options}
      selected={value}
      onChange={onChange}
      placeholder={placeholder}
      searchPlaceholder={searchPlaceholder}
      emptyPlaceholder={emptyPlaceholder}
      className={className}
    />
  );
};

export default BusStationCombobox;
