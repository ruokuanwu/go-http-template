import { config } from "@/config";

export class ApiError extends Error {
  constructor(message: string, public readonly status: number, public readonly code: number) {
    super(message);
    this.name = "ApiError";
  }
}

export async function apiClient<T>(path: string, parse: (data: unknown) => T, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body) headers.set("Content-Type", "application/json");
  const response = await fetch(`${config.baseUrl}${path}`, { ...init, headers, credentials: "include" });
  const body: unknown = await response.json();
  if (typeof body !== "object" || body === null || !("code" in body) || typeof body.code !== "number" ||
    !("message" in body) || typeof body.message !== "string") {
    throw new Error("Invalid API response.");
  }
  if (!response.ok || body.code !== 0) {
    throw new ApiError(body.message || "Request failed.", response.status, body.code);
  }
  if (!("data" in body)) throw new Error("Invalid API response.");
  return parse(body.data);
}
