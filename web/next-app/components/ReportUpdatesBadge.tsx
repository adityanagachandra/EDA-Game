'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { listMyReports } from '../lib/api'

// A small count beside the avatar. The page that clears it dispatches
// report-updates-changed, so the badge clears without a reload.
export default function ReportUpdatesBadge({ enabled }: { enabled: boolean }) {
  const [unread, setUnread] = useState(0)

  useEffect(() => {
    if (!enabled) {
      setUnread(0)
      return
    }
    let cancelled = false
    const load = () => {
      listMyReports()
        .then((data) => {
          if (!cancelled) setUnread(data.unread)
        })
        .catch(() => {
          // A transient failure must not put an error in the header. The page
          // itself reports the failure when the learner opens it.
        })
    }
    load()
    window.addEventListener('report-updates-changed', load)
    return () => {
      cancelled = true
      window.removeEventListener('report-updates-changed', load)
    }
  }, [enabled])

  if (!enabled || unread === 0) return null

  return (
    <Link
      href="/my-reports"
      aria-label={unread === 1 ? '1 report update' : `${unread} report updates`}
      className="flex min-h-[44px] min-w-[24px] items-center justify-center font-mono text-[10px] text-mathua-blue"
      title="Open my reports"
    >
      <span className="border border-mathua-blue px-1.5 py-0.5 leading-none">{unread}</span>
    </Link>
  )
}
