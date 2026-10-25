import { beforeEach, describe, expect, it, vi } from 'vitest'

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
import { getAdminReportImage, submitReport } from '../lib/api'

const ok = () =>
  new Response(JSON.stringify({ id: 1, status: 'open', has_image: true }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })

describe('submitReport screenshot transport', () => {
  beforeEach(() => {
    vi.mocked(authedFetch).mockReset()
    vi.mocked(authedFetch).mockImplementation(async () => ok())
  })

  it('sends a screenshot as multipart with the message in the same request', async () => {
    const shot = new File([new Uint8Array([0x89, 0x50, 0x4e, 0x47])], 'shot.png', { type: 'image/png' })
    await submitReport({
      kind: 'message',
      source: 'contact',
      reason: 'other',
      detail: 'The page breaks on a phone.',
      image: shot,
    })

    const init = vi.mocked(authedFetch).mock.calls[0][1] as RequestInit
    expect(init.body).toBeInstanceOf(FormData)
    const form = init.body as FormData
    expect(form.get('detail')).toBe('The page breaks on a phone.')
    expect(form.get('image')).toBe(shot)
    expect(init.headers).not.toHaveProperty('Content-Type')
  })

  it('keeps the original JSON request when there is no screenshot', async () => {
    await submitReport({
      kind: 'message',
      source: 'contact',
      reason: 'other',
      detail: 'The page breaks on a phone.',
    })

    const init = vi.mocked(authedFetch).mock.calls[0][1] as RequestInit
    expect(typeof init.body).toBe('string')
    expect(JSON.parse(String(init.body))).toMatchObject({
      kind: 'message',
      detail: 'The page breaks on a phone.',
    })
  })

  it('loads a moderation screenshot as a blob instead of a public image URL', async () => {
    vi.mocked(authedFetch).mockResolvedValue(new Response(new Blob(['png']), { status: 200 }))
    const blob = await getAdminReportImage(42)

    expect(blob).toBeInstanceOf(Blob)
    expect(String(vi.mocked(authedFetch).mock.calls[0][0])).toContain('/api/admin/reports/42/image')
  })
})
