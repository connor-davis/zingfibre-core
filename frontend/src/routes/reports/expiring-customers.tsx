import { ErrorComponent, createFileRoute } from '@tanstack/react-router';
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
  type ReportExpiringCustomer,
  type ReportExpiringCustomers,
  getApiReportsExpiringCustomers,
} from '@/api-client';
import Pagination from '@/components/pagination';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
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

export const Route = createFileRoute('/reports/expiring-customers')({
  component: RouteComponent,
  validateSearch: z.object({
    poi: z.string().default(''),
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
  loaderDeps: ({ search: { poi, page, pageSize } }) => ({
    poi,
    page,
    pageSize,
  }),
  loader: async ({ deps: { poi, page, pageSize } }) => {
    const { data } = await getApiReportsExpiringCustomers({
      client: apiClient,
      query: {
        poi,
        page,
        pageSize,
      },
      throwOnError: true,
    });

    return {
      expiringCustomers: data?.data,
      pages: data?.pages ?? 1,
    } as {
      expiringCustomers?: ReportExpiringCustomers;
      pages: number;
    };
  },
});

export const columns = [
  {
    id: 'Expires On',
    accessorKey: 'Expiration',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Expires On
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => (
      <div>{format(parseISO(row.getValue('Expires On')), 'PPP')}</div>
    ),
  },
  {
    id: 'Full Name',
    accessorKey: 'FullName',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Full Name
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Full Name')}</div>,
    footer: (props) => props.column.id,
  },
  {
    id: 'Email',
    accessorKey: 'Email',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Email
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Email')}</div>,
  },
  {
    id: 'Username',
    accessorKey: 'RadiusUsername',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Username
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Username')}</div>,
  },
  {
    id: 'Phone Number',
    accessorKey: 'PhoneNumber',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Phone Number
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Phone Number')}</div>,
  },
  {
    id: 'Address',
    accessorKey: 'Address',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Address
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Address')}</div>,
  },
  {
    id: 'Last Duration',
    accessorKey: 'LastPurchaseDuration',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Last Duration
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Duration')}</div>,
  },
  {
    id: 'Last Speed',
    accessorKey: 'LastPurchaseSpeed',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Last Speed
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Speed')}</div>,
  },
] as ColumnDef<ReportExpiringCustomer>[];

function RouteComponent() {
  const { poi } = Route.useLoaderDeps();
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
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Expiring Customers Report</Label>
        </div>
        <div className="flex items-center gap-3">
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
              href={`${import.meta.env.VITE_API_URL ?? 'http://localhost:4000'}/api/exports/expiring-customers?poi=${poi}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Export
            </a>
          </Button>
        </div>
      </div>
      <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
        <div className="flex flex-col rounded-md border overflow-x-auto bg-accent">
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

        <Pagination pages={pages} />
      </div>
    </div>
  );
}
