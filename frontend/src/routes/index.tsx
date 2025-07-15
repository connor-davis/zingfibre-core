import { createFileRoute } from '@tanstack/react-router';

import AuthenticationGuard from '@/components/guards/authentication-guard';
import { Label } from '@/components/ui/label';

export const Route = createFileRoute('/')({
  component: () => (
    <AuthenticationGuard>
      <App />
    </AuthenticationGuard>
  ),
});

function App() {
  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Dashboard</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>
    </div>
  );
}
