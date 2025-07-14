import {
  getApiAuthenticationCheckQueryKey,
  postApiAuthenticationLoginMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { type ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type LoginRequest } from '@/api-client';
import EnableMfaForm from '@/components/authentication/enable-mfa-form';
import VerifyMfaForm from '@/components/authentication/verify-mfa-form';
import { useAuthentication } from '@/components/providers/authentication-provider';
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
import { apiClient } from '@/lib/utils';

export default function AuthenticationGuard({
  children,
  disabled = false,
}: {
  children: ReactNode;
  disabled?: boolean;
}) {
  const queryClient = useQueryClient();

  if (disabled) return children;

  const { user, isError, isLoading } = useAuthentication();

  const loginForm = useForm<LoginRequest>({
    defaultValues: {
      Email: '',
      Password: '',
    },
  });

  const loginUser = useMutation({
    ...postApiAuthenticationLoginMutation({
      client: apiClient,
    }),
    onError: (error) => {
      loginForm.reset();

      toast.error('Login failed', {
        description: error.message,
        duration: 2000,
      });
    },
    onSuccess: () => {
      loginForm.reset();

      toast.success('Login successful', {
        description: 'You have been logged in successfully.',
        duration: 2000,
      });

      return queryClient.invalidateQueries({
        queryKey: getApiAuthenticationCheckQueryKey({
          client: apiClient,
        }),
      });
    },
  });

  if (isLoading)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <Label className="text-muted-foreground">
          Checking authentication.
        </Label>
      </div>
    );

  if (isError || !user)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <div className="flex flex-col w-full md:max-w-96 items-center justify-center gap-10">
          <div className="flex flex-col w-full h-auto gap-5 items-center justify-center text-center">
            <img
              src="/zing-logo.png"
              alt="Zing Logo"
              className="w-52 dark:invert"
            />

            <p className="text-muted-foreground">
              You need to authenticate to continue to the application.
            </p>
          </div>

          <Form {...loginForm}>
            <form
              onSubmit={loginForm.handleSubmit((values) =>
                loginUser.mutate({
                  body: values,
                })
              )}
              className="space-y-6 w-full"
            >
              <FormField
                control={loginForm.control}
                name="Email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input placeholder="Email" {...field} />
                    </FormControl>
                    <FormDescription>
                      Please enter your email address.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={loginForm.control}
                name="Password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Password"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Please enter your password.
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
      </div>
    );

  if (!user.MfaEnabled) return <EnableMfaForm />;

  if (!user.MfaVerified) return <VerifyMfaForm />;

  return children;
}
