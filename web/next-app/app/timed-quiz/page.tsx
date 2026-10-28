'use client'

import Link from 'next/link'
import AppShell from '../../components/app-shell'
import Footer from '../../components/Footer'
import TimedQuizHost from './TimedQuizHost'

export default function TimedQuizPage() {
  return (
    <AppShell>
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          <span className="flex mb-4">
            <Link href="/profile" className="text-mathua-secondary text-sm hover:text-mathua-primary">
              ← Back
            </Link>
          </span>
          <TimedQuizHost />
          <p className="text-mathua-muted text-xs font-mono text-center mt-6">
            A timed quiz measures fluency across concepts you have already studied. It is separate from the mastery check.
          </p>
        </section>
      </div>
      <Footer />
    </AppShell>
  )
}
