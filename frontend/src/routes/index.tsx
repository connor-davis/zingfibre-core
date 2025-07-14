import { createFileRoute } from '@tanstack/react-router';

import AuthenticationGuard from '@/components/guards/authentication-guard';

export const Route = createFileRoute('/')({
  component: () => (
    <AuthenticationGuard>
      <App />
    </AuthenticationGuard>
  ),
});

function App() {
  return <div></div>;
}
