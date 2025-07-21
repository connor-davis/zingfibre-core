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
import { CalendarIcon, ChevronDownIcon, FilterIcon } from 'lucide-react';
import { useState } from 'react';

import { format, parseISO, subDays } from 'date-fns';
import z from 'zod';

import {
  type ErrorResponse,
  type ReportRecharge,
  type ReportRecharges,
  getApiReportsRecharges,
} from '@/api-client';
import Pagination from '@/components/pagination';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { DebounceInput } from '@/components/ui/debounce-input';
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Label } from '@/components/ui/label';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/reports/recharges')({
  component: RouteComponent,
  validateSearch: z.object({
    poi: z.string().default(''),
    search: z.string().default(''),
    startDate: z.string().default(subDays(new Date(), 7).toISOString()),
    endDate: z.string().default(new Date().toISOString()),
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
  loaderDeps: ({
    search: { poi, search, startDate, endDate, page, pageSize },
  }) => ({
    poi,
    search,
    startDate,
    endDate,
    page,
    pageSize,
  }),
  loader: async ({
    deps: { poi, search, startDate, endDate, page, pageSize },
  }) => {
    const { data } = await getApiReportsRecharges({
      client: apiClient,
      query: {
        poi,
        search,
        startDate,
        endDate,
        page,
        pageSize,
      },
      throwOnError: true,
    });

    return {
      expiringCustomers: data?.data,
      pages: data?.pages ?? 1,
    } as {
      expiringCustomers?: ReportRecharges;
      pages: number;
    };
  },
});

export const columns = [
  {
    id: 'Created On',
    accessorKey: 'DateCreated',
    header: () => {
      return <div>Created On</div>;
    },
    cell: ({ row }) => (
      <div>{format(parseISO(row.getValue('Created On')), 'dd/MM/yyyy')}</div>
    ),
  },
  {
    id: 'Email',
    accessorKey: 'Email',
    header: () => {
      return <div>Email</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Email')}</div>,
  },
  {
    id: 'Full Name',
    accessorKey: 'FullName',
    header: () => {
      return <div>Full Name</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Full Name')}</div>,
  },
  {
    id: 'Item',
    accessorKey: 'ItemName',
    header: () => {
      return <div>Item</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Item')}</div>,
  },
  {
    id: 'Amount',
    accessorKey: 'Amount',
    header: () => {
      return <div>Amount</div>;
    },
    cell: ({ row }) => {
      return (
        <div>
          {row.getValue<number | undefined>('Amount')?.toLocaleString('en-ZA', {
            style: 'currency',
            currency: 'ZAR',
          }) ?? 'R 0'}
        </div>
      );
    },
  },
  {
    id: 'Method',
    accessorKey: 'Method',
    header: () => {
      return <div>Method</div>;
    },
    cell: ({ row }) => {
      return <div>{row.getValue('Method')}</div>;
    },
  },
  {
    id: 'Status',
    accessorKey: 'Successful',
    header: () => {
      return <div>Status</div>;
    },
    cell: ({ row }) => (
      <div>{row.getValue('Status') === true ? 'Success' : 'Failed'}</div>
    ),
  },
  {
    id: 'Service Id',
    accessorKey: 'ServiceId',
    header: () => {
      return <div>Service Id</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Service Id')}</div>,
  },
  {
    id: 'Build Name',
    accessorKey: 'BuildName',
    header: () => {
      return <div>Build Name</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Build Name')}</div>,
  },
  {
    id: 'Build Type',
    accessorKey: 'BuildType',
    header: () => {
      return <div>Build Type</div>;
    },
    cell: ({ row }) => <div>{row.getValue('Build Type')}</div>,
  },
] as ColumnDef<ReportRecharge>[];

function RouteComponent() {
  const routerState = useRouterState();
  const router = useRouter();

  const { poi, search, startDate, endDate } = Route.useLoaderDeps();
  const { expiringCustomers, pages } = Route.useLoaderData();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const table = useReactTable({
    data: expiringCustomers ?? [],
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
          <Label className="text-lg">Recharges Report</Label>
        </div>
        <div className="hidden lg:flex items-center gap-3">
          <DebounceInput
            type="text"
            className="w-64"
            placeholder="Search for recharge"
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

          <Popover>
            <PopoverTrigger>
              <Button
                id="date"
                variant="outline"
                className={cn(
                  'w-full justify-start text-left font-normal',
                  !startDate && !endDate && 'text-muted-foreground'
                )}
              >
                <CalendarIcon className="mr-2 size-4" />
                {startDate ? (
                  endDate ? (
                    <>
                      {format(parseISO(startDate), 'LLL dd, y')} -{' '}
                      {format(parseISO(endDate), 'LLL dd, y')}
                    </>
                  ) : (
                    format(parseISO(startDate), 'LLL dd, y')
                  )
                ) : (
                  <span>Pick a date</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
              <Calendar
                mode="range"
                defaultMonth={startDate ? parseISO(startDate) : undefined}
                selected={{
                  from: startDate ? parseISO(startDate) : undefined,
                  to: endDate ? parseISO(endDate) : undefined,
                }}
                onSelect={(selected) => {
                  if (!selected) {
                    return;
                  }

                  const from = selected.from
                    ? new Date(selected.from.setHours(0, 0, 0, 0)).toISOString()
                    : new Date(new Date().setHours(0, 0, 0, 0)).toISOString();
                  const to = selected.to
                    ? new Date(
                        selected.to.setHours(23, 59, 59, 999)
                      ).toISOString()
                    : new Date(
                        new Date().setHours(23, 59, 59, 999)
                      ).toISOString();

                  router.navigate({
                    to: '/reports/recharges',
                    search: (previous) => ({
                      ...previous,
                      startDate: from,
                      endDate: to,
                    }),
                  });
                }}
                numberOfMonths={2}
              />
            </PopoverContent>
          </Popover>

          <Button asChild>
            <a
              href={`${import.meta.env.VITE_API_URL ?? 'http://localhost:4000'}/api/exports/recharges?poi=${poi}&startDate=${startDate}&endDate=${endDate}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Export
            </a>
          </Button>
        </div>
        <Drawer>
          <DrawerTrigger asChild>
            <Button variant="ghost" size="icon" className="lg:hidden">
              <FilterIcon className="size-4" />
            </Button>
          </DrawerTrigger>
          <DrawerContent>
            <DrawerHeader>
              <DrawerTitle>Filters</DrawerTitle>
              <DrawerDescription>
                Use the filters to refine the data displayed in the report.
              </DrawerDescription>
            </DrawerHeader>

            <div className="flex flex-col gap-3 p-3">
              <DebounceInput
                type="text"
                className="w-full"
                placeholder="Search for recharge"
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

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="justify-between w-full">
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

              <Popover>
                <PopoverTrigger>
                  <Button
                    id="date"
                    variant="outline"
                    className={cn(
                      'w-full justify-start text-left font-normal',
                      !startDate && !endDate && 'text-muted-foreground'
                    )}
                  >
                    <CalendarIcon className="mr-2 size-4" />
                    {startDate ? (
                      endDate ? (
                        <>
                          {format(parseISO(startDate), 'LLL dd, y')} -{' '}
                          {format(parseISO(endDate), 'LLL dd, y')}
                        </>
                      ) : (
                        format(parseISO(startDate), 'LLL dd, y')
                      )
                    ) : (
                      <span>Pick a date</span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="range"
                    defaultMonth={startDate ? parseISO(startDate) : undefined}
                    selected={{
                      from: startDate ? parseISO(startDate) : undefined,
                      to: endDate ? parseISO(endDate) : undefined,
                    }}
                    onSelect={(selected) => {
                      if (!selected) {
                        return;
                      }

                      const from = selected.from
                        ? new Date(
                            selected.from.setHours(0, 0, 0, 0)
                          ).toISOString()
                        : new Date(
                            new Date().setHours(0, 0, 0, 0)
                          ).toISOString();
                      const to = selected.to
                        ? new Date(
                            selected.to.setHours(23, 59, 59, 999)
                          ).toISOString()
                        : new Date(
                            new Date().setHours(23, 59, 59, 999)
                          ).toISOString();

                      router.navigate({
                        to: '/reports/recharges',
                        search: (previous) => ({
                          ...previous,
                          startDate: from,
                          endDate: to,
                        }),
                      });
                    }}
                    numberOfMonths={2}
                  />
                </PopoverContent>
              </Popover>

              <Button asChild>
                <a
                  href={`${import.meta.env.VITE_API_URL ?? 'http://localhost:4000'}/api/exports/recharges?poi=${poi}&startDate=${startDate}&endDate=${endDate}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Export
                </a>
              </Button>
            </div>
          </DrawerContent>
        </Drawer>
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
                  colSpan={[...(expiringCustomers ?? [])].length}
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
