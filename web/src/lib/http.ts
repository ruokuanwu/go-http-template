import type { ApiResponse } from '@/types/api'

export class ApiError extends Error {
  constructor(message: string, public readonly status: number, public readonly code?: number) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(path, options)
  const body = await response.json().catch(() => null) as ApiResponse<T> | null
  if (!response.ok || body?.code !== 0) {
    throw new ApiError(body?.message || `请求失败（${response.status}）`, response.status, body?.code)
  }
  return body.data as T
}
