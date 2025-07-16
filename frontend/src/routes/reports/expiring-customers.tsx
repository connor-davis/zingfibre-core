import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/reports/expiring-customers')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/reports/expiring-customers"!</div>;
}
