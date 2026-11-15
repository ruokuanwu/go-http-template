import { queryOptions, useQuery, type QueryClient } from "@tanstack/react-query";
import { z } from "zod";
import { apiClient } from "@/lib/api-client";

const sessionSchema = z.object({
  user: z.object({ id: z.number().int().positive(), name: z.string(), email: z.string() }),
});
export type AuthSession = z.infer<typeof sessionSchema>;
type Credentials = { email: string; password: string };
const sessionKey = ["auth", "session"] as const;

export const authClient = {
  register: (input: Credentials & { name: string }) => apiClient("/api/auth/register", sessionSchema.parse, {
    method: "POST", body: JSON.stringify(input),
  }),
  login: (input: Credentials) => apiClient("/api/auth/login", sessionSchema.parse, {
    method: "POST", body: JSON.stringify(input),
  }),
  getSession: (signal?: AbortSignal) => apiClient("/api/auth/session", (data) => sessionSchema.nullable().parse(data), { signal }),
  logout: () => apiClient("/api/auth/logout", (data) => z.object({}).parse(data), { method: "POST" }),
};

export const sessionQueryOptions = queryOptions({
  queryKey: sessionKey,
  queryFn: ({ signal }) => authClient.getSession(signal),
  staleTime: 0,
  retry: false,
});

export function useSession() { return useQuery(sessionQueryOptions); }

export async function updateSession(client: QueryClient, session: AuthSession | null) {
  await client.cancelQueries();
  // A changed identity must never inherit the previous user's private cache.
  client.removeQueries({ predicate: (query) => query.queryKey[0] !== "auth" && query.queryKey[0] !== "demo" });
  client.setQueryData(sessionKey, session);
}
