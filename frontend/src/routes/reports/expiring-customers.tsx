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

import {
  type ErrorResponse,
  type ReportExpiringCustomer,
  type ReportExpiringCustomers,
  getApiReportsExpiringCustomers,
} from '@/api-client';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
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
  loaderDeps: ({ search: { poi } }) => ({ poi }),
  loader: async ({ deps: { poi } }) => {
    const { data } = await getApiReportsExpiringCustomers({
      client: apiClient,
      query: {
        poi,
      },
      throwOnError: true,
    });

    return {
      expiringCustomers: data?.data,
    } as {
      expiringCustomers?: ReportExpiringCustomers;
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
    cell: ({ row }) => <div>{row.getValue('Expires On')}</div>,
  },
  {
    id: 'First Name',
    accessorKey: 'FirstName',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="First Name"
            value={
              (table.getColumn('First Name')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('First Name')?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('First Name')}</div>,
    footer: (props) => props.column.id,
  },
  {
    id: 'Email',
    accessorKey: 'Email',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Email"
            value={(table.getColumn('Email')?.getFilterValue() as string) ?? ''}
            onChange={(event) =>
              table.getColumn('Email')?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Email')}</div>,
  },
  {
    id: 'Username',
    accessorKey: 'RadiusUsername',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Username"
            value={
              (table.getColumn('Username')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Username')?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Username')}</div>,
  },
  {
    id: 'Phone Number',
    accessorKey: 'PhoneNumber',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Phone Number"
            value={
              (table.getColumn('Phone Number')?.getFilterValue() as string) ??
              ''
            }
            onChange={(event) =>
              table
                .getColumn('Phone Number')
                ?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Phone Number')}</div>,
  },
  {
    id: 'Address',
    accessorKey: 'Address',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Address"
            value={
              (table.getColumn('Address')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Address')?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Address')}</div>,
  },
  {
    id: 'Last Duration',
    accessorKey: 'LastPurchaseDuration',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Last Duration"
            value={
              (table.getColumn('Last Duration')?.getFilterValue() as string) ??
              ''
            }
            onChange={(event) =>
              table
                .getColumn('Last Duration')
                ?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Duration')}</div>,
  },
  {
    id: 'Last Speed',
    accessorKey: 'LastPurchaseSpeed',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Last Speed"
            value={
              (table.getColumn('Last Speed')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Last Speed')?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
            className="p-0 hover:bg-transparent"
          >
            <ArrowUpDown className="w-4 h-4" />
          </Button>
        </div>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Speed')}</div>,
  },
] as ColumnDef<ReportExpiringCustomer>[];

function RouteComponent() {
  const { poi } = Route.useLoaderDeps();
  const { expiringCustomers } = Route.useLoaderData();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const table = useReactTable({
    data: expiringCustomers ?? [],
    columns,
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
      <div className="flex flex-col w-full h-full gap-3">
        <div className="w-full h-full overflow-hidden">
          <div className="rounded-md border overflow-x-auto bg-accent">
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
        </div>
        <div className="flex items-center justify-end py-4 space-x-2">
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              Next
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
