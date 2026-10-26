import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { isLoggedIn } from '../lib/auth'
import PlanPage from '../app/plan/page'

// The planner used to end at "← Back to profile", which was its only exit and a
// dead end for the visitors the assistant sends here: a guest has no profile.
// These tests pin the exit that replaced it, because a missing next step is
// invisible to every other kind of check, the page still renders and still
// saves a plan.

vi.mock('next/link', () => ({
  default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...rest}>{children}</a>
  ),
}))

vi.mock('next/navigation', () => ({
  useSearchParams: () => new URLSearchParams(''),
}))

vi.mock('../lib/auth', () => ({ isLoggedIn: vi.fn() }))

// The editor, shell and footer are irrelevant here and pull in their own data
// fetches; the page's own copy is the subject.
vi.mock('../components/PlanEditor', () => ({ default: () => <div data-testid="editor" /> }))
vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => null }))

describe('plan page next steps', () => {
  beforeEach(() => {
    vi.mocked(isLoggedIn).mockReturnValue(false)
  })

  it('offers the diagnostic and Learn to a visitor without an account', () => {
    render(<PlanPage />)
    expect(screen.getByRole('link', { name: 'Take the diagnostic' })).toHaveAttribute('href', '/onboard')
    expect(screen.getByRole('link', { name: 'Start learning' })).toHaveAttribute('href', '/learn')
    expect(screen.queryByRole('link', { name: 'Back to profile' })).toBeNull()
  })

  it('adds the profile exit only when there is a profile', () => {
    vi.mocked(isLoggedIn).mockReturnValue(true)
    render(<PlanPage />)
    expect(screen.getByRole('link', { name: 'Back to profile' })).toHaveAttribute('href', '/profile')
  })
})
