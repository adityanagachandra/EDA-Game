# Node elimination by star-mesh transform

The workhorse of local RC reduction is eliminating one internal node at a time. Removing a node is Gaussian elimination on the network's nodal equations, which for a resistive star becomes the star-mesh transform: new resistors between every pair of former neighbours. The eliminated node's capacitance is spread onto those neighbours. TICER-style reducers apply this to the nodes that settle fastest, and their cost and fill-in behave exactly like sparse matrix elimination.

## The star-mesh resistors

A node $x$ connected to neighbours by conductances $g_1, \dots, g_k$, with no other path, can be eliminated exactly at DC. Solving its KCL equation for its voltage and substituting gives new conductances $g_{ij} = g_i g_j / \sum_m g_m$ between every pair of neighbours, or in resistances $R_{ij} = R_i R_j \sum_m (1/R_m)$. For three resistors of 70, 80 and 180 ohm, $\sum 1/R = 0.03234$ S, so the new resistor between the 70 and 180 ohm neighbours is $70 \times 180 \times 0.03234 = 407.5$ ohm. With two neighbours the formula reduces to the series sum $R_1 + R_2$. The transform is the Schur complement of the node's row and column in the conductance matrix.

## Redistributing the capacitance

The eliminated node's capacitance cannot simply vanish. TICER moves it to the neighbours in proportion to their conductances: neighbour $i$ receives $C_x g_i/\sum_m g_m$. With 7 fF at a node connected through 90, 140 and 80 ohm, the 90 ohm neighbour's share is $(1/90)/0.03075 = 0.361$, so it gains 2.53 fF. The shares sum to one, so total capacitance is preserved exactly, and, as the chain case shows, this split also preserves the Elmore delay at downstream nodes. It is a first-order approximation for nodes with significant capacitance; that is why only nodes with small time constants are eliminated.

## Which nodes to eliminate

A node's time constant is $\tau = C_x / \sum_m g_m$, its capacitance times the parallel resistance of its connections. A node with $\tau$ much smaller than the signal transitions tracks its neighbours almost instantly, so eliminating it barely changes waveforms. TICER marks nodes below a threshold and eliminates them; nodes with large $\tau$ shape the response and stay. For 10 fF with 1900 and 900 ohm in parallel, $\tau = 10 \times 611 = 6.1$ ps; for 5.5 fF with 100 and 1100 ohm, $\tau = 0.5$ ps. Order matters for cost: eliminating a node of degree $k$ creates up to $k(k-1)/2$ resistors, so reducers prefer low-degree nodes, exactly the minimum-degree heuristic of sparse LU, to avoid fill-in that would make the reduced network denser than the original.

## Going deeper

Sheehan, "TICER: realizable reduction of extracted RC circuits", ICCAD (1999); Davis, "Direct Methods for Sparse Linear Systems" (2006), on elimination orderings.
