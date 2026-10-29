# k-feasible cuts

Cut-based mapping replaces pattern graphs with cuts. A cut of a node is a set of nodes that separates it from the primary inputs; if it has at most k members, the logic between the cut and the node is a function of at most k inputs, which a k-LUT, or a library cell found by Boolean matching, can implement. Cuts let matches cross fanout points freely, so the mapper works on the DAG directly.

## Definition of a cut

In an AIG, a **cut** C of node n is a set of nodes such that every path from a primary input to n passes through a node of C. The **trivial cut** {n} always qualifies. A cut is **k-feasible** if |C| <= k. The nodes between C and n form the cut's cone, and n is a Boolean function of the cut's leaves. For n3 = AND(n1, n2) with n1 = AND(a, b) and n2 = AND(b, c), the sets {n3}, {n1, n2}, {n1, b, c}, {a, b, n2} and {a, b, c} are all cuts; {a, b} is not, because the path from c to n3 avoids it.

## Counting cuts by merging

Cuts are enumerated bottom-up. A primary input has only its trivial cut. For an AND node n with fanins u and w, every union of one cut of u with one cut of w is a cut of n; keep the unions with at most k leaves, add the trivial cut {n}, and discard any cut that strictly contains another cut of n, since its leaves add nothing. In the example with k = 3: cuts(n1) = {n1}, {a, b}; cuts(n2) = {n2}, {b, c}; the four unions are {n1, n2}, {n1, b, c}, {a, b, n2}, {a, b, c}, all feasible; with {n3} that is 5 cuts.

## How the count grows

The number of cuts grows quickly with k and with reconvergence. A node whose fanins have $c_1$ and $c_2$ cuts can have up to $c_1 c_2 + 1$ cuts, and over a deep network the total is exponential in k in the worst case; for k = 6 an AIG node can have hundreds of cuts. Reconvergence also creates dominated unions: if the same input reaches n by two paths, many unions coincide. Raising k from 3 to 4 adds every cut with exactly four leaves that was previously filtered. This growth is why practical mappers keep only a few priority cuts per node instead of all of them.

## Going deeper

Cong, Wu and Ding, "Cut ranking and pruning: enabling a general and efficient FPGA mapping solution", FPGA 1999; Pan and Lin, "A new retiming-based technology mapping algorithm for LUT-based FPGAs", FPGA 1998, for cut enumeration.
