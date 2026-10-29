# Simulator changes by tier

A SPICE-class simulator splits cleanly into the four engine layers. The immutable core is the numerical kernel: MNA assembly, sparse LU, Newton and the time loop. Policies are the integration methods and convergence strategies. Parameters are the options card. Hooks are the probes and measurements an agent reads. Assigning each change request to the smallest tier that can serve it is what keeps the kernel small and the change cheap.

## Options are Tier 1

Most simulator requests are configuration. Faster with some loss of accuracy: change reltol, trtol or the accuracy mode. Less ringing: method=gear with a chosen maximum order. Runs at five temperatures and three supplies: sweep settings in the deck. A tighter current tolerance for a low-power circuit: abstol. None of these recompiles anything; each is a reviewable one-line diff that is trivially reverted. The agent's first question for any request is whether an option already expresses it. Options still need verification, because a tolerance change moves results everywhere, but the verification is a rerun against golden values, not a code review.

## Plug-ins at the device and integrator seams

When no option suffices, the change belongs at an extension point. A proprietary transistor model enters at the device-model interface, often compiled from Verilog-A, without touching the solver. A custom convergence strategy, such as ramping supplies in a fixed order, belongs behind an operating-point homotopy interface. A new integrator needs a seam that lets it supply companion coefficients and its own error estimate and order, so step control keeps working. Seams should be coarse: a device plug-in called once per device per iteration through a virtual call can be dominated by call overhead, so batching devices of one model per call lets the plug-in vectorize and keeps the hot loop fast.

## Core edits to the solver

Tier 3 changes touch the numerical kernel: how the sparse matrix is ordered and factored, how Newton tests convergence, how the time loop accepts steps. They affect every simulation, so they require a human-approved diff and the full golden regression with QoR comparison on reference decks. A stream of Tier 3 requests for the same kind of change signals a missing seam: the varying behavior should move behind an extension point. When an agent proposes a core edit, such as a device equation written into the matrix assembly loop, the harness should steer it to the lower tier that can host the same behavior, here the device-model interface.

## Going deeper

The ngspice and Xyce documentation on options and device interfaces; McAndrew et al., "Best practices for compact modeling in Verilog-A" (IEEE Journal of the Electron Devices Society, 2015).
