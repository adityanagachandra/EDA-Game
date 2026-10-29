# Subgradient updates of multipliers

The dual function of a relaxed problem is concave but often not differentiable, especially when the primal has discrete choices. Subgradient ascent updates the multipliers using the constraint violations as a subgradient, projects them back to nonnegative values, and shrinks the step over time. It is simple, scales to millions of multipliers, and is the classic outer loop of Lagrangian-relaxation sizing, though it converges slowly and needs care with step sizes.

## A multiplier update from the violation

For a constraint $a \le T$ with multiplier $\lambda$, the value $a - T$ at the current relaxed solution is a subgradient of the dual function with respect to $\lambda$. Ascent moves the multiplier along it:
$$\lambda \leftarrow \max\bigl(0,\ \lambda + t(a - T)\bigr).$$
Example: $\lambda = 0.8$, step $t = 0.008$ per ps, required time 710 ps, arrival 850 ps. The violation is 140 ps, so $\lambda \leftarrow 0.8 + 1.12 = 1.92$.

The intuition is a price adjustment. A violated constraint gets more expensive, so the next relaxed solve spends more area to speed it up. A constraint with slack gets cheaper, so resources move away from it. Unlike gradient descent there is no guarantee that the dual value increases each step, so implementations track the best bound seen so far.

## Projection onto nonnegative multipliers

Multipliers of inequality constraints must stay nonnegative, so the raw update is projected: any negative result becomes zero. Example: $\lambda = 1.2$, $t = 0.006$, required time 900 ps and arrival 600 ps. The constraint has 300 ps of slack, so the raw update is $1.2 + 0.006 \cdot (-300) = -0.6$, and the projection gives $\lambda = 0$.

A zero multiplier means the constraint is currently ignored by the relaxed problem: its path gets no weight, and its gates drift toward minimum size until the path becomes critical again and its multiplier returns. That oscillation is a known behaviour of subgradient methods. In LR sizing the projection is onto the flow-conservation polytope rather than onto the nonnegative orthant alone, for example by rescaling the multipliers entering each node so that flow balances, working back from the outputs.

## Step-size rules and convergence

With a constant step, subgradient ascent only reaches a neighbourhood of the dual optimum. Convergence needs diminishing steps whose sum diverges, for example $t_k = t_0/k$. Two iterations: $\lambda_0 = 1.5$, $t_0 = 0.004$, violations 200 ps then -120 ps. Then $\lambda_1 = \max(0, 1.5 + 0.004 \cdot 200) = 2.3$ and $\lambda_2 = \max(0, 2.3 + 0.002 \cdot (-120)) = 2.06$.

The Polyak step $t_k = (q^* - q(\lambda_k))/\|s_k\|^2$ converges faster but needs an estimate of the optimal dual value $q^*$, often replaced by the best feasible primal cost. Convergence is typically sublinear, which is why LR sizers run a fixed budget of iterations, track the best feasible sizing and the best bound, and stop when the gap is small. Bundle and cutting-plane methods stabilize the steps when more accuracy is needed.

## Going deeper

Bertsekas, "Nonlinear Programming" (3rd ed., 2016), chapter on dual methods and subgradient optimization; Held, Wolfe and Crowder, "Validation of subgradient optimization" (Mathematical Programming, 1974); Chen, Chu and Wong (IEEE TCAD, 1999) for the sizing application.
