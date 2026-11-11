import { Link } from 'react-router'
import { Button } from '@/components/ui/button'

export default function NotFoundPage() {
  return <section className="space-y-4 py-12 text-center"><h1 className="text-2xl font-semibold">404 · 页面不存在</h1><p className="text-muted-foreground">请检查地址，或返回概览。</p><Button asChild><Link to="/">返回概览</Link></Button></section>
}
