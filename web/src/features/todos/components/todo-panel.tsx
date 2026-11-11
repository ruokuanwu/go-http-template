import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { EmptyState, ErrorState, LoadingState } from '@/components/common/async-state'
import { cn } from '@/lib/utils'
import { useTodos, useTodoMutation } from '../hooks/use-todos'
import type { Todo } from '../types'

export default function TodoPanel() {
  const query = useTodos()
  const mutation = useTodoMutation()
  const [title, setTitle] = useState('')
  const todos = query.data ?? []
  const loading = query.isPending
  const busy = mutation.isPending
  const error = mutation.error ?? query.error

  function add(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!title.trim() || busy) return
    mutation.mutate({ type: 'create', title: title.trim() }, { onSuccess: () => setTitle('') })
  }

  function toggle(todo: Todo) {
    mutation.mutate({ type: 'update', id: todo.id, completed: !todo.completed })
  }

  function remove(id: number) {
    mutation.mutate({ type: 'delete', id })
  }

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle><h1 className="text-2xl">任务清单</h1></CardTitle>
          <CardDescription>记录待办，专注眼前的小事。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-5">
          <form onSubmit={add} className="flex gap-2">
            <Input aria-label="新任务" placeholder="添加一项任务…" value={title} onChange={event => setTitle(event.target.value)} disabled={busy || loading} />
            <Button type="submit" disabled={busy || loading || !title.trim()}>添加</Button>
          </form>
          {error && (
            <ErrorState error={error} disabled={busy || query.isFetching} onRetry={() => { mutation.reset(); void query.refetch() }} />
          )}
          <div aria-live="polite" aria-busy={loading || busy}>
            {loading ? <LoadingState /> : (
              query.isError && !query.data ? null : todos.length === 0 ? <EmptyState message="暂无任务，添加第一项吧。" /> : (
                <ul className="divide-y">
                  {todos.map(todo => (
                    <li key={todo.id} className="flex items-center gap-3 py-3">
                      <label className="flex min-w-0 flex-1 items-center gap-3">
                        <input type="checkbox" checked={todo.completed} disabled={busy} onChange={() => toggle(todo)} className="size-4 shrink-0 accent-primary" />
                        <span className={cn('break-all text-sm', todo.completed && 'text-muted-foreground line-through')}>{todo.title}</span>
                      </label>
                      <Button type="button" variant="ghost" size="sm" disabled={busy} aria-label={`删除任务：${todo.title}`} onClick={() => remove(todo.id)}>删除</Button>
                    </li>
                  ))}
                </ul>
              )
            )}
          </div>
        </CardContent>
      </Card>
    </>
  )
}
