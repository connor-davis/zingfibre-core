import { getApiPoisOptions } from '@/api-client/@tanstack/react-query.gen';
import { useQuery } from '@tanstack/react-query';
import { Link, useRouter, useSearch } from '@tanstack/react-router';

import type { PointOfInterest } from '@/api-client';
import { apiClient } from '@/lib/utils';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';
import { SidebarTrigger } from './ui/sidebar';

export default function Header() {
  const router = useRouter();
  const routerState = router.state;

  const { poi } = useSearch({ from: '__root__' });

  const { data: poiData, isLoading: isLoadingPointsOfInterest } = useQuery({
    ...getApiPoisOptions({
      client: apiClient,
    }),
  });

  return (
    <div className="flex items-center justify-between gap-3">
      <div className="flex items-center gap-3 p-3">
        <SidebarTrigger />

        <Link to="/">
          <img
            src="/zing-logo.png"
            alt="Zing Logo"
            className="h-7 dark:hidden"
          />

          <img
            src="/zing-logo-dark.png"
            alt="Zing Logo"
            className="h-7 hidden dark:block"
          />
        </Link>
      </div>

      {!isLoadingPointsOfInterest && (
        <div className="flex items-center gap-3 p-3">
          <Select
            value={poi ?? '--'}
            onValueChange={(value) => {
              router.navigate({
                to: routerState.location.pathname,
                search: (previous) => ({
                  ...previous,
                  poi: value === '--' ? undefined : value,
                }),
              });
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select Point of Interest" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="--">Global</SelectItem>
              {[...((poiData?.data ?? []) as Array<PointOfInterest>)].map(
                (poi) => (
                  <SelectItem key={poi.ID!} value={poi.Key!}>
                    {poi.Name!}
                  </SelectItem>
                )
              )}
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  );
}
