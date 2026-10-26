'use client'

import { useEffect, useState } from 'react'
import { AppSidebar, type AppSidebarProps } from './app-sidebar'
import { ChevronsRight } from './icons'
import BottomTabs from './BottomTabs'
import { SidebarInset, SidebarProvider, SidebarTrigger } from './ui/sidebar'
import { useAuthState } from '../hooks/useAuthState'
import { isStaffRole } from '../lib/auth'
import { getDueReviews, getSettings } from '../lib/api'
import { resolveAvatar } from '../lib/dicebear'

type AppShellProps = {
  children: React.ReactNode
} & Partial<AppSidebarProps>

// The product chrome. Public pages keep Header; the signed-in app uses the
// sidebar that /profile already proved. Every page that adopts this shell gets
// the same mobile trigger, the same rail and the same bottom tabs, so no branch
// can quietly ship without navigation.
export default function AppShell({ children, ...overrides }: AppShellProps) {
  const { loggedIn, user } = useAuthState()
  const [dueReviews, setDueReviews] = useState<number | undefined>(undefined)
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)

  useEffect(() => {
    if (!loggedIn || !user) {
      setAvatarUrl(undefined)
      setAvatarPreset(null)
      return
    }
    let live = true
    getSettings()
      .then((settings) => {
        if (!live) return
        const resolved = resolveAvatar(user, settings)
        setAvatarPreset(resolved.preset ?? null)
        setAvatarUrl(resolved.url)
      })
      .catch(() => {
        // The generated seed avatar is still a usable fallback.
      })
    return () => {
      live = false
    }
  }, [loggedIn, user])

  useEffect(() => {
    if (!loggedIn) {
      setDueReviews(0)
      return
    }
    let live = true
    getDueReviews()
      .then((data) => {
        if (live) setDueReviews(data.count)
      })
      .catch(() => {
        // The badge is an extra. A failed count must not take the shell down.
      })
    return () => {
      live = false
    }
  }, [loggedIn])

  return (
    <>
      <SidebarProvider>
        <AppSidebar
          name={overrides.name ?? user?.name ?? 'Guest'}
          studentId={overrides.studentId ?? user?.student_id ?? 'guest'}
          isAdmin={overrides.isAdmin ?? isStaffRole(user?.role)}
          username={overrides.username ?? user?.username}
          level={overrides.level ?? user?.level}
          streak={overrides.streak ?? user?.current_streak}
          avatarUrl={overrides.avatarUrl ?? avatarUrl}
          avatarPreset={overrides.avatarPreset !== undefined ? overrides.avatarPreset : avatarPreset}
          dueReviews={overrides.dueReviews ?? dueReviews}
        />
        <SidebarInset>
          <div className="mb-2 flex justify-start md:hidden px-4 sm:px-6 pt-8">
            <SidebarTrigger variant="ghost" className="flex-col gap-0.5">
              <span aria-hidden="true" className="font-mono text-sm text-mathua-blue">λ</span>
              <ChevronsRight size={10} className="text-mathua-muted" />
            </SidebarTrigger>
          </div>
          {children}
        </SidebarInset>
      </SidebarProvider>
      <BottomTabs />
    </>
  )
}
