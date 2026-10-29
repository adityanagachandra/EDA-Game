# Golden and QoR differential tests for a placer

No change to a placer, whether a parameter, a plug-in or a core edit, should be promoted on faith. Golden tests prove that kernels compute what they should; QoR differential tests compare the changed placer with a baseline across a regression suite and prove that only the intended metrics moved, by the intended amounts. Writing the gates precisely is what makes agent-generated changes safe to accept.

## Golden tests for placer kernels

Golden tests check deterministic pieces against known answers: HPWL of hand-built nets, B2B weights, density of a small bin grid, the Abacus result on a three-cell row, the legality checker on legal and illegal examples, and the incremental wirelength delta against full recomputation for random moves. They run in seconds, so the harness runs them after every edit. For a change at a seam, the plug-in's own golden tests join the suite. Determinism is a precondition: with fixed seeds and ordered reductions, the same input must give bit-identical output, otherwise a golden mismatch cannot be told apart from noise.

## Measuring QoR change

A QoR comparison runs the baseline and the changed placer on the same designs with the same downstream flow and reports metric changes. The basic number is the percent change: HPWL from 458000 to 455000 um is $-0.66\%$. A suite summary uses the geometric mean of new-over-baseline ratios, which treats a 1% change on a small design and on a large one equally and keeps one large design from dominating; for ratios 0.99, 1.01 and 0.97, the geometric mean is about 0.9899. Metrics include HPWL, overflow, legalization displacement, routed wirelength and DRC count where available, timing (WNS, TNS) and runtime, since a QoR gain that doubles runtime may not be acceptable.

## Promotion gates

A gate turns measurements into a decision. A typical one: on every regression design, HPWL may grow by at most 1% over baseline, overflow must stay at most 0.1, legality must hold, and runtime may grow by at most 10%. The change is blocked if any design fails, and the report lists which designs and metrics broke. For a targeted change, such as a rule plug-in, the gate also checks that the intended effect appears on the target design. Gates should be strict enough to catch regressions and loose enough to tolerate run-to-run noise, which is measured by running the baseline several times. Every promoted change is versioned and reversible, so a regression found later can be traced and rolled back.

## Going deeper

The agent-first EDA notes in this curriculum, section on verification of engine changes; the OpenROAD-flow-scripts metrics and regression documentation.
