import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRoute } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import z from 'zod';

import AuthenticationGuard from '@/components/guards/authentication-guard';
import Header from '@/components/header';
import { AuthenticationProvider } from '@/components/providers/authentication-provider';
import { ThemeProvider } from '@/components/providers/theme-provider';
import AppSidebar from '@/components/sidebar/app-sidebar';
import { SidebarProvider } from '@/components/ui/sidebar';
import { Toaster } from '@/components/ui/sonner';
import { queryClient } from '@/lib/utils';

export const Route = createRootRoute({
  component: RootComponent,
  validateSearch: z.object({
    poi: z.string().optional(),
  }),
});

function RootComponent() {
  return (
    <ThemeProvider defaultTheme="system" defaultAppearance="zing">
      <QueryClientProvider client={queryClient}>
        <AuthenticationProvider>
          <div className="flex flex-col w-screen h-screen bg-background">
            <SidebarProvider>
              <AuthenticationGuard>
                <div className="flex w-screen h-screen overflow-hidden">
                  <AppSidebar />

                  <div className="flex flex-col w-full h-full overflow-hidden">
                    <Header />
                    <Outlet />
                    {import.meta.env.MODE === 'development' && (
                      <TanStackRouterDevtools position="bottom-right" />
                    )}
                  </div>
                </div>
              </AuthenticationGuard>
            </SidebarProvider>
          </div>
        </AuthenticationProvider>
      </QueryClientProvider>

      <Toaster />
    </ThemeProvider>
  );
}
