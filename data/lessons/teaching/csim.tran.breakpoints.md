# Breakpoints and step restarts

Step control assumes smooth waveforms. Circuits are full of corners: PWL and PULSE sources, switching events, the start of the simulation. Simulators handle them with breakpoints, time points the step must land on exactly, and with special treatment right after them. Initial conditions are the other boundary of a transient run, and SPICE offers several ways to set them, each with different meaning.

## Breakpoints at source corners

A PWL corner or the start and end of a PULSE edge is a discontinuity in the derivative of a source. Stepping across it smears the corner and corrupts the divided-difference error estimates, which assume smoothness, so the step controller would waste rejections finding it. SPICE instead registers every source corner as a breakpoint and shortens any step that would cross one. A PULSE with a 1 ns period has four corners per period, so a 100 ns run carries about 400 breakpoints from that one source. Decks with many fast, unsynchronized sources pack breakpoints densely and force small steps everywhere, which is one motivation for event-driven fast SPICE.

## Restarting after a discontinuity

At a breakpoint the past history no longer describes the waveform, so simulators restart integration: they cut the step to a small value, reset the order of multistep methods, and often take a backward-Euler step, whose L-stability damps the fast modes the edge excites. Discontinuities that are not known in advance, such as a switch model changing state when its control crosses a threshold, are found by event location: detect the crossing in a step, interpolate or bisect to find its time, and place a breakpoint there. When Newton fails at a new time point, the step is rejected and retried smaller; when the step falls below the minimum, the run aborts with a timestep-too-small error.

## Initial conditions and UIC

By default, transient starts from the DC operating point, which is consistent with every algebraic constraint. With uic on the .tran card the operating point is skipped and integration starts from the given initial conditions, with unspecified values at zero; inconsistent conditions then produce a jump on the first step that can look like a real transient. Without uic, .ic values are held during the operating-point solve and released at $t = 0$. A .nodeset is weaker still: only an initial guess for the operating point. An oscillator started exactly at its unstable equilibrium never oscillates in a deterministic simulator, so designers kick it with an initial condition.

## Going deeper

The ngspice manual, sections on .tran, .ic and .nodeset; Najm, "Circuit Simulation" (Wiley, 2010), chapter 5.
