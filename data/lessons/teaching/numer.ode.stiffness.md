# Stiff equations

Circuits are stiff: a femtofarad parasitic on a small resistor has a picosecond time constant, while the signals of interest evolve over nanoseconds or microseconds. Explicit integrators must resolve the fastest mode for stability even after it has decayed. Stiffness is why every circuit simulator is built on implicit methods and sparse solves.

## Widely separated time constants

For a linear system $x' = Ax$ with real negative eigenvalues, the stiffness ratio is
$$\frac{\max|\lambda_i|}{\min|\lambda_i|}.$$
Take $A = -\begin{pmatrix} 1000 & 2 \\ 2 & 5 \end{pmatrix}$ per ns. The symmetric matrix has eigenvalues $(a + d)/2 \pm \sqrt{((a - d)/2)^2 + b^2} = 502.5 \pm 497.504$, so -1000.004 and -4.996 per ns, a stiffness ratio near 200.

Physically: a 1 fF node on a 1 kΩ path has $\tau = 1$ ps, while a 1 pF load through 10 kΩ has $\tau = 10$ ns, a ratio of $10^4$. Real netlists with extracted parasitics routinely span ratios of $10^6$ or more.

## The explicit step limit on a stiff circuit

Forward Euler is stable only if $h \le 2/|\lambda_{max}|$. For the matrix above, $h \le 2/1000.004$ ns $= 2.0$ ps, even though the slow mode, with time constant 200 ps, is the only thing still moving after a few picoseconds.

Explicit Runge-Kutta methods have larger but still bounded stability regions, so they share the problem: the step is set by the fastest mode, not by accuracy. Any error component in the fast mode, introduced by roundoff or an input edge, grows without bound as soon as the step exceeds the limit, so the limit cannot be relaxed once the fast transient is over.

## Why implicit methods win despite the solve

Simulate to $5\tau_{slow}$ with $\tau_{fast} = 2$ ps and $\tau_{slow} = 100$ ns. Forward Euler at its limit $h = 2\tau_{fast} = 4$ ps needs $500\text{ ns}/4\text{ ps} = 125{,}000$ steps. Backward Euler, stable for any $h$, takes steps chosen for accuracy, perhaps 100 to 200 over the same interval.

Each implicit step costs a sparse factorization and solve, perhaps 10 to 100 times the cost of an explicit step, but the step count falls by a factor of about a thousand, so the implicit method wins overall. That is the central economics of circuit simulation. Fast-SPICE tools recover some explicit-style speed by partitioning the circuit and treating weakly coupled blocks with different steps, but the stiff core is always handled implicitly.

## Going deeper

Hairer and Wanner, "Solving Ordinary Differential Equations II" (2nd ed., Springer, 1996), chapter IV.1; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
