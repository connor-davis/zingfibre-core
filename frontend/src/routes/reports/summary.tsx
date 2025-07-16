import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/reports/summary')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/reports/summary"!</div>;
}
