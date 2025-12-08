import { ErrorComponent, Link, createFileRoute } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';

import {
  type DynamicQuery,
  type ErrorResponse,
  getApiDynamicQueriesById,
} from '@/api-client';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
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
  const { id } = Route.useParams();
  const { dynamicQuery } = Route.useLoaderData();

  const abortControllerRef = useRef<AbortController | null>(null);

  const [status, setStatus] = useState(dynamicQuery?.Status ?? 'in_progress');
  const [statusDetail, setStatusDetail] = useState(
    'The dynamic query is being generated.'
  );

  const generateDynamicQuery = async () => {
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
          const dataMatch = part.match(/^data: (.+)$/m);

          const eventType = eventMatch ? eventMatch[1].trim() : null;
          const dataString = dataMatch ? dataMatch[1].trim() : null;
          const dataValue = JSON.parse(dataString || '');

          console.log(eventType, dataValue);
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
  }, []);

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

      {status === 'in_progress' && (
        <div className="flex flex-col w-full h-full items-center justify-center gap-3">
          <div className="flex flex-col w-auto h-auto gap-0 items-center justify-center -rotate-60">
            <Skeleton className="w-20 h-12 rounded-md" />
            <Skeleton className="w-5 h-20 rounded-none" />
            <Skeleton className="w-10 h-5 rounded-md" />
          </div>

          <Label className="text-muted-foreground">{statusDetail}</Label>
          <Label className="text-destructive/20">Do not close this page.</Label>
        </div>
      )}
    </div>
  );
}
