import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import MyReportsPage from '../app/my-reports/page'
import { listMyReports, markMyReportsRead } from '../lib/api'
import type { QuestionReport } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  listMyReports: vi.fn(),
  markMyReportsRead: vi.fn(),
}))

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({
    loggedIn: true,
    user: { student_id: 'stu_1', name: 'Ada Example', username: 'ada' },
  }),
}))

vi.mock('../components/Header', () => ({ default: () => <header /> }))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/BottomTabs', () => ({ default: () => <nav /> }))
vi.mock('../components/KatexContent', () => ({
  default: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}))

function report(over: Partial<QuestionReport> & { id: number }): QuestionReport {
  return {
    reporter_id: 'stu_1',
    concept_id: 'algebra.linear',
    kind: 'question',
    question: '2x = 4',
    expected: '2',
    explanation: '',
    lesson_id: '',
    source: 'practice',
    session_id: '',
    attempt_id: '',
    reason: 'wrong_answer',
    detail: '',
    status: 'open',
    resolution: '',
    resolved_by: '',
    resolved_at: '',
    created_at: '2026-10-20T00:00:00Z',
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(markMyReportsRead).mockResolvedValue(undefined)
})

describe('the learner report history', () => {
  it('shows a verdict and the moderator reason', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 1,
      reports: [
        report({
          id: 7,
          status: 'resolved',
          resolution: 'The expected answer was corrected.',
          resolved_at: '2026-10-25T12:00:00Z',
          unread: true,
        }),
      ],
    })

    render(<MyReportsPage />)

    expect(await screen.findByText('Resolved')).toBeTruthy()
    expect(screen.getByText('The expected answer was corrected.')).toBeTruthy()
    expect(screen.getByText(/1 report has a new verdict/)).toBeTruthy()
    expect(markMyReportsRead).toHaveBeenCalledTimes(1)
  })

  it('marks the history read once the page opens it', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 2,
      reports: [report({ id: 1, unread: true }), report({ id: 2, unread: true })],
    })
    render(<MyReportsPage />)

    await waitFor(() => expect(markMyReportsRead).toHaveBeenCalledTimes(1))
    expect(screen.getByText(/2 reports have new verdicts/)).toBeTruthy()
  })

  it('still shows the history when the read marker fails', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 1,
      reports: [
        report({
          id: 9,
          status: 'resolved',
          resolution: 'Fixed.',
          resolved_at: '2026-10-25T12:00:00Z',
          unread: true,
        }),
      ],
    })
    vi.mocked(markMyReportsRead).mockRejectedValue(new Error('marker failed'))

    render(<MyReportsPage />)

    expect(await screen.findByText('Resolved')).toBeTruthy()
    expect(screen.getByText('Fixed.')).toBeTruthy()
    await waitFor(() => expect(markMyReportsRead).toHaveBeenCalledTimes(1))
    expect(screen.queryByText(/marker failed/)).toBeNull()
  })

  it('explains the empty state without pretending there is history', async () => {
    vi.mocked(listMyReports).mockResolvedValue({ unread: 0, reports: [] })
    render(<MyReportsPage />)

    expect(await screen.findByText(/You have not reported a problem yet/)).toBeTruthy()
    expect(markMyReportsRead).not.toHaveBeenCalled()
  })
})
