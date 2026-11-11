import { Link } from 'react-router'
import { useTodos } from '@/features/todos/hooks/use-todos'
import { ErrorState, LoadingState } from '@/components/common/async-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function DashboardPage() {
  const query = useTodos()
  const todos = query.data ?? []
  const completed = todos.filter(todo => todo.completed).length
  return (
    <section className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div><h1 className="text-2xl font-semibold">概览</h1><p className="mt-2 text-sm text-muted-foreground">全部任务的当前统计。</p></div>
        <Button asChild><Link to="/todos">管理任务</Link></Button>
      </div>
      {query.error && <ErrorState error={query.error} onRetry={() => void query.refetch()} disabled={query.isFetching} />}
      {query.isPending ? <LoadingState /> : query.data && (
        <div className="grid gap-4 sm:grid-cols-3">
          {[['全部任务', todos.length], ['待完成', todos.length - completed], ['已完成', completed]].map(([label, count]) => (
            <Card key={label}><CardHeader><CardTitle>{label}</CardTitle></CardHeader><CardContent className="text-3xl font-semibold">{count}</CardContent></Card>
          ))}
        </div>
      )}
    </section>
  )
}
