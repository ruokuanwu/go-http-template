import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { usePreferences, type Theme } from '@/stores/preferences'

export default function SettingsPage() {
  const { theme, setTheme } = usePreferences()
  return (
    <section className="mx-auto max-w-xl space-y-6">
      <h1 className="text-2xl font-semibold">设置</h1>
      <Card>
        <CardHeader><CardTitle>外观</CardTitle><CardDescription>主题偏好保存在当前浏览器中。</CardDescription></CardHeader>
        <CardContent>
          <fieldset className="flex flex-wrap gap-4">
            <legend className="mb-3 text-sm font-medium">主题模式</legend>
            {([['light', '浅色'], ['dark', '深色'], ['system', '跟随系统']] satisfies [Theme, string][]).map(([value, label]) => (
              <label key={value} className="flex items-center gap-2 text-sm"><input type="radio" name="theme" value={value} checked={theme === value} onChange={() => setTheme(value)} className="accent-primary" />{label}</label>
            ))}
          </fieldset>
        </CardContent>
      </Card>
    </section>
  )
}
