import { deleteApiDynamicQueriesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { type ReactNode, useState } from 'react';

import { toast } from 'sonner';

import type { ErrorResponse } from '@/api-client';
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

export default function DeleteDynamicQueryDialog({
  id,
  name,
  children,
}: {
  id: string;
  name: string;
  children: ReactNode;
}) {
  const router = useRouter();

  const [dynamicQueryName, setDynamicQueryName] = useState('');

  const deleteDynamicQuery = useMutation({
    ...deleteApiDynamicQueriesByIdMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) => {
      setDynamicQueryName('');

      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      });
    },
    onSuccess: () => {
      toast.success('Success', {
        description: 'The dynamic query has been deleted.',
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
          <DialogTitle>Delete Dynamic Report</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this dynamic report? This action
            cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col w-full h-auto gap-3">
          <Input
            type="text"
            value={dynamicQueryName}
            onChange={(e) => setDynamicQueryName(e.target.value)}
            placeholder="Enter name"
          />
          <Label className="text-sm text-muted-foreground">
            Please enter {name} to confirm.
          </Label>
        </div>

        <DialogFooter>
          <Button
            disabled={dynamicQueryName !== name}
            onClick={() => deleteDynamicQuery.mutate({ path: { id } })}
          >
            Continue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
