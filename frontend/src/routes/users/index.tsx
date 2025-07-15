import {
  Link,
  Navigate,
  createFileRoute,
  useRouter,
  useSearch,
} from '@tanstack/react-router';

import { capitalCase } from 'change-case';
import { z } from 'zod';

import { type User, getApiUsers } from '@/api-client';
import CreateUserDialog from '@/components/dialogs/users/create';
import DeleteUserDialog from '@/components/dialogs/users/delete';
import RoleGuard from '@/components/guards/role-guard';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/users/')({
  component: () => <RouteComponent />,
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  loaderDeps: ({ search: { page, search } }) => ({ page, search }),
  loader: async ({ deps: { page, search } }) => {
    const { data } = await getApiUsers({
      client: apiClient,
      query: {
        page,
        search,
      },
    });

    return {
      ...data,
      users: data?.data as User[],
      pages: data ? (data.pages ? data.pages : 1) : 1,
    } as {
      users?: User[];
      pages?: number;
    };
  },
});

function RouteComponent() {
  const router = useRouter();

  const { page } = useSearch({ from: '/users/' });
  const { users, pages } = Route.useLoaderData();

  if (pages && pages !== 0 && page > pages) {
    return <Navigate to="/users" search={{ page: pages }} />;
  }

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Users</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search users..."
            className="w-64"
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/users',
                search: {
                  page,
                  search,
                },
              });
            }}
          />

          <RoleGuard value="admin">
            <CreateUserDialog>
              <Button variant="ghost">Add</Button>
            </CreateUserDialog>
          </RoleGuard>
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {users?.length ? (
          users.map((user, index) => (
            <div
              key={user.ID}
              className={cn(
                'flex items-center justify-between p-3',
                index + 1 < users.length ? 'border-b' : ''
              )}
            >
              <div className="flex items-center gap-3">
                <div className="flex flex-col">
                  <Label className="text-sm">{user.Email}</Label>
                  <Label className="text-xs text-muted-foreground">
                    {capitalCase(
                      Array.isArray(user.Role)
                        ? user.Role.join(', ')
                        : (user.Role ?? 'user')
                    )}
                  </Label>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <RoleGuard value="admin">
                  <Link to="/users/$id" params={{ id: user.ID! }}>
                    <Button variant="ghost">Edit</Button>
                  </Link>
                  <DeleteUserDialog id={user.ID!} email={user.Email!}>
                    <Button variant="ghost">Delete</Button>
                  </DeleteUserDialog>
                </RoleGuard>
              </div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No users found.
            </Label>
          </div>
        )}
      </div>

      {pages && (
        <div className="flex items-center justify-end w-full p-3">
          <Label className="text-xs text-muted-foreground">
            Page {page} of {pages}
          </Label>

          <Link to="/users" search={{ page: page - 1 }} disabled={page === 1}>
            <Button variant="outline" className="ml-3" disabled={page === 1}>
              Previous
            </Button>
          </Link>
          <Link
            to="/users"
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
