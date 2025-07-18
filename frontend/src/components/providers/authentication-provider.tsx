import { getApiAuthenticationCheckOptions } from '@/api-client/@tanstack/react-query.gen';
import { useQuery } from '@tanstack/react-query';
import { type ReactNode, createContext, useContext } from 'react';

import type { User } from '@/api-client';
import { apiClient } from '@/lib/utils';

const AuthenticationContext = createContext<{
  user?: User;
  isLoading: boolean;
  isError: boolean;
}>({
  user: undefined,
  isLoading: false,
  isError: false,
});

export const AuthenticationProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const { data, isLoading, isError } = useQuery({
    ...getApiAuthenticationCheckOptions({
      client: apiClient,
    }),
    retry: 0,
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
  });

  return (
    <AuthenticationContext.Provider
      value={{ user: data?.data as User, isLoading, isError }}
    >
      {children}
    </AuthenticationContext.Provider>
  );
};

export const useAuthentication = () => useContext(AuthenticationContext);
