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
import { ChevronDownIcon, ChevronUpIcon, FilterIcon } from 'lucide-react';
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
    search: z.string().default(''),
    page: z.coerce.number().default(1),
    pageSize: z.coerce.number().default(10),
    sort: z.string().default('expiration_asc'),
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
  loaderDeps: ({ search: { poi, search, page, pageSize, sort } }) => ({
    poi,
    search,
    page,
    pageSize,
    sort,
  }),
  loader: async ({ deps: { poi, search, page, pageSize, sort } }) => {
    const { data } = await getApiReportsExpiringCustomers({
      client: apiClient,
      query: {
        poi,
        search,
        page,
        pageSize,
        sort,
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
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('expiration')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'expiration_asc'
                    : 'expiration_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'expiration_asc'
                    ? 'expiration_desc'
                    : 'expiration_asc',
              }),
            });
          }}
        >
          Expires On
          {sort.startsWith('expiration') ? (
            sort === 'expiration_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => (
      <div>{format(parseISO(row.getValue('Expires On')), 'dd/MM/yyyy')}</div>
    ),
  },
  {
    id: 'Full Name',
    accessorKey: 'FullName',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('full_name')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'full_name_asc'
                    : 'full_name_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'full_name_asc' ? 'full_name_desc' : 'full_name_asc',
              }),
            });
          }}
        >
          Full Name
          {sort.startsWith('full_name') ? (
            sort === 'full_name_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Full Name')}</div>,
  },
  {
    id: 'Email',
    accessorKey: 'Email',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('email')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc') ? 'email_asc' : 'email_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort: sort === 'email_asc' ? 'email_desc' : 'email_asc',
              }),
            });
          }}
        >
          Email
          {sort.startsWith('email') ? (
            sort === 'email_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Email')}</div>,
  },
  {
    id: 'Radius Username',
    accessorKey: 'RadiusUsername',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('radius_username')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'radius_username_asc'
                    : 'radius_username_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'radius_username_asc'
                    ? 'radius_username_desc'
                    : 'radius_username_asc',
              }),
            });
          }}
        >
          Radius Username
          {sort.startsWith('radius_username') ? (
            sort === 'radius_username_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Radius Username')}</div>,
  },
  {
    id: 'Phone Number',
    accessorKey: 'PhoneNumber',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('phone_number')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'phone_number_asc'
                    : 'phone_number_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'phone_number_asc'
                    ? 'phone_number_desc'
                    : 'phone_number_asc',
              }),
            });
          }}
        >
          Phone Number
          {sort.startsWith('phone_number') ? (
            sort === 'phone_number_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Phone Number')}</div>,
  },
  {
    id: 'Address',
    accessorKey: 'Address',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('address')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc') ? 'address_asc' : 'address_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort: sort === 'address_asc' ? 'address_desc' : 'address_asc',
              }),
            });
          }}
        >
          Address
          {sort.startsWith('address') ? (
            sort === 'address_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Address')}</div>,
  },
  {
    id: 'Last Duration',
    accessorKey: 'LastPurchaseDuration',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('last_purchase_duration')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'last_purchase_duration_asc'
                    : 'last_purchase_duration_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'last_purchase_duration_asc'
                    ? 'last_purchase_duration_desc'
                    : 'last_purchase_duration_asc',
              }),
            });
          }}
        >
          Last Duration
          {sort.startsWith('last_purchase_duration') ? (
            sort === 'last_purchase_duration_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Duration')}</div>,
  },
  {
    id: 'Last Speed',
    accessorKey: 'LastPurchaseSpeed',
    header: () => {
      const routerState = useRouterState();
      const router = useRouter();

      const { sort } = Route.useLoaderDeps();

      return (
        <Button
          variant={'ghost'}
          onClick={() => {
            if (!sort.startsWith('last_purchase_speed')) {
              return router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  sort: sort.endsWith('asc')
                    ? 'last_purchase_speed_asc'
                    : 'last_purchase_speed_desc',
                }),
              });
            }

            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                sort:
                  sort === 'last_purchase_speed_asc'
                    ? 'last_purchase_speed_desc'
                    : 'last_purchase_speed_asc',
              }),
            });
          }}
        >
          Last Speed
          {sort.startsWith('last_purchase_speed') ? (
            sort === 'last_purchase_speed_asc' ? (
              <ChevronUpIcon className="size-4" />
            ) : (
              <ChevronDownIcon className="size-4" />
            )
          ) : (
            ''
          )}
        </Button>
      );
    },
    cell: ({ row }) => <div>{row.getValue('Last Speed')}</div>,
  },
] as ColumnDef<ReportExpiringCustomer>[];

function RouteComponent() {
  const routerState = useRouterState();
  const router = useRouter();

  const { poi, search } = Route.useLoaderDeps();
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
          <Label className="text-lg">Expiring Customers Report</Label>
        </div>
        <div className="hidden lg:flex items-center gap-3">
          <DebounceInput
            type="text"
            className="w-64"
            placeholder="Search for expiring customer"
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
                placeholder="Search for expiring customer"
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
