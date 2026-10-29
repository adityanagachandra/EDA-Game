# Continuation and homotopy methods

When Newton fails even with damping, simulators solve an easier problem first and deform it into the real one. Continuation, or homotopy, methods do this systematically, and gmin stepping, source stepping and pseudo-transient analysis are all instances. Understanding them as one idea lets a tool designer choose and tune the right aid for a hard operating point.

## Embedding the problem in a parameter family

Define $H(x, s)$ with $H(x, 0) = 0$ easy to solve and $H(x, 1) = F(x)$. Step $s$ from 0 to 1 and, at each step, solve $H(x, s_{i+1}) = 0$ by Newton starting from the solution at $s_i$. Because the solution moves continuously with $s$, the previous point lies close to the next one, inside Newton's region of quadratic convergence.

A simple example is the convex homotopy $H(x, s) = sF(x) + (1 - s)(x - x_0)$, whose solution at $s = 0$ is $x_0$. The stepper adapts $\Delta s$: grow it after a Newton solve converges in 2 or 3 iterations, shrink it (and retry) after a failure. Natural-parameter continuation of this kind works as long as the path is a smooth curve along which $s$ keeps increasing.

## Source stepping and gmin stepping as homotopies

Source stepping scales every independent source by $s$. At $s = 0$ a circuit without sources sits at the trivial solution, all node voltages zero; ramping $s$ to 1 brings the real excitation in gradually.

Gmin stepping adds a conductance $g$ from every node to ground, starts with a large value, say $10^{-2}$ S, and reduces it by a factor (often 10) each step to the base gmin, then to its final small value. A large $g$ makes the Jacobian strongly diagonally dominant and the circuit nearly linear around zero, so Newton converges; each reduction moves back toward the real circuit.

Pseudo-transient continuation adds a capacitor from each node to ground and runs a transient to steady state: backward Euler adds $C/h$ to the Jacobian diagonal, which acts like damping, and at steady state the capacitor currents vanish, so the DC answer is unchanged.

## Step size along the path and turning points

The solution path can fold back: as $s$ increases, the curve reaches a turning point beyond which no solution exists on that branch, and continues with $s$ decreasing. At the fold the Jacobian $\partial H/\partial x$ is singular, so natural-parameter continuation stalls there however small the step.

Pseudo-arclength continuation treats $(x, s)$ together and parameterizes by arclength: it appends the equation that the new point lies a fixed distance along the path tangent. The augmented Jacobian stays nonsingular at folds, so the path can be followed around them.

Circuits with positive feedback (latches, Schmitt triggers, oscillators biased at an unstable point) have several DC solutions, and continuation returns whichever the path reaches. That is why simulators offer nodesets and initial conditions, and why a converged operating point is one valid solution, not necessarily the one the designer intended.

## Going deeper

Allgower and Georg, "Introduction to Numerical Continuation Methods" (SIAM Classics, 2003); Najm, "Circuit Simulation" (Wiley, 2010), chapter on DC convergence.
