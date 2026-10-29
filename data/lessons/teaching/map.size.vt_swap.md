# Vt swapping under a slack budget

Most cell libraries come in several threshold-voltage flavours with the same footprint. High-Vt cells leak far less but switch more slowly. After timing closes, a leakage-recovery pass swaps cells on non-critical paths to higher Vt, spending slack for leakage. The decision is a knapsack problem, and the exponential leakage law makes it worth solving well.

## The leakage law

Subthreshold leakage scales as $\exp(-V_t / (n v_T))$, with subthreshold slope factor n and thermal voltage $v_T \approx$ 26 mV at room temperature. Raising $V_t$ by $\Delta V_t$ multiplies leakage by $\exp(-\Delta V_t / (n v_T))$. With $\Delta V_t$ = 100 mV and n = 1.4 the factor is $\exp(-100/36.4) = 0.064$, so a 136 nW cell drops to about 8.7 nW. The delay increase is far smaller in relative terms, typically a modest percentage of the cell delay, which is what makes swapping attractive.

## One path: a knapsack

Gates on one path, and on no other critical path, each offer a swap with delay cost $\Delta d_i$ and leakage saving $\Delta L_i$. Choose a set with total added delay at most the path's slack S and maximum total saving. This is the 0-1 knapsack problem: NP-hard in general, solvable by dynamic programming over integer delay budgets in time proportional to n S, and trivially by enumeration for a handful of gates. Greedy by saving per picosecond is a common heuristic but can miss the optimum when one large-saving gate does not fit with the others.

## Shared gates and many paths

A gate usually sits on many paths, and its added delay counts against all of them. With two paths, the selection must keep both slacks non-negative: a knapsack with two capacity constraints. In a real netlist there is one constraint per timing path, so tools work with slack at each node from static timing, swap greedily in order of leakage saved per unit of slack consumed, and re-time incrementally after each batch, undoing swaps that create violations. Swaps also change slews and loads seen by neighbours, so the incremental timer, not a precomputed delay, must judge each move.

## Going deeper

Wei, Chen, Roy et al., "Design and optimization of dual-threshold circuits for low-voltage low-power applications", IEEE TVLSI 1999; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter on leakage.
