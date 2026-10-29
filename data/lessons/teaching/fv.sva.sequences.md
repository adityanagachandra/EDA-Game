# SVA sequences and repetition

SVA properties are built from sequences: patterns of signal values over consecutive cycles, joined by cycle delays and repetitions. Sampled-value functions compare a signal with its value in the previous cycle. Counting where a sequence matches on a trace is how you check that a property says what you mean, and the same matching semantics is what formal tools compile into automata or monitor circuits.

## Sampled-value functions

The system functions rose, fell and stable (each written with a leading dollar sign in SVA source) compare the value of a signal sampled at the current clock edge with the value sampled at the previous edge. rose(a) is true when $a$ goes from 0 to 1, fell(a) when it goes from 1 to 0, and stable(a) when it is unchanged. At the first sample there is no previous one, so a starting value must be assumed; in this course it is 0. Take $a$ = 0110 on cycles 0 to 3. rose(a) is true at cycle 1; fell(a) at cycle 3; stable(a) at cycle 0 (0 against the assumed 0) and cycle 2 (1 against 1), so twice. These functions express edge-triggered protocol rules, for example \"after rose(req), req stays high until ack\", without extra registers in the RTL.

## Concatenation with cycle delays

The sequence a ##1 b matches starting at cycle $t$ if $a$ is 1 at $t$ and $b$ is 1 at $t + 1$; a ##2 b requires $b$ at $t + 2$ with any value in between; a ##1 b ##1 c needs $a$, $b$ and $c$ on three consecutive cycles. To count matching start cycles, slide the pattern along the trace and require it to end within the trace. With $a$ = 01110110000 and $b$ = 01001011111 (cycles 0 to 10), a ##1 b matches where $a[t] = 1$ and $b[t + 1] = 1$. $a$ is 1 at cycles 1, 2, 3, 5 and 6; $b$ at the next cycle is 0, 0, 1, 1 and 1. So matches start at cycles 3, 5 and 6: 3 matches. In an implication such as s |-> p, every start cycle where $s$ matches launches an attempt of $p$.

## Consecutive repetition

b[*2] means $b$ on two consecutive cycles, and b[*1:2] means one or two consecutive cycles of $b$; so a ##1 b[*2] ##1 c needs $a$ at $t$, $b$ at $t + 1$ and $t + 2$, and $c$ at $t + 3$. A range creates several possible matches from one start: a ##1 b[*1:2] ##1 c matches from $t$ if either $b$ at $t + 1$ and $c$ at $t + 2$, or $b$ at $t + 1$ and $t + 2$ and $c$ at $t + 3$. Count a start cycle once if at least one alternative matches within the trace. With $a$ = 100100, $b$ = 011010 and $c$ = 000101: from cycle 0, $b$ at 1 and 2, then $c$ at 3 matches the two-repetition branch; from cycle 3, $b$ at 4 then $c$ at 5 matches the one-repetition branch. Two start cycles match. Formal tools turn each such sequence into a small automaton whose states track the partial matches.

## Going deeper

IEEE Std 1800 (SystemVerilog), the chapter on assertions; Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015).
