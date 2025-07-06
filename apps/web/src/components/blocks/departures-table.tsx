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
        <div className="space-y-2">
          <div className="flex items-center gap-3">
            <span className="font-semibold">{row.departureAt}</span>
            <div className="bg-primary h-2 w-2 rounded-full"></div>
            <span>{row.fromStation.name}</span>
          </div>
          <div className="ml-6 flex items-center gap-3">
            <div className="h-8 w-2 border-l-2 border-gray-300"></div>
          </div>
          <div className="flex items-center gap-3">
            <span className="font-semibold">{row.arriveAt}</span>
            <div className="bg-primary h-2 w-2 rounded-full"></div>
            <span>{row.toStation.name}</span>
          </div>
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

  columnHelper.accessor((row) => row, {
    id: 'price',
    header: 'Cena',
    cell: () => <span className="font-semibold">1.3 EUR</span>,
  }),
];

type TimetableTableProps = {
  data: Departure[];
};

const DeparturesTable: FC<TimetableTableProps> = ({ data }) => {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="overflow-auto rounded border">
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
