# K-feasible cuts

A cut of a node is a set of nodes that separates it from the primary inputs, and the node's function over a cut is a small Boolean function. With at most $k$ leaves, a cut is a window that fits in a $k$-input truth table. Cuts are how rewriting finds local functions to replace, how LUT mappers choose LUTs, and how resubstitution and SAT sweeping localize their work.

## Definition

A **cut** of node $n$ is a set $C$ of nodes such that every path from a primary input to $n$ passes through a node of $C$. It is **$k$-feasible** if $|C| \le k$. The trivial cut $\{n\}$ is always a cut. Given the cut, the node's function can be computed over the leaves only, by simulating the cone between $C$ and $n$, giving a truth table of $2^{|C|}$ bits. For $n_4 = AND(d, n_1)$ with $n_1 = AND(b, d)$, the 3-feasible cuts are $\{n_4\}$, $\{d, n_1\}$ and $\{b, d\}$ (the merged set $\{d, b, d\}$ has only two distinct leaves), so 3 cuts.

## Bottom-up enumeration

Cuts are enumerated in topological order. A primary input has only its trivial cut. For an AND node $n$ with fanins $a$ and $b$, every nontrivial cut of $n$ is a union $u \cup v$ with $u$ a cut of $a$ and $v$ a cut of $b$. So $\mathrm{cuts}(n) = \{\{n\}\} \cup \{u \cup v : u \in \mathrm{cuts}(a), v \in \mathrm{cuts}(b), |u \cup v| \le k\}$, with duplicate sets removed. Reconvergence makes unions smaller than the sum of their parts, which is why counting needs care. The number of cuts can grow exponentially in $k$, so practical enumeration keeps $k$ at 4 to 6 for rewriting and mapping, and stores truth tables as machine words.

## Dominance and priority cuts

A cut $C'$ **dominates** $C$ if $C' \subset C$: any function expressible over $C$ restricted to that cone is determined by $C'$ too, so $C$ is redundant for mapping and is dropped. Removing dominated cuts shrinks the sets propagated upward. Large designs still produce too many cuts, so **priority cuts** keep only a few best cuts per node (by depth, then area flow), accepting that some optimal choices are lost. That trade made cut-based mapping scale to millions of nodes, and it is ABC's default in the if mapper. For a node with 5 three-feasible cuts, typically all or nearly all are non-dominated; in larger windows dominance removes many.

## Going deeper

Cong, Wu and Ding, "Cut ranking and pruning: enabling a general and efficient FPGA mapping solution", FPGA (1999); Mishchenko, Cho, Chatterjee and Brayton, "Combinational and sequential mapping with priority cuts", ICCAD (2007).
