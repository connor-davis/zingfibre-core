import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/reports/recharges')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/reports/recharges"!</div>;
}
