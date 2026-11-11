import { useEffect } from 'react'
import { create } from 'zustand'
import { useMediaQuery } from '@/hooks/use-media-query'

export type Theme = 'light' | 'dark' | 'system'

type PreferencesState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

function readTheme(): Theme {
  const saved = localStorage.getItem('theme')
  return saved === 'light' || saved === 'dark' ? saved : 'system'
}

export const usePreferences = create<PreferencesState>((set) => ({
  theme: readTheme(),
  setTheme: (theme) => {
    localStorage.setItem('theme', theme)
    set({ theme })
  },
}))

export function useThemeSync() {
  const theme = usePreferences((state) => state.theme)
  const systemDark = useMediaQuery('(prefers-color-scheme: dark)')

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark' || (theme === 'system' && systemDark))
  }, [theme, systemDark])
}
