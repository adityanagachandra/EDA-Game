'use client'

export interface BriefingItem {
  title: string
  detail: string
}

// Pre-test expectations, MA-grounded: diagnostics run 25–45 adaptive
// questions (~30–45 min per engine bounds), quizzes are closed-book with a
// slow-answer threshold. Stated once up front so the test stays distraction-free.
export const DIAGNOSTIC_BRIEFING: BriefingItem[] = [
  {
    title: '25–45 questions · about 30–45 minutes',
    detail: 'The test adapts to you. Clear answers mean fewer questions.',
  },
  {
    title: 'Answer from what you know, no searching',
    detail: 'Honest answers put you at the right starting point. Looked-up answers put you in the wrong place.',
  },
  {
    title: 'Speed counts',
    detail: 'Fast correct answers show mastery and shorten the test.',
  },
  {
    title: 'Pause anytime',
    detail: 'Leave and pick up where you stopped.',
  },
]

export const QUIZ_BRIEFING: BriefingItem[] = [
  {
    title: 'Up to 5 questions · closed book',
    detail: 'No countdown is shown. A slow answer is still graded. Keep books and notes shut.',
  },
  {
    title: 'Covers what you are still learning',
    detail: 'Questions come from concepts you have learned but not yet mastered.',
  },
  {
    title: 'A wrong answer points to the concept',
    detail: 'The verdict stays short during the check. Concepts worth revisiting are listed at the end and open in Learn for the full explanation.',
  },
  {
    title: 'XP for each correct answer · retake when due',
    detail: 'A retake opens again once you have earned enough XP for the next check.',
  },
]

export const REVIEW_BRIEFING: BriefingItem[] = [
  {
    title: 'Only what is due for review',
    detail: 'Questions come from concepts whose memory is fading.',
  },
  {
    title: 'Answer from what you remember, no searching',
    detail: 'Honest answers recalibrate your next review date.',
  },
  {
    title: 'Speed counts',
    detail: 'Fast correct answers show the memory is solid.',
  },
  {
    title: 'XP for each correct answer',
    detail: 'Reviews record progress and XP like practice.',
  },
]

export default function BriefingCard({ eyebrow, items }: { eyebrow: string; items: BriefingItem[] }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-4 text-left">
      <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted mb-3">
        {eyebrow}
      </p>
      <ul className="space-y-3">
        {items.map(item => (
          <li key={item.title} className="min-w-0">
            <p className="font-mono text-xs text-mathua-primary">{item.title}</p>
            <p className="font-mono text-[11px] text-mathua-muted mt-0.5 break-words [overflow-wrap:anywhere]">
              {item.detail}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}
