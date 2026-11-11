export type ApiResponse<T> = {
  data: T
  code: 0
  message: string
} | {
  code: number
  message: string
  data?: never
}

