import {
  getApiAuthenticationCheckQueryKey,
  postApiAuthenticationMfaVerifyMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';

import { REGEXP_ONLY_DIGITS } from 'input-otp';
import { toast } from 'sonner';

import type { ErrorResponse } from '@/api-client';
import { useAuthentication } from '@/components/providers/authentication-provider';
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/utils';

export default function VerifyMfaForm() {
  const queryClient = useQueryClient();

  const { user, isLoading } = useAuthentication();

  if (isLoading) return null;
  if (!user) return null;

  const [code, setCode] = useState<string>('');

  const verifyMfa = useMutation({
    ...postApiAuthenticationMfaVerifyMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'MFA has been verified.',
        duration: 2000,
      });

      return queryClient.invalidateQueries({
        queryKey: getApiAuthenticationCheckQueryKey({
          client: apiClient,
        }),
      });
    },
  });

  return (
    <div className="flex flex-col items-center justify-center w-screen h-screen p-3">
      <div className="flex flex-col w-full md:max-w-96 items-center justify-center gap-5 md:gap-10 p-5 md:p-10 m-5 md:m-10 border rounded-md bg-popover">
        <div className="flex flex-col space-y-3 text-center items-center justify-center">
          <img
            src="/zing-logo.png"
            alt="Zing Logo"
            className="w-52 dark:hidden"
          />

          <img
            src="/zing-logo-dark.png"
            alt="Zing Logo"
            className="w-52 hidden dark:block"
          />

          <Label className="text-sm">Welcome back, {user.Email}!</Label>

          <Label className="text-xs text-muted-foreground">
            Please enter your 6-digit MFA code to continue.
          </Label>
        </div>

        <div className="flex flex-col w-full h-auto space-y-3">
          <InputOTP
            maxLength={6}
            pattern={REGEXP_ONLY_DIGITS}
            onChange={(value: string) => /^\d+$/g.test(value) && setCode(value)}
            onComplete={() => verifyMfa.mutate({ body: { code } })}
          >
            <InputOTPGroup>
              <InputOTPSlot
                index={0}
                className="bg-background text-foreground"
              />
              <InputOTPSlot
                index={1}
                className="bg-background text-foreground"
              />
              <InputOTPSlot
                index={2}
                className="bg-background text-foreground"
              />
            </InputOTPGroup>
            <InputOTPSeparator />
            <InputOTPGroup>
              <InputOTPSlot
                index={3}
                className="bg-background text-foreground"
              />
              <InputOTPSlot
                index={4}
                className="bg-background text-foreground"
              />
              <InputOTPSlot
                index={5}
                className="bg-background text-foreground"
              />
            </InputOTPGroup>
          </InputOTP>
        </div>
      </div>
    </div>
  );
}
