# Condition and expression coverage

Branch coverage asks whether each decision went both ways. It does not ask why. A condition such as `(a && b) || c` can take both outcomes while one of its atomic conditions never mattered. **Condition coverage** and **MC/DC** (modified condition/decision coverage) look inside the decision. MC/DC comes from safety-critical software, and it is the strongest structural test of boolean logic in common use, with a clean minimum test count.

## Condition coverage

An atomic **condition** is a boolean operand that contains no boolean operator: in `(a && b) || c` the conditions are $a$, $b$ and $c$. Condition coverage requires each condition to be observed both true and false, giving $2n$ condition-value pairs for $n$ conditions. Given test vectors, scan each condition's column for a 0 and a 1. Example: for `(a && b) || c` with vectors $abc$ = 110, 011, 100, condition $a$ is seen as 1 and 0, $b$ as 1 and 0, and $c$ as 0 (in 110 and 100) and 1 (in 011). All 6 pairs are covered. With only 110 and 100, $a$ is never 0 and $c$ is never 1: 4 of 6 pairs. Condition coverage alone does not require the decision to take both outcomes; tools often combine it with branch coverage.

## Computing coverage from test vectors

Short-circuit evaluation complicates measurement: in `a && b`, $b$ is not evaluated when $a$ is 0, so some tools count only conditions that were actually evaluated. When the expression is treated as fully evaluated, the computation is a table scan. Example: decision `a || b || c` with vectors 000, 100, 010: $a$ seen 0 and 1, $b$ seen 0 and 1, $c$ seen only 0. That is 5 of 6 pairs, 83.3%. Expression coverage, a related metric, asks which rows of the expression's truth table occurred; for $n$ conditions that is up to $2^n$ rows, which is why tools usually stop at condition pairs or at MC/DC instead.

## MC/DC and its minimum test count

MC/DC demands, for each condition, a pair of tests in which only that condition changes and the decision outcome changes: evidence that the condition independently affects the result (the unique-cause form). For `a && b && c`, the vector 111 (true) paired with 011, 101 and 110 (each false) demonstrates all three conditions with 4 tests. In general a pure AND or pure OR of $n$ conditions needs exactly $n + 1$ tests, far fewer than the $2^n$ rows of the truth table, and no correct MC/DC set can be smaller than $n + 1$. For mixed expressions, find for each condition two vectors at Hamming distance 1 in that condition with different outcomes. MC/DC is required for the most critical avionics software, and it applies equally to control logic in RTL.

## Going deeper

Hayhurst, Veerhusen, Chilenski and Rierson, "A Practical Tutorial on Modified Condition/Decision Coverage", NASA/TM-2001-210876; Piziali, "Functional Verification Coverage Measurement and Analysis" (2004).
