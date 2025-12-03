import { postApiDynamicQueriesMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import type { ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type CreateDynamicQuery, type ErrorResponse } from '@/api-client';
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
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/utils';

export default function CreateDynamicQueryDialog({
  children,
}: {
  children: ReactNode;
}) {
  const router = useRouter();

  const createDynamicQueryForm = useForm<CreateDynamicQuery>({
    defaultValues: {
      Name: undefined,
      Prompt: undefined,
    },
  });

  const createDynamicQuery = useMutation({
    ...postApiDynamicQueriesMutation({
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
        description: 'The dynamic query has been created successfully.',
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
          <DialogTitle>Create Dynamic Report</DialogTitle>
          <DialogDescription>
            Use this form to create a new dynamic report. Please fill in all
            required fields.
          </DialogDescription>
        </DialogHeader>

        <Form {...createDynamicQueryForm}>
          <form
            onSubmit={createDynamicQueryForm.handleSubmit((values) =>
              createDynamicQuery.mutate({
                body: values,
              })
            )}
            className="flex flex-col gap-3"
          >
            <FormField
              control={createDynamicQueryForm.control}
              name="Name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Name" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter the dynamic reports name.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createDynamicQueryForm.control}
              name="Prompt"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Prompt</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="What would you like the dynamic report to show."
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Please enter the dynamic report prompt.
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
