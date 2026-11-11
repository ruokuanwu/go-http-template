import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createTodo, deleteTodo, listTodos, updateTodo } from '../api'

const todoKey = ['todos'] as const

export function useTodos() {
  return useQuery({ queryKey: todoKey, queryFn: listTodos })
}

export function useTodoMutation() {
  const client = useQueryClient()
  return useMutation({
    mutationFn: async (action:
      | { type: 'create'; title: string }
      | { type: 'update'; id: number; completed: boolean }
      | { type: 'delete'; id: number }) => {
      switch (action.type) {
        case 'create': return createTodo(action.title)
        case 'update': return updateTodo(action.id, action.completed)
        case 'delete': return deleteTodo(action.id)
      }
    },
    onSuccess: () => client.invalidateQueries({ queryKey: todoKey }),
  })
}
