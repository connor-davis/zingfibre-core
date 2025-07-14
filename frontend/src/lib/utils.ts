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
    baseUrl: import.meta.env.VITE_API_URL || 'http://localhost:4000',
    credentials: 'include',
  })
);
