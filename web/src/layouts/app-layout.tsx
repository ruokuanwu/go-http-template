import { NavLink, Outlet } from 'react-router'
import { cn } from '@/lib/utils'

const links = [
  { to: '/', label: '概览', end: true },
  { to: '/todos', label: '任务管理', end: false },
  { to: '/settings', label: '设置', end: false },
]

export default function AppLayout() {
  return (
    <div className="min-h-screen">
      <header className="border-b bg-card">
        <div className="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-4 px-4 py-4">
          <span className="font-semibold">Go HTTP · Web 模板</span>
          <nav aria-label="主导航" className="flex gap-1">
            {links.map(link => (
              <NavLink key={link.to} to={link.to} end={link.end} className={({ isActive }) => cn('rounded-md px-3 py-2 text-sm transition-colors hover:bg-accent', isActive ? 'bg-secondary font-medium text-secondary-foreground' : 'text-muted-foreground')}>
                {link.label}
              </NavLink>
            ))}
          </nav>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-8 sm:py-12"><Outlet /></main>
    </div>
  )
}
