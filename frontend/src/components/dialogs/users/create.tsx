import { postApiAuthenticationRegisterMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import type { ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type CreateUser } from '@/api-client';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/utils';

export default function CreateUserDialog({
  children,
}: {
  children: ReactNode;
}) {
  const router = useRouter();

  const createUserForm = useForm<CreateUser>({
    defaultValues: {
      Email: undefined,
      Password: undefined,
      Role: undefined,
    },
  });

  const createUser = useMutation({
    ...postApiAuthenticationRegisterMutation({
      client: apiClient,
    }),
    onError: (error) => {
      toast.error('Failed', {
        description: error.message,
        duration: 2000,
      });
    },
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been created successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <Dialog>
      <DialogTrigger>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create User</DialogTitle>
          <DialogDescription>
            Use this form to create a new user account. Please fill in all
            required fields.
          </DialogDescription>
        </DialogHeader>

        <Form {...createUserForm}>
          <form
            onSubmit={createUserForm.handleSubmit((values) =>
              createUser.mutate({
                body: values,
              })
            )}
          >
            <FormField
              control={createUserForm.control}
              name="Email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" placeholder="Email" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter the user's email address.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createUserForm.control}
              name="Password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="Password" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter the user's password.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createUserForm.control}
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
      </DialogContent>
    </Dialog>
  );
}
