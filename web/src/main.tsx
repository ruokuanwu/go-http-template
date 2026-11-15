import React from "react";
import ReactDOM from "react-dom/client";
import { QueryCache, MutationCache, QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import "./index.css";
import { ApiError } from "@/lib/api-client";
import { updateSession } from "@/features/auth/auth-client";

function handleAuthError(error: Error) {
  if (error instanceof ApiError && error.status === 401 && error.code === 30001) {
    void updateSession(queryClient, null).then(async () => {
      await router.navigate({ to: "/login", replace: true });
      await router.invalidate();
    });
  }
}

const queryClient = new QueryClient({
  queryCache: new QueryCache({ onError: handleAuthError }),
  mutationCache: new MutationCache({ onError: handleAuthError }),
  defaultOptions: { queries: { staleTime: 60_000, retry: (count, error) => !(error instanceof ApiError && error.status === 401) && count < 1 } },
});
const router = createRouter({ routeTree, context: { queryClient }, scrollRestoration: true });

queryClient.getQueryCache().subscribe((event) => {
  if (event.type === "updated" && event.action.type === "success" &&
    event.query.queryKey[0] === "auth" && event.query.queryKey[1] === "session" &&
    event.query.state.data === null && router.state.location.pathname.startsWith("/dashboard")) {
    void router.navigate({ to: "/login", replace: true });
  }
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>,
);
