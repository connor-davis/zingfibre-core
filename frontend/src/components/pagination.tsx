import { useRouter, useRouterState, useSearch } from '@tanstack/react-router';
import { useEffect } from 'react';

import { Button } from './ui/button';
import { DebounceNumberInput } from './ui/debounce-number-input';
import { Label } from './ui/label';

export default function Pagination({ pages }: { pages: number }) {
  const routerState = useRouterState();
  const router = useRouter();
  const { page, pageSize } = useSearch({ strict: false });

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
    <div className="flex items-center justify-between gap-3">
      <div className="flex items-center gap-3">
        <Label className="text-xs text-muted-foreground">Showing</Label>

        <div>
          <DebounceNumberInput
            className="w-24 h-9 rounded-r-none"
            min={1}
            max={100}
            value={pageSize}
            onValueChange={(value) => {
              router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  pageSize: value,
                }),
              });
            }}
          />
        </div>

        <Label className="text-xs text-muted-foreground">items per page</Label>
      </div>

      <div className="flex items-center gap-3">
        <Label className="text-xs text-muted-foreground">
          Page {page} of {pages}
        </Label>

        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            router.navigate({
              to: routerState.location.pathname,
              search: (previous) => ({
                ...previous,
                page: (page ?? 1) - 1,
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
