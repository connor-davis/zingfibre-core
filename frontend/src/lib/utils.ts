import { QueryClient } from '@tanstack/react-query';

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

import {
  type ClientOptions,
  createClient,
  createConfig,
} from '@/api-client/client';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const queryClient = new QueryClient();

export const apiClient = createClient(
  createConfig<ClientOptions>({
    baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:6173',
    credentials: 'include',
  })
);
