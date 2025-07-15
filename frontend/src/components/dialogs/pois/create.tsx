import { postApiPoisMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import type { ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type CreatePointOfInterest, type ErrorResponse } from '@/api-client';
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
import { apiClient } from '@/lib/utils';

export default function CreatePoiDialog({ children }: { children: ReactNode }) {
  const router = useRouter();

  const createPoiForm = useForm<CreatePointOfInterest>({
    defaultValues: {
      Name: undefined,
      Key: undefined,
    },
  });

  const createPoi = useMutation({
    ...postApiPoisMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) => {
      toast.error(error.error, {
        description: error.details,
        duration: 2000,
      });
    },
    onSuccess: () => {
      toast.success('Success', {
        description: 'The point of interest has been created successfully.',
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
          <DialogTitle>Create Point of Interest</DialogTitle>
          <DialogDescription>
            Use this form to create a new point of interest. Please fill in all
            required fields.
          </DialogDescription>
        </DialogHeader>

        <Form {...createPoiForm}>
          <form
            onSubmit={createPoiForm.handleSubmit((values) =>
              createPoi.mutate({
                body: values,
              })
            )}
          >
            <FormField
              control={createPoiForm.control}
              name="Key"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Key</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Key" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter the key for the point of interest.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createPoiForm.control}
              name="Name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Name" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter the name for the point of interest.
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
