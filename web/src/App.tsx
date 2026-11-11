import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router/dom'
import { PreferencesProvider } from '@/stores/preferences'
import { router } from '@/router'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 30_000, retry: false } },
})

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <PreferencesProvider><RouterProvider router={router} /></PreferencesProvider>
    </QueryClientProvider>
  )
}
