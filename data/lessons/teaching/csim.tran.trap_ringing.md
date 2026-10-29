# Trapezoidal ringing

The trapezoidal rule is A-stable, yet on stiff circuits it produces a sawtooth that alternates in sign from step to step. The cause is its amplification factor, which tends to -1 for very fast modes instead of to 0. Recognizing this ringing, and knowing why simulators switch to backward Euler at breakpoints or offer Gear methods, is part of reading transient waveforms critically.

## The trapezoidal amplification factor

For the test mode $x' = -x/\tau$, a trapezoidal step multiplies the solution by
$$R = \frac{1 - h/2\tau}{1 + h/2\tau}.$$
For $h \ll \tau$ this is close to $e^{-h/\tau}$. For $h = 2\tau$ it is exactly 0. For $h > 2\tau$ it is negative, so the computed value flips sign every step, and as $h/\tau \to \infty$ it approaches -1: the fast mode barely decays. With $\tau = 0.1$ ns and $h = 0.15$ ns, $R = 0.25/1.75 = 0.143$; with $h = 1$ ns, $R = -4/6 = -0.667$. The method stays stable, $|R| < 1$, but a fast parasitic mode that should vanish in one step lingers as an oscillation.

## Ringing over several steps

After $k$ steps the fast mode is multiplied by $R^k$. With $h/\tau = 4$, $R = -1/3$: starting from 2 V, the computed values are -0.667, 0.222, -0.074, 0.025 V, while the true value is $2e^{-16}$, effectively zero, after four steps. With $h/\tau = 40$, $R = -0.905$, and the sawtooth survives dozens of steps. On a waveform plot this appears as a zigzag on a node that should be smooth, often right after a sharp edge, on nodes with tiny parasitic capacitances or stiff inductive loops. Shrinking the step removes it, but at the cost the large step was meant to avoid.

## Damping with a backward Euler step

Backward Euler's factor $1/(1 + h/\tau)$ tends to 0 for fast modes: it is L-stable and kills stiff components in one step. SPICE-class simulators exploit this by taking a backward-Euler step right after a breakpoint, where a discontinuity has just excited the fast modes, then returning to trapezoidal. With $h/\tau = 1.5$, one backward-Euler step multiplies by 0.4 and each trapezoidal step by $0.25/1.75 = 0.143$: starting from 2.5 V, one of each gives 0.143 V. The alternative is to use a Gear method throughout, damping ringing everywhere but adding numerical damping to genuine oscillations, such as those of a high-Q resonator.

## Going deeper

Hairer and Wanner, "Solving Ordinary Differential Equations II" (Springer, 1996), on A- and L-stability; Najm, "Circuit Simulation" (Wiley, 2010), chapter 5.
