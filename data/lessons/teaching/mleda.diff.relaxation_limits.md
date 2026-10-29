# Limits of continuous relaxation

Differentiable placement and timing work on a continuous relaxation of a discrete problem: cells at real coordinates, smooth wirelength, smooth max, density as a potential field. The relaxation is what makes gradients exist, and it is also what a tool must undo. Legal rows and sites, discrete orientations and exact timing come back afterwards, through exact algorithms.

## What gradients cannot express

Placement legality is discrete: cells sit on rows and sites without overlap. A density penalty discourages overlap on average but cannot produce site-aligned positions, so a legalizer such as Abacus must run after global placement. Orientation, row choice and cell flipping are categorical decisions with no useful gradient. HPWL has gradients only at extreme pins, which is why smooth wirelength models exist, and an overlap count is piecewise constant, which is why density is modeled as an electrostatic field. Each smooth model is a design choice that keeps optimization possible by changing the problem slightly.

## What smoothing changes

Every smoothing parameter trades fidelity for gradient quality. Log-sum-exp overestimates a max by up to $\gamma \ln n$, so a smoothed timer is pessimistic and can rank paths differently from STA. WA wirelength underestimates span. Placers anneal these parameters toward the exact metric, but the gradient at each step optimizes the model, not the signoff number. If the delay model behind the gradients differs from signoff (pre-route wire estimates against extracted parasitics), a step that improves smoothed TNS can worsen real TNS. The relaxation must be validated by evaluating the true metric periodically.

## Recovering a discrete solution

After continuous optimization, a flow legalizes, then runs detailed placement with discrete moves such as swaps and in-row reordering, and only then routes and signs off. If global placement leaves regions overfull, legalization displaces cells far from their optimized positions and destroys the wirelength the gradients won, so the density target and penalty schedule matter as much as the optimizer. Gradient descent on a non-convex objective guarantees only progress toward a stationary point. In the four-layer engine model the differentiable optimizer is a strategy or fast mode; legality and final positions stay with the core.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4, on global placement, legalization and detailed placement; Spindler, Schlichtmann and Johannes, "Abacus: Fast Legalization of Standard Cell Circuits with Minimal Movement" (ISPD 2008).
