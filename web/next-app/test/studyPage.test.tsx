import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, within, waitFor } from '@testing-library/react'
import StudyPage from '../app/study/page'

// `/study` is the reference and exploration surface restored by ADR-054: context-bound to one
// concept, not a catalogue, and awarding nothing. These tests pin that behaviour, an unmastered
// prerequisite is readable, the practise control is separate and eligibility-gated, the DAG
// neighbourhood is drawn from real relationships, and the original target survives a return.

let params = new URLSearchParams()
const replaceMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: replaceMock, push: vi.fn() }),
  useSearchParams: () => params,
}))

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => null }))

vi.mock('../lib/auth', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/auth')>()),
  getUserInfo: () => ({ student_id: 's1' }),
}))

const { getLessonReadiness, getLessonKPs, getLessons, getLessonBody, getProgress } = vi.hoisted(() => ({
  getLessonReadiness: vi.fn(),
  getLessonKPs: vi.fn(),
  getLessons: vi.fn(),
  getLessonBody: vi.fn(),
  getProgress: vi.fn(),
}))

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getLessonReadiness,
  getLessonKPs,
  getLessons,
  getLessonBody,
  getProgress,
}))

const CID = 'alg.factor.trinomial'

beforeEach(() => {
  params = new URLSearchParams(`concept=${CID}`)
  replaceMock.mockClear()
  getLessonKPs.mockReset()
  getLessons.mockReset()
  getLessonBody.mockReset()
  getLessonReadiness.mockReset()
  getProgress.mockReset()

  getLessonKPs.mockResolvedValue({
    concept_id: CID,
    kps: [{ label: 'Factor out the GCF', section: 'S', subgoals: ['find it'], worked_example: 'x$^2$ + x = x(x+1)' }],
    diagram: '',
  })
  getLessons.mockResolvedValue({ lessons: {} })
  getLessonBody.mockResolvedValue('')
  getProgress.mockResolvedValue({})
  getLessonReadiness.mockResolvedValue({
    concept_id: CID, ready: false, weak: [], missing: [],
    eligible: false, state: 'locked',
    prerequisites: [
      { conceptId: 'alg.factor.gcf', title: 'Factor out the greatest common factor', state: 'locked', met: false },
      { conceptId: 'alg.poly.foil', title: 'Multiply two binomials using FOIL', state: 'locked', met: false },
    ],
  })
})

describe('/study reference route', () => {
  it('renders the concept and its reference material for a locked concept', async () => {
    render(<StudyPage />)
    expect(await screen.findByText('Factor trinomials of the form x² + bx + c')).toBeTruthy()
    // The material is not gated by eligibility: the worked example renders.
    expect(await screen.findByText(/Factor out the GCF/)).toBeTruthy()
    // The header names the surface and the subject.
    expect(screen.getByText(/Reference · Algebra/)).toBeTruthy()
  })

  it('lists both actual prerequisites under Depends on, each openable', async () => {
    render(<StudyPage />)
    const depends = await screen.findByRole('region', { name: 'Depends on' })
    await waitFor(() => expect(within(depends).getAllByRole('link')).toHaveLength(2))
    const hrefs = within(depends).getAllByRole('link').map((l) => l.getAttribute('href') ?? '')
    expect(hrefs.some((h) => h.includes('/study?concept=alg.factor.gcf'))).toBe(true)
    expect(hrefs.some((h) => h.includes('/study?concept=alg.poly.foil'))).toBe(true)
  })

  it('lists real dependents under Unlocks', async () => {
    render(<StudyPage />)
    const unlocks = await screen.findByRole('region', { name: 'Unlocks' })
    // alg.factor.ac_method lists alg.factor.trinomial as a prerequisite.
    const hrefs = within(unlocks).getAllByRole('link').map((l) => l.getAttribute('href') ?? '')
    expect(hrefs.some((h) => h.includes('/study?concept=alg.factor.ac_method'))).toBe(true)
  })

  it('shows neutral copy when a concept has no neighbours', async () => {
    params = new URLSearchParams('concept=abstract.group.dihedral')
    getLessonReadiness.mockResolvedValue({
      concept_id: 'abstract.group.dihedral', ready: true, weak: [], missing: [],
      eligible: true, state: 'unlocked', prerequisites: [],
    })
    render(<StudyPage />)
    expect(await screen.findByText('No prerequisites listed')).toBeTruthy()
    expect(screen.getByText('No dependent concepts listed')).toBeTruthy()
  })

  it('forwards to /learn when there is no concept to show', async () => {
    params = new URLSearchParams('')
    render(<StudyPage />)
    await waitFor(() => expect(replaceMock).toHaveBeenCalledWith('/learn'))
  })

  it('offers a separate, eligibility-gated practice control that carries the return target', async () => {
    params = new URLSearchParams(`concept=${CID}&return=${encodeURIComponent('/learn?concept=alg.factor.ac_method')}`)
    getLessonReadiness.mockResolvedValue({
      concept_id: CID, ready: false, weak: [], missing: [], state: 'in_progress', eligible: true,
      prerequisites: [],
    })
    render(<StudyPage />)
    const practise = await screen.findByRole('link', { name: /Practise this concept/ })
    const href = practise.getAttribute('href') ?? ''
    // Entering the normal learning loop, not a parallel practice engine.
    expect(href).toContain('/learn?concept=alg.factor.trinomial')
    // The original target survives for the return.
    expect(href).toContain(`return=${encodeURIComponent('/learn?concept=alg.factor.ac_method')}`)
    // And a way back is offered.
    const backLinks = await screen.findAllByRole('link', { name: /Back to Factor trinomials/ })
    expect(backLinks.length).toBeGreaterThan(0)
    expect(backLinks[0].getAttribute('href')).toBe('/learn?concept=alg.factor.ac_method')
  })

  it('sends a mastered concept to review rather than re-teaching it', async () => {
    getLessonReadiness.mockResolvedValue({
      concept_id: CID, ready: true, weak: [], missing: [], state: 'mastered', eligible: true,
      prerequisites: [],
    })
    render(<StudyPage />)
    const practise = await screen.findByRole('link', { name: /Practise this concept/ })
    expect(practise.getAttribute('href')).toBe('/review')
  })

  it('keeps navigation and practice available when content is missing', async () => {
    getLessonKPs.mockResolvedValue({ concept_id: CID, kps: [] })
    getLessons.mockResolvedValue({ lessons: {} })
    render(<StudyPage />)
    expect(await screen.findByText(/Detailed study material is not available/)).toBeTruthy()
    // The gap is not a dead end: prerequisites and practice are still offered.
    expect(screen.getByRole('link', { name: 'Explore prerequisites' })).toBeTruthy()
    expect(screen.getAllByRole('link', { name: /Practise this concept/ }).length).toBeGreaterThan(0)
  })
})
