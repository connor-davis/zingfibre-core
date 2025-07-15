import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarRail,
  SidebarSeparator,
} from '@/components/ui/sidebar';

import UserNav from './user-nav';

export default function AppSidebar() {
  return (
    <Sidebar>
      <SidebarContent className="gap-0">
        <SidebarSeparator className="mx-0" />
      </SidebarContent>
      <SidebarFooter className="px-0">
        <SidebarGroup className="py-0">
          <UserNav />
        </SidebarGroup>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
