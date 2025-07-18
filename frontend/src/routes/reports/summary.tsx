import {
  ErrorComponent,
  createFileRoute,
  useRouter,
  useRouterState,
} from '@tanstack/react-router';
import {
  type ColumnDef,
  type ColumnFiltersState,
  type SortingState,
  type VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table';
import { ArrowUpDown, ChevronDownIcon } from 'lucide-react';
import { useState } from 'react';

import { format, parseISO } from 'date-fns';
import z from 'zod';

import {
  type ErrorResponse,
  type ReportSummaries,
  type ReportSummary,
  getApiReportsSummary,
} from '@/api-client';
import Pagination from '@/components/pagination';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { DebounceNumberInput } from '@/components/ui/debounce-number-input';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Label } from '@/components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/reports/summary')({
  component: RouteComponent,
  validateSearch: z.object({
    poi: z.string().default(''),
    search: z.string().default(''),
    months: z.coerce.number().default(1),
    page: z.coerce.number().default(1),
    pageSize: z.coerce.number().default(10),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">
        Loading customers report...
      </Label>
    </div>
  ),
  errorComponent: ({ error }: { error: Error | ErrorResponse }) => {
    if ('error' in error) {
      // Render a custom error message
      return (
        <div className="flex flex-col w-full h-full items-center justify-center">
          <Alert variant="destructive" className="w-full max-w-lg">
            <AlertTitle>{error.error}</AlertTitle>
            <AlertDescription>{error.details}</AlertDescription>
          </Alert>
        </div>
      );
    }

    if ('name' in error) {
      return (
        <div className="flex flex-col w-full h-full items-center justify-center">
          <Alert variant="destructive" className="w-full max-w-lg">
            <AlertTitle>{error.name}</AlertTitle>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        </div>
      );
    }

    // Fallback to the default ErrorComponent
    return <ErrorComponent error={error} />;
  },
  wrapInSuspense: true,
  loaderDeps: ({ search: { poi, search, months, page, pageSize } }) => ({
    poi,
    search,
    months,
    page,
    pageSize,
  }),
  loader: async ({ deps: { poi, search, months, page, pageSize } }) => {
    const { data } = await getApiReportsSummary({
      client: apiClient,
      query: {
        poi,
        search,
        months,
        page,
        pageSize,
      },
      throwOnError: true,
    });

    return {
      summaries: data?.data,
      pages: data?.pages ?? 1,
    } as {
      summaries?: ReportSummaries;
      pages: number;
    };
  },
});

export const columns = [
  {
    id: 'Created On',
    accessorKey: 'DateCreated',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Date
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => (
      <div>{format(parseISO(row.getValue('Created On')), 'dd/MM/yyyy')}</div>
    ),
  },
  {
    id: 'Item',
    accessorKey: 'ItemName',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Item
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Item')}</div>,
  },
  {
    id: 'Radius Username',
    accessorKey: 'RadiusUsername',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Radius Username
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Radius Username')}</div>,
  },
  {
    id: 'Method',
    accessorKey: 'Method',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Method
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => {
      return <div>{row.getValue('Method')}</div>;
    },
  },
  {
    id: 'Amount',
    accessorKey: 'Amount',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Amount
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Amount')}</div>,
  },
  {
    id: 'Service Id',
    accessorKey: 'ServiceId',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Service Id
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Service Id')}</div>,
  },
  {
    id: 'Build Name',
    accessorKey: 'BuildName',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Build Name
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Build Name')}</div>,
  },
  {
    id: 'Build Type',
    accessorKey: 'BuildType',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Build Type
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Build Type')}</div>,
  },
] as ColumnDef<ReportSummary>[];

function RouteComponent() {
  const routerState = useRouterState();
  const router = useRouter();

  const { poi, search, months } = Route.useLoaderDeps();
  const { summaries, pages } = Route.useLoaderData();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const table = useReactTable({
    data: summaries ?? [],
    columns,
    manualPagination: true,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3 overflow-hidden">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Summary Report</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            className="w-64"
            placeholder="Search for summary"
            defaultValue={search}
            onChange={(value) => {
              router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  search: value.target.value,
                }),
              });
            }}
          />

          <DebounceNumberInput
            className="w-24 h-9 rounded-r-none"
            min={1}
            max={100}
            value={months}
            onValueChange={(value) => {
              router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  months: value,
                }),
              });
            }}
          />

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="ml-auto">
                Columns <ChevronDownIcon />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {table
                .getAllColumns()
                .filter((column) => column.getCanHide())
                .map((column) => {
                  return (
                    <DropdownMenuCheckboxItem
                      key={column.id}
                      className="capitalize"
                      checked={column.getIsVisible()}
                      onCheckedChange={(value) =>
                        column.toggleVisibility(!!value)
                      }
                    >
                      {column.id}
                    </DropdownMenuCheckboxItem>
                  );
                })}
            </DropdownMenuContent>
          </DropdownMenu>

          <Button asChild>
            <a
              href={`${import.meta.env.VITE_API_URL ?? 'http://localhost:4000'}/api/exports/summary?poi=${poi}&months=${months}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Export
            </a>
          </Button>
        </div>
      </div>
      <div className="flex flex-col rounded-md border overflow-auto bg-accent w-full h-full">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={[...(summaries ?? [])].length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex flex-col w-full h-auto">
        <Pagination pages={pages} />
      </div>
    </div>
  );
}
