# Direct sensitivity analysis

Sensitivity analysis asks how much an output moves when a parameter moves. Simulators compute it exactly, without finite differences, by differentiating the circuit equations themselves and reusing the factorization of the nominal solve. Sensitivities feed optimization and sizing, statistical design and worst-case estimation, so they are one of the most reused quantities a simulator produces.

## Differentiating the circuit equations

Differentiate $G(p)\, x(p) = b$ with respect to a parameter $p$: $G\, \partial x/\partial p = -(\partial G/\partial p)\, x$. The matrix is the nominal $G$, already factored, so each parameter costs one forward and back substitution. For a divider with $V_1$ through $R_1$ into the output and $R_2$ to ground, $v_{out} = V_1 R_2/(R_1 + R_2)$, so
$$\frac{\partial v_{out}}{\partial R_1} = -\frac{V_1 R_2}{(R_1 + R_2)^2}, \qquad \frac{\partial v_{out}}{\partial R_2} = \frac{V_1 R_1}{(R_1 + R_2)^2}.$$
With $V_1 = 2$ V, $R_1 = 3$ k$\Omega$ and $R_2 = 1$ k$\Omega$, $\partial v_{out}/\partial R_1 = -2/16 = -0.125$ V/k$\Omega$, or -125 mV per k$\Omega$. This direct method is exact to rounding.

## Normalized sensitivity

Raw derivatives have units and depend on scale. The normalized sensitivity
$$S^{v}_{p} = \frac{p}{v}\frac{\partial v}{\partial p}$$
is the percent change in the output per percent change in the parameter. For the divider, $S^{v}_{R_1} = -R_1/(R_1 + R_2)$ and $S^{v}_{R_2} = R_1/(R_1 + R_2)$; with $R_1 = 3$ k$\Omega$ and $R_2 = 1$ k$\Omega$ these are -0.75 and 0.75. They sum to zero, as they must: scaling both resistors together leaves a ratio unchanged. Normalized sensitivities let a designer rank parameters by influence regardless of units, which is the first step in deciding which tolerances matter and which device to resize.

## First-order change estimate

A sensitivity predicts small changes linearly: $\Delta v \approx (\partial v/\partial p)\Delta p$. For the divider above, a 2% increase in $R_1$ is $\Delta R_1 = 0.06$ k$\Omega$, so $\Delta v_{out} \approx -0.125 \times 0.06 = -7.5$ mV. Equivalently, the normalized sensitivity times the percent change gives the percent output change: $-0.75 \times 2\% = -1.5\%$ of 0.5 V. The estimate degrades for large changes or strongly nonlinear circuits, where the second derivative matters. Worst-case and statistical methods combine many first-order terms, so their accuracy rests on the parameter changes being small enough for linearity to hold.

## Going deeper

Director and Rohrer, "The generalized adjoint network and network sensitivities" (IEEE Transactions on Circuit Theory, 1969); Najm, "Circuit Simulation" (Wiley, 2010), chapter on sensitivity.
