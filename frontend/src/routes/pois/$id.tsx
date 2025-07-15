import { putApiPoisByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useParams } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';

import { type PointOfInterest, getApiPoisById } from '@/api-client';
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

export const Route = createFileRoute('/pois/$id')({
  component: () => <RouteComponent />,
  loader: async ({ params }) => {
    const { id } = params;

    const { data } = await getApiPoisById({
      client: apiClient,
      path: { id },
    });

    return {
      poi: data?.data,
    } as {
      poi?: PointOfInterest;
    };
  },
});

function RouteComponent() {
  const { id } = useParams({ from: '/pois/$id' });
  const { poi } = Route.useLoaderData();

  const poiForm = useForm<PointOfInterest>({
    values: poi,
  });

  const updatePoi = useMutation({
    ...putApiPoisByIdMutation({
      client: apiClient,
      path: { id },
    }),
    onError: (error) => {
      toast.error('Failed', {
        description: error.message,
        duration: 2000,
      });
    },
    onSuccess: ({ data: updatedPoi }) => {
      poiForm.reset({ ...poi, ...updatedPoi });

      toast.success('Success', {
        description: 'The point of interest has been updated.',
        duration: 2000,
      });
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/pois">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Edit Point of Interest</Label>
        </div>
      </div>

      <Form {...poiForm}>
        <form
          onSubmit={poiForm.handleSubmit((values) =>
            updatePoi.mutate({
              path: { id },
              body: values,
            })
          )}
          className="space-y-6 w-full"
        >
          <FormField
            control={poiForm.control}
            name="Key"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Key</FormLabel>
                <FormControl>
                  <Input placeholder="Key" {...field} />
                </FormControl>
                <FormDescription>
                  Please enter the point of interest's key.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={poiForm.control}
            name="Name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Name" {...field} />
                </FormControl>
                <FormDescription>
                  Please enter the point of interest's name.
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
