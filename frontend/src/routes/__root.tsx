import { QueryClientProvider } from '@tanstack/react-query';
import { Link, Outlet, createRootRoute } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import AuthenticationGuard from '@/components/guards/authentication-guard';
import { AuthenticationProvider } from '@/components/providers/authentication-provider';
import { ThemeProvider } from '@/components/providers/theme-provider';
import AppSidebar from '@/components/sidebar/app-sidebar';
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
import { Toaster } from '@/components/ui/sonner';
import { queryClient } from '@/lib/utils';

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider defaultTheme="system" defaultAppearance="zing">
      <QueryClientProvider client={queryClient}>
        <AuthenticationProvider>
          <div className="flex flex-col w-screen h-screen bg-background">
            <SidebarProvider>
              <AuthenticationGuard>
                <div className="flex w-screen h-screen overflow-hidden">
                  <AppSidebar />

                  <div className="flex flex-col w-full h-full overflow-hidden">
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

                    <Outlet />
                    {import.meta.env.MODE === 'development' && (
                      <TanStackRouterDevtools position="top-right" />
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
  ),
});
