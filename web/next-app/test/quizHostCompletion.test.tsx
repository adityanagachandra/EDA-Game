import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import QuizHost from '../app/goals/QuizHost'
import { startQuizSession, submitQuizAnswer } from '../lib/api'

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

vi.mock('../components/KatexContent', () => ({
  default: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}))

vi.mock('../components/math/MathInput', () => ({
  MathAnswerInput: ({
    value,
    onChange,
  }: {
    value?: string
    onChange?: (value: string) => void
  }) => (
    <input
      aria-label="answer"
      value={value ?? ''}
      onChange={event => onChange?.(event.target.value)}
    />
  ),
}))

const start = {
  session_id: 's1',
  concept_id: 'a',
  concept_name: 'A',
  question: '2 + 3',
  grading_type: 'numeric',
  done: false,
  questions_total: 2,
  closed_book: true,
} as Awaited<ReturnType<typeof startQuizSession>>

const secondQuestion = {
  done: false,
  correct: true,
  feedback: 'Correct!',
  xp: 5,
  concept_id: 'b',
  concept_name: 'B',
  question: '3 + 4',
  grading_type: 'numeric',
  questions_total: 2,
} as Awaited<ReturnType<typeof submitQuizAnswer>>

function finalAnswer(retakeAvailable: boolean) {
  return {
    done: true,
    correct: true,
    feedback: 'Correct!',
    xp: 5,
    retake_available: retakeAvailable,
  } as Awaited<ReturnType<typeof submitQuizAnswer>>
}

async function finishQuiz(retakeAvailable: boolean) {
  vi.mocked(startQuizSession).mockResolvedValue(start)
  vi.mocked(submitQuizAnswer)
    .mockResolvedValueOnce(secondQuestion)
    .mockResolvedValueOnce(finalAnswer(retakeAvailable))

  render(<QuizHost />)
  fireEvent.click(screen.getByRole('button', { name: /start quiz/i }))

  expect(await screen.findByText(/quiz question 1 of 2/i)).toBeTruthy()
  fireEvent.change(screen.getByLabelText('answer'), { target: { value: '5' } })
  fireEvent.click(screen.getByRole('button', { name: /check answer/i }))
  const next = await screen.findByRole('button', { name: /next/i })
  fireEvent.click(next)

  expect(await screen.findByText(/quiz question 2 of 2/i)).toBeTruthy()
  fireEvent.change(screen.getByLabelText('answer'), { target: { value: '7' } })
  fireEvent.click(screen.getByRole('button', { name: /check answer/i }))

  expect(await screen.findByText(/quiz complete/i, {}, { timeout: 3000 })).toBeTruthy()
}

describe('QuizHost completion lifecycle', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('does not offer a retake when the server says the gate is spent', async () => {
    await finishQuiz(false)

    expect(screen.getByText(/XP awarded for each correct answer/i)).toBeTruthy()
    expect(screen.queryByText(/more for faster answers/i)).toBeNull()
    expect(screen.getByText(/keep learning to earn your next mastery check/i)).toBeTruthy()
    expect(screen.queryByText(/a retake is available now/i)).toBeNull()
    expect(screen.queryByRole('button', { name: /retake quiz/i })).toBeNull()
    expect(screen.queryByText(/retake anytime/i)).toBeNull()

    const learn = screen.getByText('Keep learning →').closest('a')
    expect(learn?.getAttribute('href')).toBe('/learn')

    const history = screen.getByText(/every question you've answered/i).closest('a')
    expect(history?.getAttribute('href')).toBe('/history')
  })

  it('offers a retake only when the server says it is available', async () => {
    await finishQuiz(true)

    expect(screen.getByText(/a retake is available now/i)).toBeTruthy()
    expect(screen.queryByText(/keep learning to earn your next mastery check/i)).toBeNull()

    const retake = screen.getByRole('button', { name: /retake quiz/i })
    fireEvent.click(retake)

    await waitFor(() => expect(startQuizSession).toHaveBeenCalledTimes(2))
    expect(await screen.findByText(/quiz question 1 of 2/i)).toBeTruthy()
  })
})
