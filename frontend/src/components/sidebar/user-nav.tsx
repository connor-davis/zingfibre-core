import {
  getApiAuthenticationCheckQueryKey,
  postApiAuthenticationLogoutMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  LogOutIcon,
  MoonIcon,
  MoreVerticalIcon,
  PaintbrushIcon,
  SunIcon,
} from 'lucide-react';

import { capitalCase, constantCase } from 'change-case';
import { toast } from 'sonner';

import { useAuthentication } from '@/components/providers/authentication-provider';
import { useTheme } from '@/components/providers/theme-provider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar';
import { apiClient } from '@/lib/utils';

import { Avatar, AvatarFallback } from '../ui/avatar';

export default function UserNav() {
  const queryClient = useQueryClient();
  const { isMobile } = useSidebar();
  const { setTheme, setAppearance } = useTheme();

  const { user, isLoading } = useAuthentication();

  if (isLoading) return null;
  if (!user) return null;

  const logout = useMutation({
    ...postApiAuthenticationLogoutMutation({
      client: apiClient,
    }),
    onError: (error) =>
      toast.error('Failed', {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      return toast.success('Success', {
        description: 'You have been logged out.',
        duration: 2000,
        onAutoClose: () =>
          queryClient.invalidateQueries({
            queryKey: getApiAuthenticationCheckQueryKey(),
          }),
      });
    },
  });

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar>
                <AvatarFallback>
                  {constantCase(
                    (user.Email ?? 'none')
                      .split('@')[0]
                      .split('.')
                      .slice(0, 2)
                      .map((word) => word.charAt(0))
                      .join('')
                  )}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user.Email}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {capitalCase(
                    Array.isArray(user.Role)
                      ? user.Role.join(', ')
                      : (user.Role ?? '')
                  )}
                </span>
              </div>
              <MoreVerticalIcon className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? 'bottom' : 'right'}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar>
                  <AvatarFallback>
                    {constantCase(
                      (user.Email ?? 'none')
                        .split('@')[0]
                        .split('.')
                        .slice(0, 2)
                        .map((word) => word.charAt(0))
                        .join('')
                    )}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user.Email}</span>
                  <span className="truncate text-xs text-muted-foreground">
                    {capitalCase(
                      Array.isArray(user.Role)
                        ? user.Role.join(', ')
                        : (user.Role ?? '')
                    )}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              {/* <DropdownMenuItem asChild>
                <Link to="/account">
                  <UserCircleIcon />
                  Account
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator /> */}
              <DropdownMenuSub>
                <DropdownMenuSubTrigger className="p-0">
                  <DropdownMenuItem>
                    <PaintbrushIcon />
                    Appearance
                  </DropdownMenuItem>
                </DropdownMenuSubTrigger>
                <DropdownMenuPortal>
                  <DropdownMenuSubContent>
                    <DropdownMenuItem onClick={() => setAppearance('zing')}>
                      Zing
                    </DropdownMenuItem>
                  </DropdownMenuSubContent>
                </DropdownMenuPortal>
              </DropdownMenuSub>
              <DropdownMenuSub>
                <DropdownMenuSubTrigger className="p-0">
                  <DropdownMenuItem>
                    <SunIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                    <MoonIcon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                    Theme
                  </DropdownMenuItem>
                </DropdownMenuSubTrigger>
                <DropdownMenuPortal>
                  <DropdownMenuSubContent>
                    <DropdownMenuItem onClick={() => setTheme('light')}>
                      Light
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => setTheme('dark')}>
                      Dark
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => setTheme('system')}>
                      System
                    </DropdownMenuItem>
                  </DropdownMenuSubContent>
                </DropdownMenuPortal>
              </DropdownMenuSub>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => logout.mutate({})}>
              <LogOutIcon />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
