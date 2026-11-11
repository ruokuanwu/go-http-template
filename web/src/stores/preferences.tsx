import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { useMediaQuery } from '@/hooks/use-media-query'

export type Theme = 'light' | 'dark' | 'system'
const PreferencesContext = createContext<{ theme: Theme; setTheme: (theme: Theme) => void } | null>(null)

function readTheme(): Theme {
  const saved = localStorage.getItem('theme')
  return saved === 'light' || saved === 'dark' ? saved : 'system'
}

export function PreferencesProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(readTheme)
  const systemDark = useMediaQuery('(prefers-color-scheme: dark)')

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark' || (theme === 'system' && systemDark))
    localStorage.setItem('theme', theme)
  }, [theme, systemDark])

  return <PreferencesContext value={{ theme, setTheme }}>{children}</PreferencesContext>
}

export function usePreferences() {
  const preferences = useContext(PreferencesContext)
  if (!preferences) throw new Error('usePreferences 必须在 PreferencesProvider 内使用')
  return preferences
}
