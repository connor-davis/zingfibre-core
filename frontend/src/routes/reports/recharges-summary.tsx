import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/reports/recharges-summary')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/reports/recharges-summary"!</div>;
}
