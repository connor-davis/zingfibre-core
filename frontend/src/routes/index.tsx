import { createFileRoute } from '@tanstack/react-router';

import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
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
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { Label } from '@/components/ui/label';
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

const chartConfig = {
  'one-month': {
    label: '1 Month',
    color: 'var(--chart-1)',
  },
  'one-week': {
    label: '1 Week',
    color: 'var(--chart-2)',
  },
  'one-day': {
    label: '1 Day',
    color: 'var(--chart-3)',
  },
  intro: {
    label: 'Intro Package',
    color: 'var(--chart-4)',
  },
} satisfies ChartConfig;

function App() {
  const { items, types } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Dashboard</Label>
        </div>
        <div className="flex items-center gap-3"></div>
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
          <ChartContainer
            config={chartConfig}
            className="aspect-auto w-full h-full"
          >
            <AreaChart data={items ?? []}>
              <defs>
                <linearGradient id="fillDesktop" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="var(--color-one-month)"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="95%"
                    stopColor="var(--color-one-month)"
                    stopOpacity={0.1}
                  />
                </linearGradient>
                <linearGradient id="fillMobile" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="var(--color-one-week)"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="95%"
                    stopColor="var(--color-one-week)"
                    stopOpacity={0.1}
                  />
                </linearGradient>
                <linearGradient id="fillTablet" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="var(--color-one-day)"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="95%"
                    stopColor="var(--color-one-day)"
                    stopOpacity={0.1}
                  />
                </linearGradient>
                <linearGradient id="fillIntro" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="var(--color-intro)"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="95%"
                    stopColor="var(--color-intro)"
                    stopOpacity={0.1}
                  />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="Period" />
              <YAxis />
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent indicator="dot" />}
              />

              {[...(types ?? [])].map((type) => (
                <Area key={type} dataKey={type} />
              ))}
            </AreaChart>
          </ChartContainer>
        </CardContent>
      </Card>
    </div>
  );
}
