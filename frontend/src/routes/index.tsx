import { createFileRoute, useRouter, useSearch } from '@tanstack/react-router';
import { CalendarIcon, SigmaIcon } from 'lucide-react';

import { CartesianGrid, Line, LineChart, XAxis, YAxis } from 'recharts';
import uniqolor from 'uniqolor';
import z from 'zod';

import {
  type RechargeTypeCounts,
  getApiAnalyticsRechargeTypeCounts,
} from '@/api-client';
import AuthenticationGuard from '@/components/guards/authentication-guard';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { DebounceIncrementorInput } from '@/components/ui/debounce-incrementor';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/')({
  component: () => (
    <AuthenticationGuard>
      <App />
    </AuthenticationGuard>
  ),
  validateSearch: z.object({
    count: z.number().optional(),
    period: z.string().optional(),
    poi: z.string().optional(),
  }),
  loaderDeps: ({ search: { count, period, poi } }) => ({
    count: count ?? 1,
    period: period ?? 'months',
    poi: poi ?? '',
  }),
  loader: async ({ deps: { count, period, poi } }) => {
    const { data } = await getApiAnalyticsRechargeTypeCounts({
      client: apiClient,
      query: {
        count,
        period,
        poi,
      },
    });

    const result = {
      data: data?.data,
      pages: data ? (data.pages ? data.pages : 1) : 1,
    } as {
      data: RechargeTypeCounts;
      pages: number;
    };

    return {
      items: result.data.Items,
      types: result.data.Types,
      pages: result.pages,
    };
  },
});

function App() {
  const { count, period } = useSearch({ from: '/' });
  const router = useRouter();

  const { items, types } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Dashboard</Label>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-3">
            <SigmaIcon className="size-4" />
            <DebounceIncrementorInput
              className="w-24"
              min={1}
              defaultValue={count}
              onChange={(value) => {
                console.log(value);

                router.navigate({
                  to: '/',
                  search: (previous) => ({
                    ...previous,
                    count: value.target.valueAsNumber,
                  }),
                });
              }}
            />
          </div>

          <div className="flex items-center gap-3">
            <CalendarIcon className="size-4" />
            <Select
              defaultValue="months"
              value={period}
              onValueChange={(value) => {
                router.navigate({
                  to: '/',
                  search: (previous) => ({
                    ...previous,
                    period: value,
                  }),
                });
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select Period" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="weeks">Weeks</SelectItem>
                <SelectItem value="months">Months</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      <Card className="pt-0 w-full h-full bg-background">
        <CardHeader className="flex items-center gap-2 space-y-0 py-5 sm:flex-row">
          <div className="grid flex-1 gap-1">
            <CardTitle>Recharge Counts</CardTitle>
            <CardDescription>
              This chart shows the number of recharges made over the selected
              period, grouped by recharge type.
            </CardDescription>
          </div>
          <div className="flex items-center gap-1"></div>
        </CardHeader>
        <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6 w-full h-full">
          {items && types && (
            <ChartContainer config={{}} className="aspect-auto w-full h-full">
              <LineChart data={items ?? []}>
                <CartesianGrid strokeDasharray="3 3" />

                <XAxis dataKey="Period" />
                <YAxis />

                <ChartTooltip
                  wrapperStyle={{ pointerEvents: 'auto' }}
                  cursor={true}
                  content={
                    <ChartTooltipContent
                      indicator="line"
                      className="bg-background"
                    />
                  }
                />
                <ChartLegend content={<ChartLegendContent />} />

                {[...(types ?? [])].map((type) => (
                  <Line
                    key={type}
                    dataKey={type}
                    stroke={uniqolor(type).color}
                    connectNulls
                    type="monotone"
                  />
                ))}
              </LineChart>
            </ChartContainer>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
