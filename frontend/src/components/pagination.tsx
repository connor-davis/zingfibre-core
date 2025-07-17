import { useRouter, useRouterState, useSearch } from '@tanstack/react-router';
import { useEffect } from 'react';

import { Button } from './ui/button';

export default function Pagination({ pages }: { pages: number }) {
  const routerState = useRouterState();
  const router = useRouter();
  const { page } = useSearch({ strict: false });

  useEffect(() => {
    const disposeable = setTimeout(() => {
      if (!page || isNaN(Number(page)) || Number(page) < 1) {
        router.navigate({
          to: routerState.location.pathname,
          search: (previous) => ({
            ...previous,
            page: 1,
          }),
        });
      }

      if (pages && page && Number(page) > pages) {
        router.navigate({
          to: routerState.location.pathname,
          search: (previous) => ({
            ...previous,
            page: pages,
          }),
        });
      }
    }, 0);

    return () => {
      clearTimeout(disposeable);
    };
  }, [routerState.location.pathname, page]);

  return (
    <div className="flex items-center justify-end py-4 space-x-2">
      <div className="flex items-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                page: Math.max(page ?? 1, 1),
              }),
            })
          }
          disabled={!page || page <= 1}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                page: (page ?? 1) + 1,
              }),
            })
          }
          disabled={!page || page >= pages}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
