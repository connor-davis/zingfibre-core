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
  type ReportCustomers,
  getApiReportsCustomers,
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

export const Route = createFileRoute('/reports/customers')({
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
    const { data } = await getApiReportsCustomers({
      client: apiClient,
      query: {
        poi,
      },
      throwOnError: true,
    });

    return {
      customers: data?.data,
    } as {
      customers?: ReportCustomers[];
    };
  },
});

export const columns = [
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
    id: 'Last Name',
    accessorKey: 'Surname',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Last Name"
            value={
              (table.getColumn('Last Name')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Last Name')?.setFilterValue(event.target.value)
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
    cell: ({ row }) => <div>{row.getValue('Last Name')}</div>,
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
    id: 'Radius Username',
    accessorKey: 'RadiusUsername',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Radius Username"
            value={
              (table
                .getColumn('Radius Username')
                ?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table
                .getColumn('Radius Username')
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
    cell: ({ row }) => <div>{row.getValue('Radius Username')}</div>,
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
] as ColumnDef<ReportCustomers>[];

function RouteComponent() {
  const { poi } = Route.useLoaderDeps();
  const { customers } = Route.useLoaderData();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const table = useReactTable({
    data: customers ?? [],
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
          <Label className="text-lg">Customers Report</Label>
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
              href={`${import.meta.env.VITE_API_URL ?? 'http://localhost:4000'}/api/exports/customers?poi=${poi}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Export
            </a>
          </Button>
        </div>
      </div>
      <div className="flex flex-col w-full h-full gap-3">
        <div className="w-full overflow-hidden">
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
                      colSpan={[...(customers ?? [])].length}
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
