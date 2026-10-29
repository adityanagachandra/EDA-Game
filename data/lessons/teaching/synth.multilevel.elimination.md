# Node elimination and its literal value

Elimination removes a node by substituting its function into every fanout. It undoes sharing, which sounds wasteful, but it exposes optimization across node boundaries: after collapsing, a larger node can be re-minimized and re-factored. SIS-style scripts eliminate nodes whose removal costs few literals, then re-extract. The value of a node decides which ones go.

## The value of a node

Suppose node $t$ has a factored form with $l$ literals and appears as a single literal in $n$ fanout expressions. Eliminating it replaces each use by a copy of the factored form, so the network gains $nl$ literals and loses the $l$ literals of $t$ plus the $n$ uses: the change is $nl - n - l = (n - 1)(l - 1) - 1$. A node with $l = 3$ and $n = 5$ has value $15 - 8 = 7$. A node used once has value -1: eliminating it always saves the literal of its use. This is the mirror image of the extraction saving, which is why extraction and elimination are paired in scripts.

## Threshold elimination

The SIS command eliminate with a threshold $T$ collapses every node whose value is at most $T$. A threshold of -1 removes only nodes used once; a threshold of 0 or more also removes small shared nodes; large thresholds collapse almost everything. Given candidates $t_1$ (2 literals, 1 use, value -1), $t_2$ (4 literals, 1 use, value -1), $t_3$ (3 literals, 2 uses, value 1), $t_4$ (4 literals, 2 uses, value 2), a threshold of 1 eliminates $t_1$, $t_2$ and $t_3$. Values change as nodes are eliminated, since collapsing one node can change the literal counts and use counts of others, so the pass recomputes them as it goes.

## Exact change in SOP literals

In SOP terms elimination multiplies out. If $x$ has $c$ cubes and $L$ literals and appears in a fanout cube $qx$ where $q$ has $l_q$ literals, that cube becomes $c$ cubes with $l_q c + L$ literals in total. For $x = ac + bc$ ($c = 2$, $L = 4$) used in $ex$, $gx$ and $dgx$: before, $4 + 2 + 2 + 3 = 11$ literals; after, $(2 + 4) + (2 + 4) + (4 + 4) = 20$. The change is $+9$. SOP literal growth is larger than factored growth because multiplying out loses the factoring. Collapsing is the first step of refactoring and of node simplification with don't cares, so a temporary increase is often accepted.

## Going deeper

Sentovich et al., "SIS: A system for sequential circuit synthesis", UCB/ERL M92/41 (1992), on eliminate and the standard scripts; De Micheli (1994), the chapter on multiple-level optimization.
