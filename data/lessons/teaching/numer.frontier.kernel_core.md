# Numerical kernels as an immutable core

In an agent-first EDA engine, the numerical kernels are the immutable core: the code every analysis trusts and every agent calls but rarely edits. Deciding precisely what belongs in that core for a simulator, IR-drop or thermal solver, and what is policy, parameter or hook, is the design decision that lets agents customize solvers safely.

## What belongs in the solver core

The core holds what defines correctness: sparse factorization and triangular solves, SpMV, the Krylov recurrences, the integration formulas and their LTE estimators, the Newton update, and the definitions behind every reported number, such as what a residual is and what "converged" means. A downstream timer or signoff check reads a converged waveform as a contract; if the meaning of the tolerance changed silently, every result would change meaning with it.

Two practical rules follow. First, the hot loops stay inside the core and the seams around it are coarse, called once per solve or per batch, so the compiler can inline and vectorize. Second, a change to the core, say a new LTE formula or a different residual scaling, is a Tier 3 change: it needs a human-approved diff and a full regression, because it alters the accuracy contract of every analysis.

## Policies, parameters and hooks around a solver

Around the core sit three easier-to-change layers. Policy: which ordering (AMD, nested dissection, Markowitz), which preconditioner, which integration method among those the core provides, which damping or limiting strategy, which timestep controller. Parameters: reltol, abstol, vntol, pivot threshold, GMRES restart length, maximum Newton iterations, fast or accurate mode. Hooks: per-iteration residual norms, the node with the largest update, accepted and rejected steps, pivot growth, condition estimates.

Classifying requests by layer is the daily work. "Use BDF2 for this PLL" is a policy choice; "tighten reltol to 1e-4" is a parameter; "let the agent see which node fails Newton" is a hook. None of them touches the core, so none needs a recompile of the kernel, and most need only a configuration change.

## The determinism contract of a kernel

Agents compare runs: a baseline against a candidate, a golden result against a new build. That only works if the same inputs and settings give the same answer. So the core promises determinism: fixed reduction orders, no hidden global state, explicit seeds for anything random, and results independent of thread scheduling.

Fast paths may relax this, for example GPU atomic-add reductions, but then the engine offers two declared modes, a deterministic one for comparisons and signoff and a nondeterministic one for exploration where only rough rankings matter. A reordering of floating-point operations for speed changes results in the last digits; the release process must then show tolerance-based agreement with golden results and determinism within the new version, since bitwise equality across the change is impossible.

## Going deeper

The four-layer engine model and change tiers in docs/eda/agent-first-concepts.md of this repository; Higham, "Accuracy and Stability of Numerical Algorithms" (2nd ed., SIAM, 2002), chapter 4 on summation order.
