'use client'

import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import AppShell from '@/components/app-shell'
import Footer from '@/components/Footer'
import KatexContent from '@/components/KatexContent'
import { listMyReports, markMyReportsRead, type QuestionReport } from '@/lib/api'
import { useAuthState } from '@/hooks/useAuthState'

function statusLabel(status: string): string {
  switch (status) {
    case 'reviewing':
      return 'Reviewing'
    case 'resolved':
      return 'Resolved'
    case 'dismissed':
      return 'Dismissed'
    default:
      return 'Open'
  }
}

function statusTone(status: string): string {
  switch (status) {
    case 'resolved':
      return 'border-[var(--accent-green)]/40 text-[var(--accent-green)]'
    case 'dismissed':
      return 'border-mathua-border text-mathua-muted'
    case 'reviewing':
      return 'border-[var(--accent-blue)]/40 text-[var(--accent-blue)]'
    default:
      return 'border-[var(--accent-red)]/40 text-[var(--accent-red)]'
  }
}

function fmtDate(value?: string): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function MyReportsPage() {
  const { loggedIn, user } = useAuthState()
  const [reports, setReports] = useState<QuestionReport[]>([])
  const [unread, setUnread] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await listMyReports()
      setReports(data.reports)
      setUnread(data.unread)
      if (data.unread > 0) {
        // A failed read marker must not hide the history that just loaded. The
        // badge simply stays until the next successful visit.
        try {
          await markMyReportsRead()
          window.dispatchEvent(new Event('report-updates-changed'))
        } catch {
          // Keep the unread marker and the list.
        }
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Your reports could not be loaded.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (loggedIn) void load()
    else setLoading(false)
  }, [loggedIn, load])

  return (
    <AppShell>
      <div className="bg-mathua-bg text-mathua-primary">
        <main className="max-w-container mx-auto px-4 md:px-6 py-10 pb-28">
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div>
            <h1 className="font-mono text-xl text-mathua-primary">Notifications</h1>
            <p className="mt-1 text-sm text-mathua-secondary">
              Updates on the problems you reported and what moderation decided.
            </p>
          </div>
          <Link
            href="/contact"
            className="font-mono text-[11px] text-mathua-muted underline underline-offset-2 hover:text-mathua-blue"
          >
            Report another problem
          </Link>
        </div>

        {unread > 0 && (
          <p className="mt-6 border border-[var(--accent-blue)]/40 bg-mathua-surface px-4 py-3 font-mono text-[11px] text-mathua-secondary">
            {unread === 1 ? '1 report has a new verdict.' : `${unread} reports have new verdicts.`} The
            unread marker clears after this page opens.
          </p>
        )}

        {!loggedIn && (
          <p className="mt-8 border border-mathua-border bg-mathua-surface p-4 text-sm text-mathua-secondary">
            <Link href="/login" className="text-mathua-blue underline underline-offset-2">
              Sign in
            </Link>{' '}
            to track the reports you filed.
          </p>
        )}

        {loggedIn && loading && <p className="mt-8 font-mono text-xs text-mathua-muted">Loading reports.</p>}

        {error && (
          <p className="mt-8 border border-[var(--accent-red)]/40 bg-mathua-surface p-4 font-mono text-xs text-mathua-secondary">
            {error}
          </p>
        )}

        {loggedIn && !loading && !error && reports.length === 0 && (
          <p className="mt-8 border border-mathua-border bg-mathua-surface p-4 text-sm text-mathua-secondary">
            You have not reported a problem yet. A report is the fastest way to get a question or
            explanation fixed, and it keeps the exact answer and format you saw.
          </p>
        )}

        <ul className="mt-8 space-y-3">
          {reports.map((report) => (
            <li key={report.id} className="border border-mathua-border bg-mathua-surface p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="min-w-0">
                  <span className="font-mono text-xs text-mathua-primary">
                    {report.concept_id || 'General report'}
                  </span>
                  <span className="ml-2 font-mono text-[11px] text-mathua-muted">
                    {report.kind} · {report.reason}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  {report.unread && <span className="h-2 w-2 rounded-full bg-mathua-blue" aria-label="Unread update" />}
                  <span className={`border px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider ${statusTone(report.status)}`}>
                    {statusLabel(report.status)}
                  </span>
                  <span className="font-mono text-[10px] text-mathua-muted">#{report.id}</span>
                </div>
              </div>

              {report.question && (
                <div className="mt-3 border border-mathua-border bg-mathua-code p-3">
                  <KatexContent className="text-sm text-mathua-primary">{report.question}</KatexContent>
                </div>
              )}
              {report.detail && (
                <p className="mt-3 border-l-2 border-mathua-border pl-3 font-mono text-[11px] text-mathua-secondary">
                  {report.detail}
                </p>
              )}
              <p className="mt-3 font-mono text-[10px] text-mathua-muted">
                Sent {fmtDate(report.created_at)}
                {report.resolved_at ? ` · ${statusLabel(report.status)} ${fmtDate(report.resolved_at)}` : ''}
              </p>
              {report.resolution && (
                <p className="mt-2 border border-mathua-border bg-mathua-code p-3 font-mono text-[11px] text-mathua-secondary">
                  <span className="text-mathua-muted">Moderator: </span>
                  {report.resolution}
                </p>
              )}
            </li>
          ))}
        </ul>

        {loggedIn && user && (
          <p className="mt-8 font-mono text-[10px] text-mathua-muted">
            Reports are tied to @{user.username || user.name}. Signing out does not remove them.
          </p>
        )}
      </main>
      </div>
      <Footer />
    </AppShell>
  )
}
