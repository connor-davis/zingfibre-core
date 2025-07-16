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
  type ReportRechargeSummaries,
  type ReportRechargeSummary,
  getApiReportsRechargesSummary,
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

export const Route = createFileRoute('/reports/recharges-summary')({
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
    const { data } = await getApiReportsRechargesSummary({
      client: apiClient,
      query: {
        poi,
      },
      throwOnError: true,
    });

    return {
      expiringCustomers: data?.data,
    } as {
      expiringCustomers?: ReportRechargeSummaries;
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
    cell: ({ row }) => <div>{row.getValue('Created On')}</div>,
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
    id: 'Surname',
    accessorKey: 'Surname',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Surname"
            value={
              (table.getColumn('Surname')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Surname')?.setFilterValue(event.target.value)
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
    cell: ({ row }) => <div>{row.getValue('Surname')}</div>,
    footer: (props) => props.column.id,
  },
  {
    id: 'Item',
    accessorKey: 'ItemName',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Item"
            value={(table.getColumn('Item')?.getFilterValue() as string) ?? ''}
            onChange={(event) =>
              table.getColumn('Item')?.setFilterValue(event.target.value)
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
    cell: ({ row }) => <div>{row.getValue('Item')}</div>,
  },
  {
    id: 'Amount',
    accessorKey: 'Amount',
    header: ({ table, column }) => {
      return (
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Amount"
            value={
              (table.getColumn('Amount')?.getFilterValue() as string) ?? ''
            }
            onChange={(event) =>
              table.getColumn('Amount')?.setFilterValue(event.target.value)
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
    cell: ({ row }) => {
      return <div>R {row.getValue('Amount')}</div>;
    },
  },
  {
    id: 'Status',
    accessorKey: 'Successful',
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          className="p-0 hover:bg-transparent"
        >
          Status
          <ArrowUpDown className="w-4 h-4 ml-2" />
        </Button>
      );
    },
    cell: ({ row }) => (
      <div>{row.getValue('Status') === true ? 'Success' : 'Failed'}</div>
    ),
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
] as ColumnDef<ReportRechargeSummary>[];

function RouteComponent() {
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
          <Label className="text-lg">Recharges Summary Report</Label>
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
