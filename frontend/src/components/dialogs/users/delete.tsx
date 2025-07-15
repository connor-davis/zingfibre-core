import { deleteApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { type ReactNode, useState } from 'react';

import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/utils';

export default function DeleteUserDialog({
  id,
  email,
  children,
}: {
  id: string;
  email: string;
  children: ReactNode;
}) {
  const router = useRouter();

  const [userEmail, setUserEmail] = useState('');

  const deleteUser = useMutation({
    ...deleteApiUsersByIdMutation({
      client: apiClient,
    }),
    onError: (error) => {
      setUserEmail('');

      toast.error('Failed', {
        description: error.message,
        duration: 2000,
      });
    },
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been deleted.',
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
          <DialogTitle>Delete User</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this user? This action cannot be
            undone.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col w-full h-auto gap-3">
          <Input
            type="email"
            value={userEmail}
            onChange={(e) => setUserEmail(e.target.value)}
            placeholder="Enter email"
          />
          <Label className="text-sm text-muted-foreground">
            Please enter {email} to confirm.
          </Label>
        </div>

        <DialogFooter>
          <Button
            disabled={userEmail !== email}
            onClick={() => deleteUser.mutate({ path: { id } })}
          >
            Continue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
