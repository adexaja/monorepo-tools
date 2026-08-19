import { Outlet, createRootRoute } from "@tanstack/react-router";
import "../styles/globals.css";

export const Route = createRootRoute({
  component: () => <Outlet />,
});
