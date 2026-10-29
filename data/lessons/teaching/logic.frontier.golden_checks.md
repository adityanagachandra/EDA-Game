# Golden and differential checks for generated RTL

When an agent changes RTL, two questions need evidence: is the new design correct against the specification (golden), and did it change only what was intended (differential)? Digital logic offers unusually sharp instruments for both: combinational and sequential equivalence checking, exhaustive simulation for small blocks, and QoR comparisons after synthesis. Knowing what each instrument proves is the difference between a trusted agent and a plausible one.

## What each check proves

- **Directed or random simulation** shows the design behaves correctly on the vectors run, and nothing about the rest. Exhaustive simulation of an $n$-input combinational block needs $2^n$ vectors: fine at 16 inputs, hopeless at 64.
- **Combinational equivalence checking** proves that two netlists with matching registers compute the same next-state and output functions for every input, using SAT and structural hashing. It is the standard check after synthesis.
- **Sequential equivalence checking** proves two designs produce the same outputs from reset for every input sequence, even when their registers differ (retiming, re-encoded FSMs).
- **Property checking** proves stated properties hold, which is only as strong as the properties.

## Differential checks for an edit

For a refactoring edit (re-encoding an FSM, replacing a ripple adder by a prefix adder, retiming a pipeline) the right check is equivalence against the previous version: behavior must not change at all. For a functional edit (add a saturate mode, change a reset value) equivalence must fail, and the check becomes: equivalent under the old mode, and the new behavior matches the specification under the new one. Constrain the comparison to the unchanged modes and prove those equivalent.

On top of correctness, a differential QoR report after synthesis (area, worst slack, cell count, power estimate) shows whether the edit had the intended physical effect, and nothing else moved.

## Pitfalls

Equivalence checkers need a correspondence or a reset state, and both can hide problems. Two FSMs equivalent from reset can differ from an illegal state, which matters if the specification requires safe recovery. A design that is equivalent only because both versions share the same bug (both copied from one wrong reference) passes every differential check; that is why the golden reference must be independent.

Don't-care handling is another trap: synthesis may exploit X assignments differently in the two versions, so an equivalence check that treats X as a free choice can pass designs that simulate differently. A trustworthy harness states its X semantics and its reset assumptions alongside every pass.

## Going deeper

Kuehlmann and Krohm, "Equivalence Checking Using Cuts and Heaps" (DAC, 1997); Mishchenko et al., "Improvements to Combinational Equivalence Checking" (ICCAD, 2006); the ABC documentation for cec and dsec.
