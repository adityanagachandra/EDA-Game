# Adjoint sensitivity

When a circuit has many parameters and few outputs, the direct method's one solve per parameter becomes expensive. The adjoint method reverses the bookkeeping: one solve per output, with the transposed matrix, gives the sensitivity of that output to every parameter at once. It is the circuit-simulation ancestor of reverse-mode automatic differentiation and backpropagation.

## The adjoint solve

For an output $y = e^T x$, where $e$ selects the output node, solve the adjoint system
$$G^T \lambda = e.$$
The vector $\lambda$ is the response of every node to a unit current injected at the output in the transposed circuit. For a resistive network $G$ is symmetric, so $\lambda = G^{-1} e$. With $g_a$ from node 1 to ground, $g_b$ between nodes 1 and 2, $g_c$ from node 2 to ground and output $v_2$: $\lambda_1 = g_b/\det G$ and $\lambda_2 = (g_a + g_b)/\det G$. For $g_a = g_b = g_c = 1$ mS, $\det G = 3$, $\lambda_1 = 0.333$ k$\Omega$ and $\lambda_2 = 0.667$ k$\Omega$. The factors of $G$ are reused: solving with $G^T$ swaps the roles of $L$ and $U$.

## Sensitivity from two solutions

From the direct formula, $\partial y/\partial p = -e^T G^{-1}(\partial G/\partial p) x = -\lambda^T (\partial G/\partial p)\, x$. For a conductance $g$ between nodes $p$ and $q$, $\partial G/\partial g$ is the resistor stamp, so
$$\frac{\partial y}{\partial g} = -(\lambda_p - \lambda_q)(v_p - v_q),$$
with ground at 0. In the example with 1 mA into node 1, $v_1 = 0.667$ V and $v_2 = 0.333$ V. Then $\partial v_2/\partial g_b = -(0.333 - 0.667)(0.667 - 0.333) = 0.111$ V/mS: more conductance between the nodes raises the output. Every element's sensitivity is an inner product of the nominal and adjoint solutions at its terminals.

## Counting solves: direct versus adjoint

After the nominal factorization, the direct method needs one substitution per parameter and yields every output; the adjoint method needs one substitution per output and yields every parameter. With 300 parameters and 2 outputs, direct costs 300 solves and adjoint costs 2, each followed by cheap inner products. Optimization of a few performance metrics over many device sizes is therefore an adjoint problem. The same reasoning appears in noise analysis, which uses one adjoint solve per frequency to get every noise source's transfer to the output, and in transient adjoint sensitivity, which integrates backward in time and must store the forward trajectory.

## Going deeper

Director and Rohrer, "The generalized adjoint network and network sensitivities" (IEEE Transactions on Circuit Theory, 1969); Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994).
