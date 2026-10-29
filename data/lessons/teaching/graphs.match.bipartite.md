# Bipartite matching

Many EDA assignment problems pair items of one kind with slots of another under compatibility rules: pins to package bumps, I/O signals to pads, cells to legal sites, registers to scan slots, nets to tracks. Each is a bipartite matching. Maximum matching has fast exact algorithms, and its dual, the minimum vertex cover, certifies that no better assignment exists, which a tool can report when an assignment fails.

## Matchings and augmenting paths

A bipartite graph has left vertices (pins), right vertices (slots) and edges for allowed pairs. A **matching** is a set of edges with no shared endpoint; a maximum matching has the most edges. A matching is maximum exactly when it has no **augmenting path**, a path that alternates between unmatched and matched edges and starts and ends at unmatched vertices. Flipping such a path adds one edge.

With pins P1 to S2 only, P2 to S1, S2, S3, P3 to S1, S2, S4, P4 to S5 and P5 to S2 only, both P1 and P5 need S2, so at most one of them is placed, and the maximum is 4. Equivalently, add a source joined to every pin and a sink joined from every slot, all capacities 1: the max flow equals the maximum matching.

## Greedy versus maximum

A greedy pass that gives each pin, in order, its first free allowed slot is fast but can block later pins. If P1 takes S1 and P4, which can use only S1, comes later, P4 is stranded although P1 could have used S2. An augmenting path P4 to S1, then S1's partner P1 to S2, repairs this by shifting P1.

So the difference between the greedy size and the maximum is the number of augmenting steps still needed. Hopcroft and Karp's algorithm finds many shortest augmenting paths per phase and runs in $O(E\sqrt{V})$, fast enough for millions of sites. In legalization and detailed placement, matching-based moves reassign groups of cells to sites optimally, which greedy placement cannot do.

## Covers and Hall's condition

König's theorem says that in a bipartite graph the size of a maximum matching equals the size of a minimum **vertex cover**, a smallest set of vertices touching every edge. If 5 pins compete for 4 slots and a cover of size 4 exists, no assignment can place more than 4 pins: the cover is a certificate.

Hall's theorem gives the failure condition directly: all left vertices can be matched exactly when every subset $X$ of them has at least $|X|$ neighbours. A pin assigner that fails can therefore report a concrete obstruction, such as three pins whose combined allowed slots number only two, which tells the designer exactly which constraints to relax.

## Going deeper

Hopcroft and Karp, "An n^(5/2) algorithm for maximum matchings in bipartite graphs" (SIAM J. Computing, 1973); Lovász and Plummer, "Matching Theory" (1986).
