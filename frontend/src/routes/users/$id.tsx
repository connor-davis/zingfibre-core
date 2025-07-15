import { putApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useParams } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type User, type ZingResponse, getApiUsersById } from '@/api-client';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/users/$id')({
  component: () => <RouteComponent />,
  loader: async ({ params }) => {
    const { id } = params;

    const { data } = await getApiUsersById({
      client: apiClient,
      path: { id },
    });

    return {
      user: data?.data,
    } as {
      user?: User;
    };
  },
});

function RouteComponent() {
  const { id } = useParams({ from: '/users/$id' });
  const { user } = Route.useLoaderData();

  const userForm = useForm<User>({
    values: user,
  });

  const updateUser = useMutation({
    ...putApiUsersByIdMutation({
      client: apiClient,
      path: { id },
    }),
    onError: (error: ZingResponse) => {
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      });
    },
    onSuccess: ({ data: updatedUser }) => {
      userForm.reset({ ...user, ...updatedUser });

      toast.success('Success', {
        description: 'The user has been updated.',
        duration: 2000,
      });
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/users">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Edit User</Label>
        </div>
      </div>

      <Form {...userForm}>
        <form
          onSubmit={userForm.handleSubmit((values) =>
            updateUser.mutate({
              path: { id },
              body: values,
            })
          )}
          className="space-y-6 w-full"
        >
          <FormField
            control={userForm.control}
            name="Email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input placeholder="Email" {...field} />
                </FormControl>
                <FormDescription>
                  Please enter the user's email address.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={userForm.control}
            name="Role"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Role</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={
                    Array.isArray(field.value)
                      ? field.value.join(', ')
                      : (field.value ?? 'user')
                  }
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a role" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="admin">Admin</SelectItem>
                    <SelectItem value="staff">Staff</SelectItem>
                    <SelectItem value="user">User</SelectItem>
                  </SelectContent>
                </Select>
                <FormDescription>
                  Please select the user's role.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" className="w-full">
            Continue
          </Button>
        </form>
      </Form>
    </div>
  );
}
