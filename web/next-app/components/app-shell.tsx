'use client'

import { useEffect, useState } from 'react'
import { AppSidebar, type AppSidebarProps } from './app-sidebar'
import { ChevronsRight } from './icons'
import BottomTabs from './BottomTabs'
import { SidebarInset, SidebarProvider, SidebarTrigger } from './ui/sidebar'
import { useAuthState } from '../hooks/useAuthState'
import { isStaffRole } from '../lib/auth'
import { getDueReviews, getSettings, listMyReports } from '../lib/api'
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
  // Leaderboard navigation follows the same setting the server reads to anonymise
  // the rows. Default true so a slow or failed settings call never removes a
  // destination the learner has not actually turned off.
  const [showLeaderboards, setShowLeaderboards] = useState(true)
  const [notificationUnread, setNotificationUnread] = useState(0)

  useEffect(() => {
    if (!loggedIn || !user) {
      setAvatarUrl(undefined)
      setAvatarPreset(null)
      setShowLeaderboards(true)
      return
    }
    let live = true
    const load = () => {
      getSettings()
        .then((settings) => {
          if (!live) return
          const resolved = resolveAvatar(user, settings)
          setAvatarPreset(resolved.preset ?? null)
          setAvatarUrl(resolved.url)
          setShowLeaderboards(settings.leaderboard_privacy !== 'hidden')
        })
        .catch(() => {
          // The generated seed avatar is still a usable fallback, and a failed
          // settings read cannot be evidence that leaderboards were turned off.
          if (live) setShowLeaderboards(true)
        })
    }
    load()
    window.addEventListener('settings-changed', load)
    return () => {
      live = false
      window.removeEventListener('settings-changed', load)
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

  useEffect(() => {
    if (!loggedIn) {
      setNotificationUnread(0)
      return
    }
    let live = true
    const load = () => {
      listMyReports()
        .then((data) => {
          if (live) setNotificationUnread(data.unread)
        })
        .catch(() => {
          // A failed count must not take the shell down.
        })
    }
    load()
    window.addEventListener('report-updates-changed', load)
    return () => {
      live = false
      window.removeEventListener('report-updates-changed', load)
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
          showLeaderboards={overrides.showLeaderboards ?? showLeaderboards}
          notificationUnread={overrides.notificationUnread ?? notificationUnread}
        />
        <SidebarInset>
          <div className="mb-2 flex justify-start md:hidden px-4 sm:px-6 pt-8">
            <SidebarTrigger variant="ghost" className="flex-col gap-0.5">
              <span aria-hidden="true" className="font-mono text-sm text-mathua-blue">λ</span>
              <ChevronsRight size={10} className="text-mathua-muted" />
            </SidebarTrigger>
          </div>
          {/* A block child keeps page wrappers with mx-auto centred and full
              width. As flex items their auto margins would shrink-wrap them.
              It also fills the inset and pushes a direct footer to the bottom,
              so a short page never leaves empty background under the footer. */}
          <div className="flex w-full min-w-0 flex-1 flex-col [&>*]:w-full [&>footer]:mt-auto">{children}</div>
        </SidebarInset>
      </SidebarProvider>
      <BottomTabs />
    </>
  )
}
