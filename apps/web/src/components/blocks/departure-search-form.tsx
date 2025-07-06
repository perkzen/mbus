import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowLeftRight } from 'lucide-react';
import { DateInput } from '@/components/ui/date-input';
import type { ComboBoxItem } from '@/components/ui/api-combobox.tsx';
import { BusStationSelect } from '@/components/blocks/bus-station-select.tsx';

type DepartureSearchFormProps = {
  fromStation?: ComboBoxItem;
  setFromStation: (item?: ComboBoxItem) => void;
  toStation?: ComboBoxItem;
  setToStation: (item?: ComboBoxItem) => void;
  searchDepartures: () => void;
};

const DepartureSearchForm = ({
  fromStation,
  setFromStation,
  toStation,
  setToStation,
  searchDepartures,
}: DepartureSearchFormProps) => {
  const swapStations = () => {
    const temp = fromStation;
    setFromStation(toStation);
    setToStation(temp);
  };

  return (
    <Card className="mb-6">
      <CardContent className="p-6">
        <div className="flex flex-col items-end gap-2 md:flex-row">
          <div className="flex w-full flex-col items-start gap-2 md:max-w-2/3 md:flex-row md:items-end">
            <div className="w-full">
              <label className="mb-2 block text-sm font-medium text-gray-700">
                IZBERITE VSTOPNO POSTAJO
              </label>
              <BusStationSelect
                selectedItem={fromStation}
                onSelect={setFromStation}
                className="w-full"
              />
            </div>

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

            <div className="w-full">
              <label className="mb-2 block text-sm font-medium text-gray-700">
                IZBERITE IZSTOPNO POSTAJO
              </label>
              <BusStationSelect
                selectedItem={toStation}
                onSelect={setToStation}
                className="w-full"
              />
            </div>
          </div>

          <div className="w-full md:w-fit">
            <label className="mb-2 block text-sm font-medium text-gray-700">
              IZBERITE DATUM
            </label>
            <DateInput />
          </div>

          <Button
            className="bg-primary w-full text-white md:w-fit"
            onClick={() => searchDepartures()}
          >
            IŠČI VOZNI RED
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default DepartureSearchForm;
