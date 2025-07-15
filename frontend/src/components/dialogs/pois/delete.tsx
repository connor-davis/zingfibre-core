import { deleteApiPoisByIdMutation } from '@/api-client/@tanstack/react-query.gen';
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

export default function DeletePoiDialog({
  id,
  key,
  children,
}: {
  id: string;
  key: string;
  children: ReactNode;
}) {
  const router = useRouter();

  const [poiKey, setPoiKey] = useState('');

  const deletePoi = useMutation({
    ...deleteApiPoisByIdMutation({
      client: apiClient,
    }),
    onError: (error) => {
      setPoiKey('');

      toast.error('Failed', {
        description: error.message,
        duration: 2000,
      });
    },
    onSuccess: () => {
      toast.success('Success', {
        description: 'The point of interest has been deleted.',
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
          <DialogTitle>Delete Point of Interest</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this point of interest? This action
            cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col w-full h-auto gap-3">
          <Input
            type="text"
            value={poiKey}
            onChange={(e) => setPoiKey(e.target.value)}
            placeholder="Enter key to confirm"
          />
          <Label className="text-sm text-muted-foreground">
            Please enter {poiKey} to confirm.
          </Label>
        </div>

        <DialogFooter>
          <Button
            disabled={poiKey !== key}
            onClick={() => deletePoi.mutate({ path: { id } })}
          >
            Continue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
