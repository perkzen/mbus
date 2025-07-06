import { createFileRoute } from '@tanstack/react-router';
import DepartureSearchForm from '@/components/blocks/departure-search-form.tsx';
import { useState } from 'react';
import type { ComboBoxItem } from '@/components/ui/api-combobox.tsx';
import { useQuery } from '@tanstack/react-query';
import { departuresQueryOptions } from '@/api/hooks.ts';
import DeparturesTable from '@/components/blocks/departures-table.tsx';

export const Route = createFileRoute('/')({
  component: App,
});

function App() {
  const [fromStation, setFromStation] = useState<ComboBoxItem>();
  const [toStation, setToStation] = useState<ComboBoxItem>();

  const { data, refetch: searchDepartures } = useQuery({
    ...departuresQueryOptions({
      from: fromStation?.value || '',
      to: toStation?.value || '',
    }),
    enabled: false,
  });

  return (
    <div className="mx-auto max-w-6xl p-10">
      <DepartureSearchForm
        searchDepartures={searchDepartures}
        setFromStation={setFromStation}
        setToStation={setToStation}
        toStation={toStation}
        fromStation={fromStation}
      />
      <DeparturesTable data={data || []} />
    </div>
  );
}
