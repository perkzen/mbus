import type { ReactNode } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowLeftRight } from 'lucide-react';
import { DateInput } from '@/components/ui/date-input';
import { type DepartureSearchParams, Route } from '@/routes';
import { format } from 'date-fns';
import BusStationCombobox from '@/components/blocks/bus-station-combobox.tsx';

type DepartureSearchFormProps = {
  searchDepartures: () => void;
  isLoading: boolean;
};

const FormField = ({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) => (
  <div className="w-full">
    <label className="mb-2 block text-sm font-medium text-gray-700">
      {label}
    </label>
    {children}
  </div>
);

const DepartureSearchForm = ({
  searchDepartures,
  isLoading,
}: DepartureSearchFormProps) => {
  const { from, to, date } = Route.useSearch();
  const navigate = Route.useNavigate();

  const updateSearch = (newParams: Partial<DepartureSearchParams>) => {
    navigate({ search: { from, to, date, ...newParams } });
  };

  const swapStations = () => {
    updateSearch({ from: to, to: from });
  };

  return (
    <Card>
      <CardContent className="md:p-6">
        <div className="flex flex-col items-end gap-2 md:flex-row">
          <div className="flex w-full flex-col items-start gap-2 md:min-w-2/3 md:flex-row md:items-end">
            <FormField label="IZBERITE VSTOPNO POSTAJO">
              <BusStationCombobox
                key={from}
                value={from.toString()}
                onChange={(item) => updateSearch({ from: parseInt(item) })}
              />
            </FormField>

            <div className="flex justify-center">
              <Button
                variant="outline"
                size="icon"
                onClick={swapStations}
                className="text-primary bg-transparent"
              >
                <ArrowLeftRight className="h-4 w-4" />
              </Button>
            </div>

            <FormField label="IZBERITE IZSTOPNO POSTAJO">
              <BusStationCombobox
                key={to}
                value={to.toString()}
                onChange={(item) => updateSearch({ to: parseInt(item) })}
              />
            </FormField>
          </div>

          <FormField label="IZBERITE DATUM">
            <DateInput
              value={date ? new Date(date) : undefined}
              onChange={(newDate) => {
                if (newDate) {
                  const formatted = format(newDate, 'yyyy-MM-dd');
                  updateSearch({ date: formatted });
                }
              }}
            />
          </FormField>

          <Button
            className="bg-primary w-full text-white md:w-fit"
            onClick={searchDepartures}
            disabled={isLoading}
          >
            IŠČI VOZNI RED
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default DepartureSearchForm;
