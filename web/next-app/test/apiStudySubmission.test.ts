import { describe, expect, it, vi, beforeEach } from 'vitest'

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return {
    ...mod,
    authedFetch: vi.fn(),
    getAuthHeaders: () => ({}),
    getGuestId: () => null,
  }
})

import { authedFetch } from '../lib/auth'
import { submitStudyAnswer } from '../lib/api'

const ok = () =>
  new Response(JSON.stringify({ correct: true, feedback: 'Correct', xp: 1 }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })

describe('submitStudyAnswer submission id', () => {
  beforeEach(() => {
    vi.mocked(authedFetch).mockReset()
    vi.mocked(authedFetch).mockImplementation(async () => ok())
  })

  it('sends one non-empty submission id per logical answer', async () => {
    await submitStudyAnswer('a', '4', 5, '2+2')
    await submitStudyAnswer('a', '4', 5, '2+2')

    expect(authedFetch).toHaveBeenCalledTimes(2)
    const bodies = vi.mocked(authedFetch).mock.calls.map((call) => {
      const init = call[1] as RequestInit
      return JSON.parse(String(init.body)) as { submission_id?: string }
    })
    expect(bodies[0].submission_id).toBeTruthy()
    expect(bodies[1].submission_id).toBeTruthy()
    expect(bodies[0].submission_id).not.toBe(bodies[1].submission_id)
    expect(vi.mocked(authedFetch).mock.calls[0][1]).toMatchObject({ retryAttempts: 1 })
  })
})
