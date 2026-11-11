import { useEffect, useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { request } from '@/lib/api'

type Todo = { id: number; title: string; completed: boolean }

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [title, setTitle] = useState('')
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  async function load() {
    setLoading(true)
    setError('')
    try {
      const result = await request<{ items: Todo[] }>('/api/todos')
      setTodos(result.items)
    } catch (error) {
      setError(error instanceof Error ? error.message : '请求失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [])

  async function mutate(action: () => Promise<void>) {
    setBusy(true)
    setError('')
    try {
      await action()
    } catch (error) {
      setError(error instanceof Error ? error.message : '请求失败，请稍后重试')
    } finally {
      setBusy(false)
    }
  }

  function add(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!title.trim() || busy) return
    void mutate(async () => {
      const result = await request<Todo>('/api/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: title.trim() }),
      })
      setTodos(previous => [...previous, result])
      setTitle('')
    })
  }

  function toggle(todo: Todo) {
    void mutate(async () => {
      const result = await request<Todo>(`/api/todos/${todo.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ completed: !todo.completed }),
      })
      setTodos(previous => previous.map(item => item.id === todo.id ? result : item))
    })
  }

  function remove(id: number) {
    void mutate(async () => {
      await request<Record<string, never>>(`/api/todos/${id}`, { method: 'DELETE' })
      setTodos(previous => previous.filter(item => item.id !== id))
    })
  }

  return (
    <main className="mx-auto max-w-xl px-4 py-12 sm:py-20">
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
            <div role="alert" className="flex flex-wrap items-center justify-between gap-2 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              <span>{error}</span>
              <Button type="button" variant="outline" size="sm" disabled={busy || loading} onClick={() => void load()}>重新加载</Button>
            </div>
          )}
          <div aria-live="polite" aria-busy={loading || busy}>
            {loading ? <p className="py-6 text-center text-sm text-muted-foreground">正在加载任务…</p> : (
              todos.length === 0 ? <p className="py-6 text-center text-sm text-muted-foreground">暂无任务，添加第一项吧。</p> : (
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
    </main>
  )
}
