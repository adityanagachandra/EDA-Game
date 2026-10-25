'use client'

import { useEffect, Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useTheme } from '../../hooks/useTheme'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import Footer from '../../components/Footer'
import Loading from '../../components/Loading'
import QuizHost from './QuizHost'

// This route is the quiz. The domain picker and the goals diagnostic used to be
// rendered here as well, which meant two copies of the same screen with the same
// scope count worded differently in each -- and only one of them got fixed. They
// live in /onboard now, and a visit without ?quiz=1 goes there.
//
// The diagnostic lives in /onboard and runs inline there. An older
// DiagnosticHost in this folder was deleted on 23 October 2026: it had a test
// suite that kept passing while the flow learners actually used drifted, and
// three client fixes went into it before anyone noticed nothing rendered it. Deleting it means
// porting that coverage onto the inline implementation first, which is its own
// change and not something to do by deleting files.

function GoalsContent() {
  const { mounted } = useTheme()
  const { replace } = useRouter()
  const searchParams = useSearchParams()

  // Read the route through the router, not window.location. A Link can render
  // this component from a prefetched payload before the address bar changes,
  // and a one-shot window read then caches false and redirects a quiz click to
  // the diagnostic. useSearchParams re-renders when navigation settles.
  const isQuizMode = searchParams.get('quiz') === '1'

  useEffect(() => {
    if (!mounted || isQuizMode) return
    replace('/onboard')
  }, [mounted, isQuizMode, replace])

  if (!mounted) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  // Until the redirect lands, render nothing rather than a picker that is no
  // longer this route's job.
  if (!isQuizMode) return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />

  return (
    <>
      <Header />
      <div className="flex min-h-[calc(100svh-69px)] flex-col">
        <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0 flex-1 min-h-[calc(100svh-69px)] lg:min-h-0">
          <section className="pt-8 min-w-0 overflow-hidden">
            <QuizHost />
          </section>
        </div>
        <Footer />
      </div>
      <BottomTabs />
    </>
  )
}

export default function GoalsPage() {
  return (
    <Suspense fallback={<><div style={{ background: 'var(--bg)', minHeight: '100vh' }}><Loading label="LOADING" /></div></>}>
      <GoalsContent />
    </Suspense>
  )
}
