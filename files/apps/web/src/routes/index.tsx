import { createRoute } from "@tanstack/react-router";
import { Button } from "__PROJECT_SLUG__-ui";
import { Route as rootRoute } from "./__root";

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: HomePage,
});

function HomePage() {
  return (
    <main className="shell">
      <p className="eyebrow">__PROJECT_NAME__</p>
      <h1>Ship a focused product with a boring, durable foundation.</h1>
      <p className="lede">
        TanStack Start handles the full-stack React app, while the Go API,
        PostgreSQL, and Shoebox worker keep server-side boundaries explicit.
      </p>
      <Button>Get started</Button>
    </main>
  );
}
