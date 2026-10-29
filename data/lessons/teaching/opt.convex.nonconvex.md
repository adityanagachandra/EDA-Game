# Nonconvex landscapes, QCQP and local optima

Most real EDA formulations are nonconvex: non-overlap, discrete sizes, density, bilinear capacitance terms, and nonlinear equalities. Nonconvex problems can have many local optima and saddle points, and descent methods find the one whose basin contains the start. A tool designer must recognize what breaks convexity, predict where a local method will land, and choose a strategy, such as continuation, multi-start, convex relaxation or annealing, that makes a local answer good enough.

## Basins of attraction

Gradient descent with a small step follows $-f'(x)$ and stops at a stationary point that is a local minimum. Which one depends only on the basin the start lies in. Take a quartic with positive leading term and $f'(x) = (x + 6)(x + 4)(x + 1)$. The sign of $f'$ is negative below -6, positive on $(-6, -4)$, negative on $(-4, -1)$ and positive above -1. So $x = -6$ and $x = -1$ are local minima and $x = -4$ is a local maximum.

Descent never crosses the maximum: from $x_0 = -6.8$ it rises to -6, from -5 it falls back to -6, and from -3 it moves to -1. In placement the same picture holds in millions of dimensions, and the start, usually a quadratic placement, determines the basin. That is why analytical placers invest in a good initial solution.

## What breaks convexity

A feasible set is convex if it contains the segment between any two of its points. Non-overlap between rectangles is an either-or condition, a union of four convex pieces, so it is not convex. A discrete library $x \in \{1, 2, 4, 8\}$ is not convex: 3 lies between 2 and 4. A circle $x^2 + y^2 = 1$ is not convex, because a nonlinear equality is convex only when affine. A constraint $(x_i - x_j)^2 + (y_i - y_j)^2 \ge d^2$ excludes a disk, and its complement is not convex.

A quadratically constrained program is convex exactly when the objective's matrix and every constraint matrix $Q_i$ are positive semidefinite. A bilinear term $xy$, as when wire width multiplies length, has Hessian eigenvalues $\pm 1$, a saddle. Electrostatic density energy in ePlace-style placers is nonconvex in the positions, so results depend on the initial placement and on the penalty schedule.

## Strategies for nonconvex problems

Four strategies cover most EDA practice. Continuation solves an easy smooth problem first and tightens toward the hard one: analytical placers start with a small density weight and raise it, so wirelength sets the global structure before density forces spreading. Multi-start runs a local method from several starts and keeps the best, never the average, since averaging solutions of a nonconvex problem can land in a poor or infeasible region. Convex relaxation, such as SDP or McCormick envelopes for bilinear terms, gives a lower bound and sometimes a good starting point. Stochastic search such as simulated annealing accepts some uphill moves to cross barriers between basins.

Saddle points are the other trap. At a point with zero gradient and Hessian eigenvalues 3 and -2, moving along the negative-curvature eigenvector decreases the cost. In production a local optimum is good enough when its QoR is within tolerance of baselines or bounds on representative designs.

## Going deeper

Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapter 2 on local minimizers; Lu et al., "ePlace: electrostatics-based placement using fast Fourier transform and Nesterov's method" (ACM TODAES, 2015); McCormick, "Computability of global solutions to factorable nonconvex programs" (Mathematical Programming, 1976).
