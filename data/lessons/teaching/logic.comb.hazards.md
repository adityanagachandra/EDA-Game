# Static hazards and glitches

A combinational circuit computes the right function once it settles, but on the way there unequal path delays can make its output glitch: a momentary pulse to the wrong value. In synchronous design glitches are tolerated, because the register samples only after settling. They still matter for power, for asynchronous control signals, for clocks and resets that must never glitch, and for anyone reading a simulation waveform.

## How a static-1 hazard happens

Take $f = a b + \bar{a} c$ with $b = c = 1$. Logically $f = a + \bar{a} = 1$ whatever $a$ does. In the circuit, $\bar{a}$ comes from an inverter and arrives later than $a$. When $a$ falls from 1 to 0, the term $ab$ turns off at once, but $\bar{a}c$ turns on one inverter delay later. For that interval both terms are 0 and the output dips to 0: a static-1 hazard, a glitch on an output that should stay at 1.

The dual, a static-0 hazard, is a momentary 1 on an output that should stay 0. A dynamic hazard is several transitions where one was expected, and needs at least three paths of different delays.

## Covering adjacent implicants removes the hazard

In a sum-of-products circuit, a static-1 hazard on a single-input change occurs when the two adjacent minterms are covered by *different* product terms and no single term covers both. The fix is the consensus term: for $ab + \bar{a}c$, add $bc$. When $b = c = 1$, $bc$ holds the output at 1 during the transition, so the glitch disappears. The function is unchanged, since $bc$ is the consensus of the other two terms.

A two-level SOP that contains every prime implicant is free of static-1 hazards for single-input changes. That conflicts with logic minimization, which deliberately drops redundant primes. Hazard-free logic is therefore a separate synthesis mode, needed in asynchronous design but not in synchronous design.

## When glitches matter in a synchronous design

A flip-flop samples its D input only around the clock edge, after the combinational logic has settled, so glitches on data paths do not change the computed state. Static timing analysis does not model them at all. This is one of the strongest arguments for synchronous design.

Glitches still cost something. Each spurious transition charges and discharges capacitance, so glitch power can be a sizable share of dynamic power in deep arithmetic logic such as array multipliers, where many paths of different lengths reconverge. And a glitch on a signal used as a clock, an asynchronous reset or a gated clock enable is a functional bug, because those inputs act immediately. That is why clock-gating cells contain a latch, and why lint tools forbid combinational logic on asynchronous resets.

## Going deeper

Unger, "Asynchronous Sequential Switching Circuits" (1969); Mano and Ciletti, "Digital Design" (5th ed., 2013), the section on hazards; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), on glitching power.
