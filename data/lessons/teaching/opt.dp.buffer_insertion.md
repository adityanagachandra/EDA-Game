# Van Ginneken buffer insertion

Van Ginneken's algorithm inserts buffers on a routing tree to maximize slack at the driver. It is a bottom-up dynamic program over candidate solutions, each summarized by the load it presents upstream and the required time it guarantees. Pruning dominated candidates keeps the lists short, and the algorithm is optimal for the Elmore delay model and a given set of buffer positions. It is the basis of buffering in physical synthesis, with extensions for multiple buffer types, slew and cost.

## Propagating a candidate through a wire

A candidate is a pair $(C, q)$: the downstream capacitance it presents and the latest required time at its root. Moving up through a wire with resistance $R_w$ and capacitance $C_w$, modelled as a pi segment, the wire's Elmore delay is $R_w(C_w/2 + C)$, so
$$q' = q - R_w\left(\tfrac{C_w}{2} + C\right), \qquad C' = C + C_w.$$
Example: $(C, q) = (5\text{ fF}, 113\text{ ps})$ through $R_w = 179\ \Omega$ and $C_w = 29$ fF. The delay is $179 \cdot 19.5 = 3490\ \Omega\cdot$fF $= 3.49$ ps, so $q' = 109.51$ ps and $C' = 34$ fF. Units matter: ohms times femtofarads are femtoseconds, so divide by 1000 for picoseconds. At a branch point, the two children's candidates merge: loads add and the required time is the smaller of the two.

## Inserting a buffer

At a legal buffer position every candidate can be buffered. A buffer with intrinsic delay $D_b$, output resistance $R_b$ and input capacitance $C_b$ drives the downstream load C and isolates it, so
$$(C, q) \mapsto \bigl(C_b,\ q - D_b - R_b C\bigr).$$
Example: candidates (58, 476), (46, 456) and (18, 276) in fF and ps, and a buffer with $D_b = 21$ ps, $R_b = 176\ \Omega$, $C_b = 6$ fF. The buffered required times are $476 - 21 - 10.21 = 444.79$, $456 - 21 - 8.10 = 426.90$ and $276 - 21 - 3.17 = 251.83$. All three buffered candidates have the same load, 6 fF, so only the best, 444.79 ps, survives. The unbuffered candidates also stay in the list, so each position adds at most one new candidate per buffer type.

## Pruning and choosing at the driver

Candidate $(C_1, q_1)$ dominates $(C_2, q_2)$ when $C_1 \le C_2$ and $q_1 \ge q_2$: it is no harder to drive and no worse in timing, so the other can be discarded. Sorting by load, a candidate survives only if its required time beats every lighter one. From (10, 120), (20, 150), (25, 140), (40, 210) and (50, 200), the survivors are (10, 120), (20, 150) and (40, 210). Pruning keeps the list linear in the number of buffer positions, giving an $O(n^2)$ algorithm, later improved to $O(n \log n)$ by Shi and Li.

At the source the driver with resistance $R_d$ adds $R_d C$, and the slack is $\max(q - R_d C)$. With candidates (9, 318) and (108, 305) and $R_d = 657\ \Omega$: $318 - 5.91 = 312.09$ and $305 - 70.96 = 234.04$, so the answer is 312.09 ps, and a traceback recovers the buffer positions behind it.

## Going deeper

van Ginneken, "Buffer placement in distributed RC-tree networks for minimal Elmore delay" (ISCAS, 1990); Lillis, Cheng and Lin, "Optimal wire sizing and buffer insertion for low power and a generalized delay model" (IEEE JSSC, 1996); Shi and Li, "A fast algorithm for optimal buffer insertion" (IEEE TCAD, 2005).
