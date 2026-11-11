import { Button } from '@/components/ui/button'

export function LoadingState() {
  return <p role="status" className="py-8 text-center text-sm text-muted-foreground">正在加载…</p>
}

export function EmptyState({ message }: { message: string }) {
  return <p className="py-8 text-center text-sm text-muted-foreground">{message}</p>
}

export function ErrorState({ error, onRetry, disabled = false }: {
  error: Error; onRetry: () => void; disabled?: boolean
}) {
  return (
    <div role="alert" className="flex flex-wrap items-center justify-between gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
      <span>{error.message}</span>
      <Button type="button" variant="outline" size="sm" disabled={disabled} onClick={onRetry}>重新加载</Button>
    </div>
  )
}
