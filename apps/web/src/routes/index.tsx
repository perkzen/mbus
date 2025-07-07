import { createFileRoute, stripSearchParams } from '@tanstack/react-router';
import DepartureSearchForm from '@/components/blocks/departure-search-form.tsx';
import { useQuery } from '@tanstack/react-query';
import { departuresQueryOptions } from '@/api/query-options.ts';
import DeparturesTable from '@/components/blocks/departures-table.tsx';
import { z } from 'zod';
import { format } from 'date-fns';

const defaultSearchParams = {
  from: 0,
  to: 0,
  date: format(new Date(), 'yyyy-MM-dd'),
};

const departureSearchSchema = z.object({
  from: z.number().catch(defaultSearchParams.from),
  to: z.number().catch(defaultSearchParams.to),
  date: z
    .string()
    .regex(/^\d{4}-\d{2}-\d{2}$/)
    .catch(defaultSearchParams.date),
});

export type DepartureSearchParams = z.infer<typeof departureSearchSchema>;

export const Route = createFileRoute('/')({
  component: App,
  validateSearch: (search: Record<string, unknown>) =>
    departureSearchSchema.parse(search),
  search: {
    middlewares: [stripSearchParams(defaultSearchParams)],
  },
});

function App() {
  const searchParams = Route.useSearch();

  const {
    data,
    refetch: searchDepartures,
    isFetching,
  } = useQuery(departuresQueryOptions(searchParams));

  return (
    <div className="mx-auto max-w-6xl p-10">
      <DepartureSearchForm
        searchDepartures={() => searchDepartures()}
        isLoading={isFetching}
      />
      <DeparturesTable data={data} isLoading={isFetching} />
    </div>
  );
}
