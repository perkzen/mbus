import { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Autocomplete } from '@/components/ui/autocomplete';
import { Button } from '@/components/ui/button';
import { ArrowLeftRight } from 'lucide-react';
import { DateInput } from '@/components/ui/date-input';

const DepartureSearchForm = () => {
  const [fromStation, setFromStation] = useState('');
  const [toStation, setToStation] = useState('');
  const [_selectedDate, _setSelectedDate] = useState('05.07.2025');

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
              <Autocomplete
                options={[]}
                value={fromStation}
                onChange={setFromStation}
                className={'w-full'}
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
              <Autocomplete
                options={[]}
                value={toStation}
                onChange={setToStation}
                className={'w-full'}
              />
            </div>
          </div>

          <div className="w-full md:w-fit">
            <label className="mb-2 block text-sm font-medium text-gray-700">
              IZBERITE DATUM
            </label>
            <DateInput />
          </div>

          <Button className="bg-primary w-full text-white md:w-fit">
            IŠČI VOZNI RED
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default DepartureSearchForm;
