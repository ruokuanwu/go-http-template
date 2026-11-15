export const config = {
  baseUrl: (import.meta.env.VITE_API_BASE_URL || window.location.origin).replace(/\/$/, ""),
  debug: import.meta.env.DEBUG,
} as const;
