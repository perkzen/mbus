import type { FC } from 'react';
import type { Departure } from '@/api/departures.ts';
import { Loader2 } from 'lucide-react';

type DepartureListProps = {
  data?: Departure[];
  isLoading?: boolean;
};

const DepartureList: FC<DepartureListProps> = ({ data = [], isLoading }) => {
  const isEmpty = !isLoading && data.length === 0;

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 p-8 text-gray-500">
        <Loader2 className="h-6 w-6 animate-spin" />
        Nalagam vozni red...
      </div>
    );
  }

  if (isEmpty) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 p-10 text-center text-gray-500">
        <img
          src="/assets/lost.svg"
          alt="No departures"
          className="h-80 w-auto"
        />
        <p className="text-lg font-medium">Ni odhodov za prikaz.</p>
        <p className="text-sm text-gray-400">
          Poskusite z drugim datumom ali smerjo.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      {data.map((departure) => (
        <div
          key={`${departure.fromStation.name}-${departure.toStation.name}-${departure.departureAt}`}
          className="border p-4 shadow-sm"
        >
          <div className="grid grid-cols-[auto_min-content_1fr] items-center gap-x-2">
            <span className="font-semibold">{departure.departureAt}</span>
            <div className="bg-primary h-2 w-2 rounded-full" />
            <span>{departure.fromStation.name}</span>

            <span></span>
            <div className="h-8 w-0.5 justify-self-center bg-gray-300" />
            <span></span>

            <span className="font-semibold">{departure.arriveAt}</span>
            <div className="bg-primary h-2 w-2 rounded-full" />
            <span>{departure.toStation.name}</span>
          </div>

          <div className="mt-4 flex items-center justify-between text-sm">
            <div className="flex items-center gap-2">
              <div className="bg-primary text-background flex h-10 w-10 items-center justify-center">
                {departure.line}
              </div>
              <div>
                <div className="font-semibold">{departure.duration}</div>
                <div className="text-xs text-gray-500">
                  Prevoznik: Marprom d.o.o.
                </div>
              </div>
            </div>
            <div className="text-right">
              <div className="font-semibold">{departure.distance} km</div>
              <div className="font-semibold">1.3 EUR</div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default DepartureList;
