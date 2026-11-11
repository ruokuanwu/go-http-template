import { createBrowserRouter } from 'react-router'
import AppLayout from '@/layouts/app-layout'
import DashboardPage from '@/pages/dashboard'
import TodosPage from '@/pages/todos'
import SettingsPage from '@/pages/settings'
import NotFoundPage from '@/pages/not-found'

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'todos', element: <TodosPage /> },
      { path: 'settings', element: <SettingsPage /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
])
