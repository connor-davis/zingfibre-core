import {
  postApiAuthenticationPasswordResetMutation,
  putApiUsersByIdMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useParams,
} from '@tanstack/react-router';
import { ArrowLeftIcon, ChevronsUpDown } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type ErrorResponse, type User, getApiUsersById } from '@/api-client';
import RoleGuard from '@/components/guards/role-guard';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
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
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp';
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
  loader: async ({ params }) => {
    const { id } = params;

    const { data } = await getApiUsersById({
      client: apiClient,
      path: { id },
      throwOnError: true,
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

  const passwordResetForm = useForm<{ Code: string; NewPassword: string }>({
    defaultValues: {
      Code: '',
      NewPassword: '',
    },
  });

  const updateUser = useMutation({
    ...putApiUsersByIdMutation({
      client: apiClient,
      path: { id },
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been updated.',
        duration: 2000,
      });
    },
  });

  const resetPassword = useMutation({
    ...postApiAuthenticationPasswordResetMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The password has been reset.',
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

      <RoleGuard value="admin">
        <Collapsible className="flex w-full flex-col gap-2">
          <CollapsibleTrigger asChild>
            <Button
              variant="outline"
              className="flex items-center justify-between gap-4 px-4"
            >
              <h4 className="text-sm font-semibold">Password Reset</h4>

              <ChevronsUpDown />
            </Button>
          </CollapsibleTrigger>
          <CollapsibleContent className="flex flex-col gap-2">
            <Form {...passwordResetForm}>
              <form
                onSubmit={passwordResetForm.handleSubmit((values) =>
                  resetPassword.mutate({
                    body: values,
                  })
                )}
                className="w-full grid grid-cols-1 md:grid-cols-2 gap-6 border bg-accent rounded-md p-3"
              >
                <FormField
                  control={passwordResetForm.control}
                  name="Code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>MFA Code</FormLabel>
                      <FormControl>
                        <InputOTP maxLength={6} {...field}>
                          <InputOTPGroup>
                            <InputOTPSlot index={0} />
                            <InputOTPSlot index={1} />
                            <InputOTPSlot index={2} />
                          </InputOTPGroup>
                          <InputOTPSeparator />
                          <InputOTPGroup>
                            <InputOTPSlot index={3} />
                            <InputOTPSlot index={4} />
                            <InputOTPSlot index={5} />
                          </InputOTPGroup>
                        </InputOTP>
                      </FormControl>
                      <FormDescription>
                        Please enter the MFA code from your authenticator app.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={passwordResetForm.control}
                  name="NewPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>New Password</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder="New Password"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Please enter the user's new password.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" className="w-full col-span-2">
                  Reset Password
                </Button>
              </form>
            </Form>
          </CollapsibleContent>
        </Collapsible>
      </RoleGuard>
    </div>
  );
}
