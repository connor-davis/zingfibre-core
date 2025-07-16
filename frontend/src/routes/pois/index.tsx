import {
  Link,
  Navigate,
  createFileRoute,
  useRouter,
  useSearch,
} from '@tanstack/react-router';

import { z } from 'zod';

import { type PointOfInterest, getApiPois } from '@/api-client';
import CreatePoiDialog from '@/components/dialogs/pois/create';
import DeletePoiDialog from '@/components/dialogs/pois/delete';
import RoleGuard from '@/components/guards/role-guard';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/pois/')({
  component: () => <RouteComponent />,
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  loaderDeps: ({ search: { page, search } }) => ({ page, search }),
  loader: async ({ deps: { page, search } }) => {
    const { data } = await getApiPois({
      client: apiClient,
      query: {
        page,
        search,
      },
    });

    return {
      ...data,
      pois: data?.data,
      pages: data ? (data.pages ? data.pages : 1) : 1,
    } as {
      pois?: PointOfInterest[];
      pages?: number;
    };
  },
});

function RouteComponent() {
  const router = useRouter();

  const { page } = useSearch({ from: '/pois/' });
  const { pois, pages } = Route.useLoaderData();

  if (pages && pages !== 0 && page > pages) {
    return <Navigate to="/pois" search={{ page: pages }} />;
  }

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Points of Interest</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search points of interest..."
            className="w-64"
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/pois',
                search: {
                  page,
                  search,
                },
              });
            }}
          />

          <RoleGuard value={['admin', 'staff']}>
            <CreatePoiDialog>
              <Button variant="ghost">Add</Button>
            </CreatePoiDialog>
          </RoleGuard>
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {pois?.length ? (
          pois.map((poi, index) => (
            <div
              key={poi.ID}
              className={cn(
                'flex items-center justify-between p-3',
                index + 1 < pois.length ? 'border-b' : ''
              )}
            >
              <div className="flex items-center gap-3">
                <div className="flex flex-col">
                  <Label className="text-sm">{poi.Name}</Label>
                  <Label className="text-xs text-muted-foreground">
                    {poi.Key}
                  </Label>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <RoleGuard value={['admin', 'staff']}>
                  <Link to={`/pois/$id`} params={{ id: poi.ID! }}>
                    <Button variant="ghost">Edit</Button>
                  </Link>
                </RoleGuard>

                <RoleGuard value="admin">
                  <DeletePoiDialog id={poi.ID!} poiKey={poi.Key!}>
                    <Button variant="ghost">Delete</Button>
                  </DeletePoiDialog>
                </RoleGuard>
              </div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No points of interest found.
            </Label>
          </div>
        )}
      </div>

      {pages && (
        <div className="flex items-center justify-end w-full p-3">
          <Label className="text-xs text-muted-foreground">
            Page {page} of {pages}
          </Label>

          <Link to="/pois" search={{ page: page - 1 }} disabled={page === 1}>
            <Button variant="outline" className="ml-3" disabled={page === 1}>
              Previous
            </Button>
          </Link>
          <Link
            to="/pois"
            search={{ page: page + 1 }}
            disabled={page === pages}
          >
            <Button
              variant="outline"
              className="ml-1"
              disabled={page === pages}
            >
              Next
            </Button>
          </Link>
        </div>
      )}
    </div>
  );
}
