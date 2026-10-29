# Depth-optimal LUT mapping by cuts

Depth-optimal LUT mapping assigns every node of a logic network the fewest LUT levels any mapping can give it. The modern formulation enumerates $k$-feasible cuts and labels nodes in topological order; FlowMap reaches the same labels with a max-flow test. This is the core of the ABC and VPR-flow LUT mappers, and it is one of the few mapping problems solvable optimally in polynomial time.

## Enumerating k-feasible cuts

A cut of node $v$ is a set of nodes such that every path from a primary input to $v$ passes through the set; it is $k$-feasible if it has at most $k$ leaves, so the cone it bounds fits one $k$-LUT. For a two-input gate $v = f(x, y)$ the cuts are the trivial cut $\{v\}$ plus every union $C_x \cup C_y$ of a cut of $x$ and a cut of $y$ with at most $k$ leaves. A primary input has only its trivial cut. Take $g_1 = \mathrm{AND}(a, b)$ and $g_2 = \mathrm{OR}(g_1, c)$ with $k = 3$. The cuts of $g_1$ are $\{g_1\}$ and $\{a, b\}$, so $g_2$ has $\{g_2\}$, $\{g_1, c\}$ and $\{a, b, c\}$: three cuts. Visiting nodes in topological order computes all cuts bottom-up.

## Depth-optimal labels

Give primary inputs label 0. For each gate in topological order, its label is the minimum over its non-trivial $k$-feasible cuts $C$ of $1 + \max_{u \in C} \ell(u)$. The label is the depth of $v$ in the best possible mapping, because the LUT rooted at $v$ implements the cone of some cut, and its inputs must be produced first. For a chain $g_1 = f(a, b)$, $g_2 = f(g_1, c)$, $g_3 = f(g_2, d)$ with $k = 3$: $g_1$ and $g_2$ get label 1 (cut $\{a, b, c\}$ for $g_2$), and $g_3$ has cuts $\{g_2, d\}$ and $\{g_1, c, d\}$ with best value $1 + 1 = 2$. FlowMap proves that this minimum is found by a max-flow test asking whether a cut exists whose leaves all have label below some value, which avoids enumerating every cut.

## Covering and what depth saves

Labels fix the depth; the cover then picks, from each primary output backwards, the cut that achieves the label, and recursively covers that cut's leaves. Compared with a one-gate-per-LUT mapping, whose depth is the gate level, the saving is the gate level minus the label. A node 6 gates deep with $k = 3$ and label 3 saves three LUT levels, and on an FPGA every level saved removes a LUT and a routing hop from the path. Enumerating all cuts is exponential in the worst case, so practical mappers keep only a few priority cuts per node, which loses optimality rarely. Depth-optimal covers can duplicate logic at reconvergent fanout, so the next step is area recovery that keeps depth while choosing cheaper cuts.

## Going deeper

Cong and Ding, "FlowMap: An Optimal Technology Mapping Algorithm for Delay Optimization in Lookup-Table Based FPGA Designs", IEEE TCAD (1994); Mishchenko, Cho, Chatterjee and Brayton, "Combinational and Sequential Mapping with Priority Cuts", ICCAD (2007).
