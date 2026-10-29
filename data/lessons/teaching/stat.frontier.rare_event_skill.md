# Rare-event estimation as a skill

In an agent-first flow, a rare-event estimator is not a script someone runs by hand. It is a skill: a packaged engine with a manifest, a stable interface to the circuit simulator, and golden statistical tests, which an agent selects, configures and calls. Treating it that way forces precise answers to what the estimator assumes, what can be changed without risk, and what evidence a new variant must produce before anyone trusts its numbers.

## What the manifest must say

An agent chooses a skill by reading manifests, so a rare-event manifest must state what decides fitness for a design: the failure-region structures it handles (one region or several), the dimension range it was validated on, the sigma range it targets, the simulator interface it needs, and the outputs it returns. Those outputs must include more than a point estimate: a confidence interval, an effective sample size and the number of simulator calls spent, so a downstream agent can decide whether to spend more budget or switch method. Faced with a bitcell that fails by read upset in one direction and write failure in another, an agent reading manifests picks a method that advertises multiple-region support over a single mean-shift sampler, a cheap text decision made before any expensive simulation.

## Where each change lands

The four-layer model sorts change requests. The **core** computes likelihood-ratio weights, the unbiased estimator and its variance; a mistake there corrupts every result, so edits are Tier 3 with human review. The **strategy** layer holds proposal distributions behind a sampler seam: a mixture proposal with one shift per failure region is a Tier 2 plug-in. **Parameters**, such as the sigma target, the sample budget and the confidence level, are Tier 1 configuration. **Hooks** expose the running estimate, its interval and the weight statistics. The sampler seam must pass the parameter dimension and distributions, known failure points and the budget, in batches of samples rather than one at a time, so the simulator can parallelize and the hot path stays fast.

## Golden statistical tests

A rare-event skill cannot be tested by comparing to yesterday's number digit for digit: a correct change to the sampler legitimately moves the estimate. Golden tests use synthetic limit states with known answers, such as a hyperplane in $d$ dimensions where $p = Q(\beta)$ exactly, a union of two hyperplanes, and a curved boundary computed by quadrature. The test asserts that across many seeds the intervals cover the truth at close to the stated rate, and that the error is within a band whose false-alarm rate is known. A differential test on a reference bitcell then checks that a new variant agrees with the baseline within their combined intervals at lower cost. Fixed seeds make every run reproducible, so a failing case can be rerun exactly.

## Going deeper

Kanj, Joshi and Nassif, DAC 2006; Singhee and Rutenbar, IEEE TCAD, 2009; the agent-first engine material in this curriculum's agentic domain.
