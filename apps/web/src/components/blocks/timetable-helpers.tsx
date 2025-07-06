import { Button } from '@/components/ui/button';
import { ArrowRight, ChevronLeft, ChevronRight } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

const TimetableDateNavigation = () => {
  return (
    <div className="mb-6 flex items-center justify-center gap-4">
      <Button variant="ghost" size="icon">
        <ChevronLeft className="h-4 w-4" />
      </Button>
      <div className="flex gap-2">
        <Button variant="ghost" className="text-gray-400">
          04.07.2025
        </Button>
        <Button className="bg-primary text-white">05.07.2025</Button>
        <Button variant="ghost" className="text-gray-400">
          06.07.2025
        </Button>
      </div>
      <Button variant="ghost" size="icon">
        <ChevronRight className="h-4 w-4" />
      </Button>
    </div>
  );
};

const Direction = () => {
  return (
    <div className="mb-6 flex items-center">
      <div className="flex items-center gap-2 text-2xl font-semibold">
        <span>Ihova</span>
        <ArrowRight className="h-6 w-6 text-gray-400" />
        <span>Maribor AP</span>
      </div>
    </div>
  );
};

const Timetable = () => {
  return (
    <div>
      <Direction />
      <TimetableDateNavigation />
      {/* Timetable */}
      <Card className="">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="">
                <tr>
                  <th className="p-4 text-left font-medium text-gray-700">
                    Odhod/prihod
                  </th>
                  <th className="p-4 text-left font-medium text-gray-700">
                    Trajanje
                  </th>
                  <th className="p-4 text-left font-medium text-gray-700">
                    Kilometri
                  </th>
                  <th className="p-4 text-left font-medium text-gray-700">
                    Cena
                  </th>
                  <th className="p-4 text-right font-medium"></th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-t">
                  <td className="p-4">
                    <div className="space-y-2">
                      <div className="flex items-center gap-3">
                        <span className="font-semibold">04:43</span>
                        <div className="bg-primary h-2 w-2 rounded-full"></div>
                        <span>Ihova</span>
                      </div>
                      <div className="ml-6 flex items-center gap-3">
                        <div className="h-8 w-2 border-l-2 border-gray-300"></div>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className="font-semibold">05:27</span>
                        <div className="bg-primary h-2 w-2 rounded-full"></div>
                        <span>Maribor AP</span>
                      </div>
                    </div>
                  </td>
                  <td className="p-4">
                    <div>
                      <div className="font-semibold">00:44</div>
                      <div className="text-sm text-gray-600">
                        Prevoznik: Arriva d.o.o.
                      </div>
                    </div>
                  </td>
                  <td className="p-4">
                    <span className="font-semibold">32 km</span>
                  </td>
                  <td className="p-4">
                    <span className="font-semibold">1.3 EUR</span>
                  </td>
                  <td className="p-4 text-right">
                    <Button
                      variant="outline"
                      className="text-primary border-primary bg-transparent"
                    >
                      PRIKAŽI POT
                    </Button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Timetable;
