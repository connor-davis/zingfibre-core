import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRoute } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import { AuthenticationProvider } from '@/components/providers/authentication-provider';
import { ThemeProvider } from '@/components/providers/theme-provider';
import { Toaster } from '@/components/ui/sonner';
import { queryClient } from '@/lib/utils';

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider defaultTheme="system" defaultAppearance="zing">
      <QueryClientProvider client={queryClient}>
        <AuthenticationProvider>
          <div className="flex flex-col w-screen h-screen bg-muted">
            <Outlet />
            {import.meta.env.MODE === 'development' && (
              <TanStackRouterDevtools />
            )}
          </div>
        </AuthenticationProvider>
      </QueryClientProvider>

      <Toaster />
    </ThemeProvider>
  ),
});
