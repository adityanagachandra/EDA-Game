# Seams for solvers and preconditioners

A seam is the narrow, stable interface where new behaviour plugs into an engine without touching its core. For numerical engines the most valuable seams are around the linear solver and the preconditioner, because that is where customers' matrices differ most. Where the seam sits, and what context it passes, decides whether solver changes are configuration, plug-ins or core edits.

## The linear solver and preconditioner seams

A good linear-solver seam mirrors the structure of sparse direct solvers: analyze(pattern), factor(values), solve(rhs). Splitting analysis from factorization exposes the reuse that matters in simulation: the pattern is fixed across Newton iterations and timesteps, so a solver plug-in orders once and refactors many times. A preconditioner seam is setup(A) and apply(r), where apply returns an approximation to $A^{-1}r$ once per Krylov iteration.

Calls are coarse, per matrix or per solve, never per entry, so the plug-in owns its inner loops and the core's performance is preserved. Plug-ins receive plain data structures (CSR arrays, vectors), keep no hidden global state, and are deterministic given their inputs. A plug-in that caches state between calls without declaring it breaks reproducibility and golden tests.

## What context the seam must expose

A plug-in can only choose well if the seam tells it what it is solving. Useful context: structural flags (symmetric, positive definite, M-matrix), so it can pick Cholesky or CG over LU or GMRES; whether the pattern changed since the last call and how much the values changed, so it can reuse a factorization or preconditioner; the outer Newton tolerance and current nonlinear residual, so an inner iterative solve is only as accurate as needed (inexact Newton); and domain metadata such as node coordinates and metal layers, so a power-grid preconditioner can coarsen by layer.

Context should be data, not callbacks into the core. Missing context forces plug-ins to guess, for instance re-analysing every call or solving every system to $10^{-14}$, which wastes the speed the seam was meant to unlock.

## Classifying solver change requests by tier

Tier 1 is configuration: switch the IR-drop preconditioner from Jacobi to incomplete Cholesky, which the tool ships; change the GMRES restart from 30 to 60; tighten a tolerance. Tier 2 is a plug-in at an existing seam: a block-Jacobi preconditioner whose blocks are power-grid layers, or a new ordering for a family of analog netlists. Tier 3 is a core edit: declaring Newton convergence on the residual alone, changing the LTE estimator, altering how gmin is applied.

The tier is a test of the seams. If most requests over a year needed Tier 3, behaviour that varies is trapped in the core and the seam is in the wrong place. The right response is to move that variability behind an interface once, under human review, so that future requests of that kind land in Tier 1 or 2.

## Going deeper

Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapter 1 on the analyze-factor-solve structure; Balay et al., "PETSc users manual" (Argonne National Laboratory), on KSP and PC interfaces as a model of solver and preconditioner seams.
