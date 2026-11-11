import { request } from '@/lib/http'
import type { Todo } from './types'

export async function listTodos() {
  return (await request<{ items: Todo[] }>('/api/todos')).items
}

export function createTodo(title: string) {
  return request<Todo>('/api/todos', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
  })
}

export function updateTodo(id: number, completed: boolean) {
  return request<Todo>(`/api/todos/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ completed }),
  })
}

export function deleteTodo(id: number) {
  return request<Record<string, never>>(`/api/todos/${id}`, { method: 'DELETE' })
}
