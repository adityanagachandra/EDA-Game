# Solving a small MNA system

After stamping comes solving. A simulator factors the sparse MNA matrix with LU and substitutes; by hand, a two-node system yields to Cramer's rule. Solving small systems exactly builds the intuition needed to sanity-check simulator output and to derive the companion-model steps, sensitivities and moments that later concepts compute.

## Two-node nodal equations

A current source $I$ into node 1, conductance $g_1$ from node 1 to ground, $g_2$ between nodes 1 and 2, and $g_3$ from node 2 to ground give
$$\begin{pmatrix} g_1 + g_2 & -g_2 \\ -g_2 & g_2 + g_3 \end{pmatrix} \begin{pmatrix} v_1 \\ v_2 \end{pmatrix} = \begin{pmatrix} I \\ 0 \end{pmatrix}.$$
With $\det G = (g_1 + g_2)(g_2 + g_3) - g_2^2$, Cramer's rule gives $v_1 = I(g_2 + g_3)/\det G$. For $I = 1$ mA, $g_1 = 1$, $g_2 = 1$, $g_3 = 2$ mS: $\det G = 2 \cdot 3 - 1 = 5$, so $v_1 = 3/5 = 0.6$ V.

## Solving for the second node

The second unknown follows the same way: replace the second column by $b$, giving $v_2 = I g_2 / \det G$. In the example, $v_2 = 1/5 = 0.2$ V. A quick check is KCL at node 2: the current through $g_2$, $g_2(v_1 - v_2) = 0.4$ mA, must equal the current to ground, $g_3 v_2 = 0.4$ mA. In a simulator the same solve is done by LU: factor $G = LU$, then forward-substitute $Ly = b$ and back-substitute $Ux = y$. The factorization is reused when only $b$ changes, which is what makes sensitivity and multiple-source analyses cheap.

## A system with a voltage source row

When V1 holds node 1 at $V_s$, the unknowns are $v_1$, $v_2$ and the branch current. The branch row is simply $v_1 = V_s$, so it can be substituted: the node-2 row becomes $-g_2 V_s + (g_2 + g_3) v_2 = I_2$, where $I_2$ is any current injected at node 2. Then $v_2 = (g_2 V_s + I_2)/(g_2 + g_3)$. With $V_s = 3.5$ V, $g_2 = g_3 = 0.25$ mS and $I_2 = -0.5$ mA, $v_2 = (0.875 - 0.5)/0.5 = 0.75$ V. The node-1 row then gives the branch current. A sparse LU does this elimination implicitly, but it must pivot around the zero diagonal of the branch row.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 3 on linear equation solution; Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006).
