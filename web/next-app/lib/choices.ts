// Parse the multiple-choice options a generator embeds in the question stem.
//
// Most `multiple_choice` generators emit options as trailing "A) …", "B) …"
// blocks (see generator/precalculus mcProblem). Rendering them as buttons lets
// students answer by tapping instead of typing the exact option text into a
// free-text box. Returns null when the question has no parseable option list,
// so every caller falls back to the text input unchanged.

export interface ChoiceOption {
  letter: string
  text: string
}

const OPTION_RE = /\n\n([A-E])\)\s+([\s\S]*?)(?=\n\n[A-E]\)\s|$)/g

export function parseChoices(question: string): ChoiceOption[] | null {
  if (!question) return null
  const opts: ChoiceOption[] = []
  let m: RegExpExecArray | null
  OPTION_RE.lastIndex = 0
  while ((m = OPTION_RE.exec(question)) !== null) {
    const text = m[2].trim()
    if (text) opts.push({ letter: m[1], text })
  }
  if (opts.length < 2) return null
  // Letters must run A, B, C, … so we never mistake prose for options.
  for (let i = 0; i < opts.length; i++) {
    if (opts[i].letter !== String.fromCharCode(65 + i)) return null
  }
  return opts
}

// Generators often write option text as bare LaTeX, for example
// `\mathbb{Q} \cap [2,9]` without `$` delimiters. KatexContent only renders
// delimited math, so those options showed their backslashes literally
// (reports 285, 286, 287 and 288). Wrap an option that contains a backslash
// and no existing math delimiter. Plain prose options stay untouched.
export function optionMathText(text: string): string {
  if (!text.includes('\\')) return text
  if (text.includes('$') || text.includes('\\(') || text.includes('\\[')) return text
  return `\\(${text}\\)`
}

// The stem without the option block, for display above tappable options.
//
// The generator writes the options into the question text, and the options are
// also drawn as buttons, so a learner reads every choice twice: once as prose
// and once as a tap target (report #232). Only the rendered stem loses the
// list. The full question is still passed to the grader and to the report
// button, because the letter is resolved against the text the server served.
export function stemWithoutChoices(question: string): string {
  if (!parseChoices(question)) return question
  const cut = question.search(/\n\n[A-E]\)\s/)
  return cut === -1 ? question : question.slice(0, cut).trimEnd()
}
