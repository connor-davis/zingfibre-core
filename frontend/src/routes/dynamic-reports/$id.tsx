import { ErrorComponent, Link, createFileRoute } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';

import {
  type DynamicQuery,
  type ErrorResponse,
  getApiDynamicQueriesById,
} from '@/api-client';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
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
      dynamicQuery: data?.data,
    } as {
      dynamicQuery?: DynamicQuery;
    };
  },
});

function RouteComponent() {
  const { dynamicQuery } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/dynamic-reports">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">{dynamicQuery?.Name}</Label>
        </div>
      </div>
    </div>
  );
}
