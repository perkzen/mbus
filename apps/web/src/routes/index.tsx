import { createFileRoute } from '@tanstack/react-router'
import DepartureSearchForm from "@/components/blocks/departure-search-form.tsx";

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
      <div className="max-w-6xl mx-auto p-10">
        <DepartureSearchForm />
        {/*<DeparturesTable />*/}
      </div>
  )
}
