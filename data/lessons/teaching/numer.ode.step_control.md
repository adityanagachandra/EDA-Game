# LTE-based timestep control

A transient simulator never uses a fixed timestep. After each step it compares the LTE estimate with a tolerance, rejects bad steps, and chooses the next step from the error. Good step control is what makes a run both accurate at sharp edges and fast through quiet stretches, and it interacts with breakpoints, Newton convergence and integration order.

## Accepting or rejecting a step

After solving a step, the controller computes the LTE estimate, typically per unknown and scaled by a mixed tolerance $\text{tol} = \text{reltol}\cdot|x| + \text{abstol}$, and takes the worst ratio $r = \text{err}/\text{tol}$. If $r \le 1$ the step is accepted; otherwise it is rejected and retried with a smaller step from the last accepted point.

Example: tolerance 0.5 mV and estimates 0.31, 0.62, 0.48, 0.55 mV on four attempts: two are rejected. Every rejection wastes a full Newton solve, so controllers aim for $r$ somewhat below 1. A step can also be rejected because Newton failed to converge; then the step is usually cut by a fixed factor such as 8, as in SPICE, rather than by the error formula.

## The step-size update formula

For a method of order $p$ the LTE behaves like $Ch^{p+1}$. To bring the error to the tolerance, choose
$$h_{new} = s\,h\,\Big(\frac{\text{tol}}{\text{err}}\Big)^{1/(p+1)},$$
with a safety factor $s$ around 0.9. For trapezoidal ($p = 2$), $h = 10$ ps, err $= 2$ mV and tol $= 0.5$ mV: $h_{new} = 0.9 \cdot 10 \cdot 0.25^{1/3} = 5.67$ ps. For backward Euler with the same numbers the exponent is $1/2$: $h_{new} = 4.5$ ps.

The same formula grows the step after a very accurate one. Using the wrong exponent, for example $1/p$ instead of $1/(p+1)$, makes the controller oscillate between rejections and too-small steps.

## Breakpoints and step growth limits

Formula alone is not enough. Controllers cap growth, often at a factor of 2 per step, because the error estimate after a quiet stretch can be unreliably small. They also never step across a breakpoint, a time where a source has a corner (a PWL vertex, a pulse edge): integrating across it would smear a discontinuity that the methods assume away.

So $h_{next} = \min(h_{new},\ 2h,\ t_b - t)$. Example: $h = 50$ ps, formula gives 140 ps, next breakpoint 70 ps away: $h_{next} = \min(140, 100, 70) = 70$ ps. After the breakpoint, SPICE restarts with a small step and a backward Euler step or two, then lets the step grow again. A minimum step guards against infinite loops: hitting it is reported as "timestep too small", usually a sign of a model discontinuity.

## Going deeper

Hairer, Norsett and Wanner, "Solving Ordinary Differential Equations I" (2nd ed., Springer, 1993), section II.4; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
