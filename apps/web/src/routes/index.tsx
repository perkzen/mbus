import { createFileRoute } from '@tanstack/react-router';
import DepartureSearchForm from '@/components/blocks/departure-search-form.tsx';
import DeparturesTable from '@/components/blocks/departures-table.tsx';

export const Route = createFileRoute('/')({
  component: App,
});

function App() {
  return (
    <div className="mx-auto max-w-6xl p-10">
      <DepartureSearchForm />
      <DeparturesTable />
    </div>
  );
}
