import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import QuizHost from '../app/goals/QuizHost'
import { startQuizSession } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    startQuizSession: vi.fn(),
    submitQuizAnswer: vi.fn(),
    skipQuizQuestion: vi.fn(),
  }
})

vi.mock('next/link', () => ({
  default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...rest}>{children}</a>
  ),
}))

vi.mock('sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }))

// The quiz screen itself is not under test; these keep the intro render light.
vi.mock('../components/KatexContent', () => ({ default: ({ text }: { text: string }) => <span>{text}</span> }))
vi.mock('../components/math/MathInput', () => ({ MathAnswerInput: () => <input /> }))

// A failed request and an empty curriculum used to render the same screen.
// The engine falls back to the first five concepts in DAG order when nothing
// is picked, so "No questions available" was almost never a true statement:
// it was what a network error, a 500 or an expired token looked like.
describe('QuizHost failure state', () => {
  beforeEach(() => {
    vi.mocked(startQuizSession).mockReset()
  })

  it('reports a failed start instead of claiming there are no questions', async () => {
    vi.mocked(startQuizSession).mockRejectedValue(new Error('network down'))
    render(<QuizHost />)

    fireEvent.click(screen.getByRole('button', { name: /start quiz/i }))

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toMatch(/quiz could not start/i)
    expect(alert.textContent).toMatch(/request failed/i)
    expect(screen.queryByText(/no questions available/i)).toBeNull()
  })

  it('says the mastery check is not due instead of blaming the request', async () => {
    vi.mocked(startQuizSession).mockRejectedValue(
      Object.assign(new Error('conflict'), { status: 409, serverMessage: 'mastery check is not due yet' }),
    )
    render(<QuizHost />)

    fireEvent.click(screen.getByRole('button', { name: /start quiz/i }))

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toMatch(/not due yet/i)
    expect(alert.textContent).not.toMatch(/request failed/i)
    expect(screen.getByText(/mastery check not due yet/i)).toBeTruthy()
    expect(screen.queryByText(/no questions available/i)).toBeNull()
  })

  it('explains that no concepts are ready instead of blaming the request', async () => {
    vi.mocked(startQuizSession).mockRejectedValue(
      Object.assign(new Error('unprocessable'), {
        status: 422,
        serverMessage: 'no concepts are ready for a mastery check',
      }),
    )
    render(<QuizHost />)

    fireEvent.click(screen.getByRole('button', { name: /start quiz/i }))

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toMatch(/no concepts ready yet/i)
    expect(alert.textContent).not.toMatch(/request failed/i)
    expect(screen.getByText(/no mastery check ready yet/i)).toBeTruthy()
    expect(screen.queryByText(/no questions available/i)).toBeNull()
  })

  it('offers a retry that can succeed', async () => {
    vi.mocked(startQuizSession).mockRejectedValueOnce(new Error('network down'))
    render(<QuizHost />)
    fireEvent.click(screen.getByRole('button', { name: /start quiz/i }))
    await screen.findByRole('alert')

    vi.mocked(startQuizSession).mockResolvedValueOnce({
      session_id: 's1',
      concept_id: 'arith.add.single',
      concept_name: 'Add single digits',
      question: '2 + 3',
      grading_type: 'numeric',
      done: false,
      questions_total: 5,
      time_limit_seconds: 30,
      closed_book: true,
    } as Awaited<ReturnType<typeof startQuizSession>>)

    fireEvent.click(screen.getByRole('button', { name: /try again/i }))
    expect(await screen.findByText(/quiz question 1 of 5/i)).toBeTruthy()
  })
})
