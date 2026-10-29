# Damped Newton steps

Newton converges quadratically near a root and can diverge or cycle far from one. Damping shortens the Newton step so that each iterate stays where the linearization is trustworthy. Simulators use several forms: clipping each node update, scaling the whole update vector, and backtracking until the residual decreases. All three keep the Newton direction but trade the step length, and with it some speed, for robustness.

## Clipping the step

The simplest damping clips each node's update to a maximum magnitude $\Delta_{max}$: $v_{k+1} = v_k + \operatorname{sign}(\Delta v)\min(|\Delta v|, \Delta_{max})$. If Newton at $v_k = 0.2$ V proposes $\Delta v = -0.9$ V with $\Delta_{max} = 0.3$ V, the next iterate is -0.1 V. If the proposal is smaller than the limit, it is applied in full, so clipping does not disturb the final quadratic phase. Clipping per node is cheap and common in device-oriented limiting, but because it changes the components of the step unequally, it can bend the update away from the Newton direction.

## Uniform damping of a vector update

To keep the direction, scale the whole update by one factor: $x_{k+1} = x_k + \alpha \Delta x$ with $\alpha = \min(1, \Delta_{max}/\|\Delta x\|_\infty)$. The largest component moves exactly $\Delta_{max}$, and the others move proportionally. For an update $(0.4, -1.2, -0.2)$ V with $\Delta_{max} = 0.5$ V, $\|\Delta x\|_\infty = 1.2$ and $\alpha = 0.4167$, so the first node moves by 0.167 V instead of 0.4 V. Because the direction is still the Newton direction, which is a descent direction for the residual norm, a small enough $\alpha$ always reduces the residual.

## Backtracking on the residual

A more principled scheme picks $\alpha$ by testing the true residual. Try $\alpha = 1, 1/2, 1/4, \ldots$ and accept the first that gives sufficient decrease, for example $\|f(x + \alpha \Delta x)\| \le (1 - \alpha/2)\|f(x)\|$. With $\|f(x)\| = 5$ mA, trial norms of 6.1 at $\alpha = 1$ (needs at most 2.5) and 3.2 at $\alpha = 1/2$ (needs at most 3.75) accept $\alpha = 0.5$. Each trial costs a device evaluation but no new factorization. Backtracking guarantees progress on the residual, but in circuits with very different node scales the residual norm needs weighting, or one large current dominates the test.

## Going deeper

Dennis and Schnabel, "Numerical Methods for Unconstrained Optimization and Nonlinear Equations" (SIAM, 1996), chapter 6; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
