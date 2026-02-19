import {
  getApiDynamicQueriesByIdResultsOptions,
  putApiDynamicQueriesByIdMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation, useQuery } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
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
import { ArrowLeftIcon } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { usePapaParse } from 'react-papaparse';

import { type ParseResult } from 'papaparse';
import { toast } from 'sonner';

import {
  type DynamicQuery,
  type DynamicQueryResult,
  type ErrorResponse,
  type UpdateDynamicQuery,
  getApiDynamicQueriesById,
} from '@/api-client';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { DataTablePagination } from '@/components/ui/data-table';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Label } from '@/components/ui/label';
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/dynamic-reports/$id')({
  component: () => <RouteComponent />,
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading dynamic report...</Label>
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
  loader: async ({ params }) => {
    const { id } = params;

    const { data } = await getApiDynamicQueriesById({
      client: apiClient,
      path: { id },
      throwOnError: true,
    });

    return {
      dynamicQuery: (data?.data ?? {}) as DynamicQuery,
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { readString } = usePapaParse();

  const [columns, setColumns] = useState<string[]>([]);
  const [data, setData] = useState<Record<string, string>[]>([]);

  const { id } = Route.useParams();
  const { dynamicQuery } = Route.useLoaderData();
  const {
    data: dynamicQueryResults,
    isLoading: isLoadingDynamicQueryResults,
    isError: isDynamicQueryResultsError,
  } = useQuery({
    ...getApiDynamicQueriesByIdResultsOptions({
      client: apiClient,
      path: {
        id,
      },
    }),
    enabled: dynamicQuery?.Status === 'complete',
  });

  const abortControllerRef = useRef<AbortController | null>(null);

  const [status, setStatus] = useState(dynamicQuery.Status);
  const [statusDetail, setStatusDetail] = useState(
    'The dynamic query is being generated.'
  );

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const updateForm = useForm<UpdateDynamicQuery>({
    defaultValues: {
      Prompt: dynamicQuery?.Prompt,
      Status: 'in_progress',
    },
  });

  const updateDynamicQuery = useMutation({
    ...putApiDynamicQueriesByIdMutation({
      client: apiClient,
      path: {
        id,
      },
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The dynamic query has been updated.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  const generateDynamicQuery = async () => {
    setStatus('in_progress');

    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    abortControllerRef.current = new AbortController();

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/api/dynamic-queries/${id}/generate`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          signal: abortControllerRef.current.signal,
          credentials: 'include',
        }
      );

      if (!response.ok || !response.body) {
        throw new Error(response.statusText);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;

        const parts = buffer.split('\n\n');

        buffer = parts.pop() || '';

        for (const part of parts) {
          const eventMatch = part.match(/^event: (.+)$/m);
          // const dataMatch = part.match(/^data: (.+)$/m);

          const eventType = eventMatch ? eventMatch[1].trim() : null;
          // const dataString = dataMatch ? dataMatch[1].trim() : null;
          // const dataValue = JSON.parse(dataString || '');

          if (eventType === 'done') {
            setStatusDetail(
              'The dynamic query has been generated successfully.'
            );

            setTimeout(() => {
              return router.invalidate();
            }, 1000);
          }
        }
      }
    } catch (error: unknown) {
      if (error instanceof DOMException && error.name === 'AbortError') {
        console.log('Request was aborted');

        return;
      }
    } finally {
      abortControllerRef.current = null;
    }
  };

  useEffect(() => {
    const disposeable = setTimeout(() => {
      if (status === 'in_progress') {
        generateDynamicQuery();
      }
    }, 1000);

    return () => {
      clearTimeout(disposeable);

      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [status]);

  useEffect(() => {
    const disposeable = setTimeout(() => {
      const data = ((dynamicQueryResults?.data ?? {}) as { data: string }).data;

      readString(data, {
        worker: true,
        header: true,
        complete: ({ data }: ParseResult<Record<string, string>>) => {
          setData(data);
          setColumns(Object.keys(data[0] ?? {}));
        },
      });
    }, 0);

    return () => {
      clearTimeout(disposeable);
    };
  }, [dynamicQueryResults]);

  const table = useReactTable({
    data: data,
    columns: columns.map((column) => ({
      accessorKey: column,
      header: () => <Label>{column}</Label>,
      cell: ({ row }) => <Label>{row.getValue(column)}</Label>,
    })) as ColumnDef<Record<string, string>, string>[],
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

  if (isLoadingDynamicQueryResults) {
    return (
      <div className="flex flex-col w-full h-full items-center justify-center">
        <Label className="text-muted-foreground">
          Loading dynamic report results...
        </Label>
      </div>
    );
  }

  if (isDynamicQueryResultsError) {
    return (
      <div className="flex flex-col w-full h-full items-center justify-center">
        <Alert variant="destructive" className="w-full max-w-lg">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            There was an error loading the dynamic report results.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3 overflow-hidden">
      <ResizablePanelGroup direction="horizontal">
        <ResizablePanel defaultSize={80} minSize={50} maxSize={80}>
          <div className="flex flex-col w-full h-full gap-3">
            <div className="flex items-center justify-between w-full h-auto">
              <div className="flex items-center gap-3">
                <Link to="/dynamic-reports">
                  <Button variant="ghost" size="icon">
                    <ArrowLeftIcon className="size-4" />
                  </Button>
                </Link>

                <Label className="text-lg">{dynamicQuery?.Name}</Label>
              </div>

              <Button asChild>
                <a
                  href={`${import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:6173'}/api/dynamic-queries/${id}/export`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Export
                </a>
              </Button>
            </div>

            {status === 'in_progress' && (
              <div className="flex flex-col w-full h-full items-center justify-center gap-3">
                <div className="flex flex-col w-auto h-auto gap-0 items-center justify-center -rotate-60">
                  <Skeleton className="w-20 h-12 rounded-md" />
                  <Skeleton className="w-5 h-20 rounded-none" />
                  <Skeleton className="w-10 h-5 rounded-md" />
                </div>

                <Label className="text-muted-foreground">{statusDetail}</Label>
                <Label className="text-destructive/20">
                  Do not close this page.
                </Label>
              </div>
            )}

            {status === 'complete' && (
              <div className="flex flex-col w-full h-full gap-3 overflow-hidden">
                <div className="flex flex-col w-full h-auto border rounded-md bg-accent overflow-y-auto">
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
                            colSpan={
                              (
                                (dynamicQueryResults?.data ??
                                  {}) as DynamicQueryResult
                              ).data.length
                            }
                            className="h-24 text-center"
                          >
                            No results.
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </div>

                <DataTablePagination table={table} />
              </div>
            )}
          </div>
        </ResizablePanel>

        <ResizableHandle className="mx-3" />

        <ResizablePanel defaultSize={20} minSize={20} maxSize={50}>
          <Form {...updateForm}>
            <form
              onSubmit={updateForm.handleSubmit((values) =>
                updateDynamicQuery.mutate({
                  path: { id },
                  body: values,
                })
              )}
              className="space-y-6 w-full"
            >
              <FormField
                control={updateForm.control}
                name="Prompt"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Prompt</FormLabel>
                    <FormControl>
                      <Textarea placeholder="Prompt" {...field}></Textarea>
                    </FormControl>
                    <FormDescription>
                      Enter your dynamic query changes here...
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full">
                Tweak Query
              </Button>
            </form>
          </Form>
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
}
