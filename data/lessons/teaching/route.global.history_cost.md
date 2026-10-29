# History costs in negotiated congestion

Present-congestion cost alone can make a router oscillate: two nets swap the same contested edge back and forth forever. History cost breaks the cycle. Each edge accumulates a penalty for every iteration it ends overused, and that memory makes chronically congested edges permanently expensive, so nets learn to avoid them even when they happen to fit for a moment.

## Accumulating history

After each iteration, every resource's history grows by its overuse: $h_n \leftarrow h_n + h_{fac}\max(0, occ_n - cap_n)$, starting from 0. History never decreases. Example: an edge of capacity 5 ends four iterations with occupancies 6, 4, 8 and 7; with $h_{fac} = 0.5$ it gains 0.5, 0, 1.5 and 1, so its history is 3.

Note the asymmetry with the present factor, which uses $occ_n + 1 - cap_n$ because it prices adding one more net now. History records what actually happened at the end of each iteration, so it uses $occ_n - cap_n$. An edge that was overused in iteration 2 stays more expensive in iteration 10 even if nobody overuses it in between.

## History moves stubborn nets

Consider a net whose short path A crosses one hot edge and whose alternative B is longer but free. Ignoring present cost, A costs $b_A + h$ and B costs $b_B$. If the hot edge ends every iteration overused by $o$, its history after $i - 1$ iterations is $(i-1) h_{fac} o$, and A first becomes more expensive in the smallest iteration $i$ with $b_A + (i-1) h_{fac} o > b_B$.

Example: $b_A = 9$, $b_B = 14$, $o = 1$ and $h_{fac} = 1$: A loses in iteration 7. Without history, if present cost alone left A slightly cheaper for this net every time it was rerouted, the net would never move and the edge would stay overused. History guarantees that persistent congestion eventually outweighs any finite path-length advantage.

## History and present cost together

In PathFinder the two terms combine as $(b_n + h_n)\,p_n$. History raises the floor of an edge's cost; the present factor multiplies it when the edge is full right now. Example: an edge with $b = 1$, capacity 4 and end-of-iteration occupancies 3 and 4 has history 0 (it was never overused). In iteration 3 with 5 nets already on it and $p_{fac} = 1$, its present factor is $1 + 2 = 3$ and its cost is 3.

The two factors must be tuned together. A large $h_{fac}$ converges quickly but can lock nets out of regions they need; a small one lets oscillation persist. Global routers such as NCTU-GR and NTHU-Route use history-based costs of this kind, often with variations such as letting history decay or scaling it by an edge's capacity.

## Going deeper

L. McMurchie and C. Ebeling, "PathFinder" (1995); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
