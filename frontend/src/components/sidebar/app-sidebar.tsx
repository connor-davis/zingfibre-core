import { Link } from '@tanstack/react-router';
import { LayoutDashboardIcon, MapPinIcon, UsersIcon } from 'lucide-react';

import UserNav from '@/components/sidebar/user-nav';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
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
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="px-0">
        <RoleGuard value={['admin', 'staff']}>
          <SidebarSeparator className="mx-0" />

          <SidebarGroup className="py-0">
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/pois">
                    <MapPinIcon />
                    <span>Points of Interest</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
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
