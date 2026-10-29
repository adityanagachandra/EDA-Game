# Reverse mode and adjoint sensitivity

When a single scalar, a delay, an IR drop or a placement cost, depends on many parameters, its gradient is best computed backwards. Reverse-mode automatic differentiation and the adjoint method of circuit sensitivity are the same idea: one extra pass, or one extra transposed solve, gives the derivative with respect to every parameter at once.

## Sweeping adjoints backwards through the graph

Reverse mode first evaluates the function and records the computational graph (the tape). It then seeds the output adjoint $\bar f = 1$ and walks the graph backwards: each node $v = g(u_1, u_2)$ adds $\bar v\,\partial g/\partial u_i$ to $\bar u_i$. Contributions to a variable used in several places add up.

For $f = (x + y)\cdot xy$ at $(2, 3)$: $a = x + y = 5$, $b = xy = 6$, $f = ab = 30$. Backwards: $\bar a = b = 6$, $\bar b = a = 5$. Then $\bar x = \bar a \cdot 1 + \bar b \cdot y = 6 + 15 = 21$ and $\bar y = \bar a + \bar b \cdot x = 6 + 10 = 16$. Check: $f = x^2y + xy^2$, $\partial f/\partial x = 2xy + y^2 = 21$. Backpropagation in neural networks is exactly this.

## Adjoint sensitivity of a linear circuit

A linear circuit solves $Gx = b$, and an output is $y = c^Tx$, say a node voltage. Then $\partial y/\partial b = G^{-T}c$, so solve once
$$G^T\lambda = c,$$
and read $\partial y/\partial b_j = \lambda_j$ for every injection $j$. For a conductance $p$ that enters $G$, $\partial y/\partial p = -\lambda^T(\partial G/\partial p)x$; for a single entry, $\partial y/\partial G_{ij} = -\lambda_ix_j$.

Example: $G = \begin{pmatrix} 3 & -1 \\ -1 & 2 \end{pmatrix}$ mS, output $y = x_2$, so $c = e_2$. Solving $G^T\lambda = e_2$: $\det = 5$, $\lambda = (1/5, 3/5) = (0.2, 0.6)$ kΩ. A 1 mA injection at node 1 raises $x_2$ by 0.2 V. The factorization of $G$ is reused for the transposed solve, so all sensitivities cost about one extra solve.

## When reverse mode pays off

Reverse mode computes the gradient of one output in a small constant multiple of one function evaluation, often quoted as at most 3 to 5, whatever the number of inputs, and a full Jacobian in $m$ passes for $m$ outputs. Forward mode needs $n$ passes. With $n = 10^6$ cell coordinates and $m = 1$ objective, forward mode is about $10^6$ evaluations and reverse about 3.

The costs are memory and complexity. The tape stores intermediate values of the whole forward pass, so long time integrations need checkpointing, recomputing segments instead of storing everything. DREAMPlace casts analytical placement as a differentiable computation and uses reverse mode through GPU deep-learning frameworks; differentiable timing and transient adjoint sensitivity in simulators apply the same principle. Forward mode remains better for few-input functions such as device models.

## Going deeper

Griewank and Walther, "Evaluating Derivatives" (2nd ed., SIAM, 2008), chapters 3 and 4; Director and Rohrer, "The generalized adjoint network and network sensitivities", IEEE Trans. Circuit Theory (1969).
