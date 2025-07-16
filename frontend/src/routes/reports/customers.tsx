import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/reports/customers')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/reports/customers"!</div>;
}
