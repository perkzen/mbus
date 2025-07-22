import { Loader2 } from 'lucide-react';
import type { Departure } from '@/api/departures.ts';
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import type { FC } from 'react';

const columnHelper = createColumnHelper<Departure>();

const columns = [
  columnHelper.accessor((row) => row, {
    id: 'id',
    header: 'Odhod/Prihod',
    cell: ({ getValue }) => {
      const row = getValue() as Departure;

      return (
        <div className="grid grid-cols-[auto_min-content_1fr] items-center gap-x-2">
          <span className="font-semibold">{row.departureAt}</span>
          <div className="bg-primary h-2 w-2 rounded-full" />
          <span>{row.fromStation.name}</span>

          <span></span>
          <div className="h-8 w-0.5 justify-self-center bg-gray-300" />
          <span></span>

          <span className="font-semibold">{row.arriveAt}</span>
          <div className="bg-primary h-2 w-2 rounded-full" />
          <span>{row.toStation.name}</span>
        </div>
      );
    },
  }),

  columnHelper.accessor((row) => row.line, {
    id: 'line',
    header: 'Linija',
    cell: (info) => (
      <div className="bg-primary text-background flex h-10 w-10 items-center justify-center">
        <div>{info.getValue()}</div>
      </div>
    ),
  }),

  columnHelper.accessor((row) => row.direction, {
    id: 'direction',
    header: 'Smer',
    cell: (info) => <div>{info.getValue()}</div>,
  }),

  columnHelper.accessor((row) => row.duration, {
    id: 'duration',
    header: 'Trajanje',
    cell: (info) => (
      <div>
        <div className="font-semibold">{info.getValue()}</div>
        <div className="text-sm text-gray-600">Prevoznik: Marprom d.o.o.</div>
      </div>
    ),
  }),

  columnHelper.accessor((row) => row.distance, {
    id: 'distance',
    header: 'Kilometri',
    cell: (info) => <span className="font-semibold">{info.getValue()} km</span>,
  }),
];

type TimetableTableProps = {
  data?: Departure[];
  isLoading?: boolean;
};

const DeparturesTable: FC<TimetableTableProps> = ({ data = [], isLoading }) => {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

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
          Poskusite z drugim datumom ali postajo.
        </p>
      </div>
    );
  }

  return (
    <div className="overflow-auto border">
      <table className="min-w-full table-auto text-sm">
        <thead className="bg-gray-100">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th className="p-4 text-left font-medium" key={header.id}>
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext()
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id} className="border-t">
              {row.getVisibleCells().map((cell) => (
                <td className="p-4" key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default DeparturesTable;
