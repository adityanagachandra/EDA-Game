# Van Ginneken candidate operations

Van Ginneken's buffer insertion algorithm represents every partial solution at a point of a net by two numbers: the load seen looking downstream and the required time there. Three operations transform these candidates as the algorithm walks from the sinks to the driver: through a wire, through a buffer, and across a branch. Each is a short formula from the Elmore delay.

## Through a wire segment

A candidate (C, q) at the downstream end of a wire of resistance R and capacitance $C_w$ becomes, at the upstream end,
$$C' = C + C_w, \qquad q' = q - R\left(\frac{C_w}{2} + C\right).$$
The wire is modelled as a pi segment: half its capacitance is charged through its own resistance, plus the full downstream load. Example: (12 fF, 580 ps) through 0.3 kOhm and 50 fF becomes (62 fF, 580 - 0.3 x 37 = 568.9 ps). Required time falls and load grows as the walk moves upstream.

## Through a buffer

Inserting a buffer with input capacitance $C_b$, drive resistance $R_b$ and intrinsic delay $d_b$ at a site with candidate (C, q) gives
$$C' = C_b, \qquad q' = q - d_b - R_b C.$$
The buffer hides the downstream load: upstream sees only its small input capacitance. Example: (91 fF, 395 ps) with $C_b$ = 6 fF, $R_b$ = 0.6 kOhm, $d_b$ = 16 ps becomes (6 fF, 395 - 16 - 54.6 = 324.4 ps). The required time drops, but so does the load the upstream wire must charge, which pays off when the upstream resistance is large. At each site the algorithm keeps both the unbuffered and the buffered candidate.

## Across a branch

At a Steiner point where two subtrees meet, every left candidate combines with every right candidate: $C = C_1 + C_2$ and $q = \min(q_1, q_2)$, since both loads hang on the same node and the tighter requirement governs. Merging lists of sizes 2 and 3 forms 6 pairs, many dominated. If both lists are sorted by load with required time increasing, a linear merge, advancing on the side with the smaller q, produces all non-dominated combinations, at most $|L_1| + |L_2| - 1$ of them. This keeps the algorithm polynomial: for one buffer type and n sites, the candidate lists stay O(n) long and the total time is $O(n^2)$.

## Going deeper

van Ginneken, "Buffer placement in distributed RC-tree networks for minimal Elmore delay", ISCAS 1990; Alpert, Mehta and Sapatnekar (eds.), "Handbook of Algorithms for Physical Design Automation" (2008), chapter on buffer insertion.
