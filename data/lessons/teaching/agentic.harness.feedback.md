# Actionable failure feedback

When a build or test fails, the harness sends the failure back to the model. What it sends decides whether the next attempt is better. Feedback must be precise enough to act on, small enough to fit in context, and designed so that the model fixes the problem rather than learning to hide it. This is an interface design problem for the harness author.

## What a build failure should return

For a compile or link failure, return the compiler's diagnostic verbatim: file, line, column, error class and the offending code, with the first error first, because later errors are often consequences of it. Add the relevant declaration when the error refers to it, for example the extension-point interface whose method signature the plug-in got wrong. Sanitizer reports from AddressSanitizer or UndefinedBehaviorSanitizer should be returned with the stack trace trimmed to the frames in the skill's own code. Do not return a 2,000-line build log: it buries the error, wastes context and makes the model guess which line matters.

## What a test failure should return

For a golden or differential failure, return the minimum a human would debug from. Name the failing test and design, the metric that failed, the expected and actual values and the tolerance: "golden placement test dp_alu32: displacement 14.2 sites, limit 8". For a differential check, list the objects that moved unexpectedly with their before and after values, such as the five endpoints whose slack changed outside the intended set. Return the first failing check rather than every downstream consequence, and say what passed, so the model does not undo a correct part of its edit.

## Feedback that invites gaming

Some feedback teaches the wrong lesson. If the harness returns the golden test's full expected output, a model can special-case that input and pass without solving the problem; held-out tests that the model never sees catch this. If the harness reports only a single pass or fail bit, the model cannot tell what to fix and tends to make random changes. If failures in the test harness itself, such as a timeout from an overloaded machine, are reported as the model's failure, the model learns to "fix" code that was correct. Good feedback separates infrastructure errors from skill errors and keeps some tests hidden.

## Going deeper

The Clang documentation on diagnostics and sanitizers; the agent-first EDA notes in this curriculum, section on the harness loop.
