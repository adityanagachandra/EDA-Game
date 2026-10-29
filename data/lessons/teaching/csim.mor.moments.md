# Moments of an RC network

Model order reduction for interconnect starts from moments: the coefficients of the Taylor expansion of a transfer function around $s = 0$. For RC networks the first moment is minus the Elmore delay, and higher moments follow from a simple recursion that reuses one factorization. Moments are cheap, they capture the low-frequency behavior that dominates delay, and they are what AWE and Krylov methods match.

## The first moment at the far end

Expand the transfer function from the driven root to node $i$ as $H_i(s) = m_0 + m_1 s + m_2 s^2 + \cdots$. For an RC tree driven by a voltage source, $m_0 = 1$ (at DC every node follows the source) and
$$m_1^{(i)} = -\sum_k R_{ik} C_k,$$
where $R_{ik}$ is the resistance shared by the paths from the root to $i$ and to $k$. This is minus the Elmore delay. For a chain of $R_1 = 100\ \Omega$ to node 1 with 20 fF and $R_2 = 200\ \Omega$ to node 2 with 10 fF, the far end has $m_1 = -(100 \cdot 20 + 300 \cdot 10)$ ohm fF $= -5000$ fs $= -5$ ps.

## The first moment at an inner node

At an inner node of a chain, capacitors beyond the node share only the resistance up to that node. At node 1 of the chain above, $C_1$ shares $R_1$ and so does $C_2$, so $m_1^{(1)} = -100(20 + 10) = -3000$ fs $= -3$ ps. The general computation is a linear solve: moments satisfy $G m_0 = b$ and $G m_q = -C m_{q-1}$, so $m_1 = -G^{-1} C m_0$. For a tree, $G^{-1}$ has the shared-resistance entries $R_{ik}$, and the recursion becomes two tree traversals: a post-order pass accumulating downstream capacitance and a pre-order pass accumulating resistance times downstream capacitance.

## The second moment

Higher moments reuse the same recursion with the previous moment as the excitation: $m_q = -G^{-1} C m_{q-1}$, so in a tree
$$m_2^{(i)} = -\sum_k R_{ik} C_k m_1^{(k)}.$$
Each capacitor is weighted by the first moment at its own node. In the chain above, $m_1^{(1)} = -3$ ps and $m_1^{(2)} = -5$ ps, so at node 2, $m_2 = 100 \cdot 20 \cdot 3 + 300 \cdot 10 \cdot 5 = 6000 + 15000 = 21000$ fs ps $= 21$ ps$^2$. Moments alternate in sign for RC trees. Each new moment costs one forward and back substitution with the factored $G$, which is why moment computation scales to huge parasitic networks.

## Going deeper

Pillage and Rohrer, "Asymptotic waveform evaluation for timing analysis" (IEEE TCAD, 1990); Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods" (McGraw-Hill, 1995).
