import { createFileRoute, redirect } from "@tanstack/react-router";
import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { sessionQueryOptions } from "@/features/auth/auth-client";

export const Route = createFileRoute("/dashboard")({
  beforeLoad: async ({ context }) => {
    const session = await context.queryClient.fetchQuery(sessionQueryOptions);
    if (!session) throw redirect({ to: "/login" });
    return { authenticatedSession: session };
  },
  errorComponent: ({ error, reset }) => <main className="page-container py-24"><h1 className="text-2xl">Unable to load your account.</h1><p role="alert" className="my-5 text-sm text-muted-foreground">{error instanceof Error ? error.message : "Check your connection and try again."}</p><button className="button-primary" onClick={reset}>Try again</button></main>,
  component: DashboardLayout,
});
