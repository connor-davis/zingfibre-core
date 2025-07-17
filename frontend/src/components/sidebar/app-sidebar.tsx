import { Link } from '@tanstack/react-router';
import {
  CircleDotIcon,
  LayoutDashboardIcon,
  RefreshCcwDotIcon,
  RefreshCcwIcon,
  SpeechIcon,
  UsersIcon,
} from 'lucide-react';

import UserNav from '@/components/sidebar/user-nav';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarSeparator,
} from '@/components/ui/sidebar';

import RoleGuard from '../guards/role-guard';

export default function AppSidebar() {
  return (
    <Sidebar collapsible="icon">
      <SidebarContent className="gap-0">
        <SidebarGroup>
          <SidebarGroupLabel>Analytics</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/">
                    <LayoutDashboardIcon />
                    <span>Dashboard</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="mx-0" />

        <SidebarGroup>
          <SidebarGroupLabel>Reports</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/reports/customers">
                    <UsersIcon />
                    <span>Customers</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/reports/expiring-customers">
                    <SpeechIcon />
                    <span>Expiring Customers</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/reports/recharges">
                    <RefreshCcwIcon />
                    <span>Recharges</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/reports/recharges-summary">
                    <RefreshCcwDotIcon />
                    <span>Recharges Summary</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/reports/summary">
                    <CircleDotIcon />
                    <span>Summary</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="mx-0" />
      </SidebarContent>

      <SidebarFooter className="px-0">
        <RoleGuard value={['admin', 'staff']}>
          <SidebarSeparator className="mx-0" />

          <SidebarGroup className="py-0">
            <SidebarGroupLabel>System</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <RoleGuard value={['admin']}>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/users">
                        <UsersIcon />
                        <span>Users</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </RoleGuard>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </RoleGuard>

        <SidebarSeparator className="mx-0" />

        <SidebarGroup className="py-0">
          <UserNav />
        </SidebarGroup>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
