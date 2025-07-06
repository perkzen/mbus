import { Button } from '@/components/ui/button';
import { ArrowRight, ChevronLeft, ChevronRight } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

const DepartureTableDateNavigation = () => {
  return (
    <div className="flex items-center justify-center gap-4 mb-6">
      <Button variant="ghost" size="icon">
        <ChevronLeft className="h-4 w-4" />
      </Button>
      <div className="flex gap-2">
        <Button variant="ghost" className="text-gray-400">
          04.07.2025
        </Button>
        <Button className="bg-orange-500 text-white border-b-4 border-orange-600">
          05.07.2025
        </Button>
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
    <div className="flex items-center mb-6">
      <div className="flex items-center gap-2 text-2xl font-semibold">
        <span>Ihova</span>
        <ArrowRight className="h-6 w-6 text-gray-400" />
        <span>Maribor AP</span>
      </div>
    </div>
  );
};

const DeparturesTable = () => {
  return (
    <div>
      <Direction />
      <DepartureTableDateNavigation />
      {/* Timetable */}
      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="text-left p-4 font-medium text-gray-700">
                    Odhod/prihod
                  </th>
                  <th className="text-left p-4 font-medium text-gray-700">
                    Trajanje
                  </th>
                  <th className="text-left p-4 font-medium text-gray-700">
                    Kilometri
                  </th>
                  <th className="text-left p-4 font-medium text-gray-700">
                    Cena
                  </th>
                  <th className="text-right p-4 font-medium"></th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-t">
                  <td className="p-4">
                    <div className="space-y-2">
                      <div className="flex items-center gap-3">
                        <span className="font-semibold">04:43</span>
                        <div className="w-2 h-2 bg-orange-500 rounded-full"></div>
                        <span>Ihova</span>
                      </div>
                      <div className="flex items-center gap-3 ml-6">
                        <div className="w-2 h-8 border-l-2 border-gray-300"></div>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className="font-semibold">05:27</span>
                        <div className="w-2 h-2 bg-orange-500 rounded-full"></div>
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
                      className="text-orange-500 border-orange-500 bg-transparent"
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

export default DeparturesTable;
