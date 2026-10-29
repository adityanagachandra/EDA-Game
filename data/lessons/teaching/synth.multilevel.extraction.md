# Extracting common kernels

Extraction finds a divisor shared by several functions, creates one node for it, and rewrites each function to use that node. It is the multi-level move that saves the most area, because logic computed once is reused many times. Kernel extraction (the gkx command in SIS) is the multi-cube version; its literal accounting decides whether a candidate is worth extracting.

## Literals before extraction

Start from the SOPs of the functions that share a divisor. For $F_1 = bcde + ade + g$ and $F_2 = bcd + ad + g$ the common divisor is $D = bc + a$: $F_1 = (bc + a)de + g$ and $F_2 = (bc + a)d + g$. Before extraction the literal count is simply the sum over all cubes: $4 + 3 + 1 = 8$ for $F_1$ and $3 + 2 + 1 = 6$ for $F_2$, 14 in total. The kernels of $F_1$ include $bc + a$ (co-kernel $de$) and those of $F_2$ include $bc + a$ (co-kernel $d$); their intersection is the two-cube divisor $D$.

## Literals after extraction

Create the node $X = D$ and rewrite each function as $XQ + R$, where $Q$ and $R$ come from weak division by $D$. Count the literals of $X$ once, then for each function count each quotient cube's literals plus one for the literal $X$, plus the remainder's literals. For the example: $X = bc + a$ costs 3. $F_1 = Xde + g$ costs $(2 + 1) + 1 = 4$ and $F_2 = Xd + g$ costs $(1 + 1) + 1 = 3$. Total $3 + 4 + 3 = 10$. A quotient of the single cube 1 still costs one literal, the $X$ itself: $F = D + h$ becomes $X + h$ with 2 literals.

## Value, greed and the extraction loop

The value of extracting $D$ is literals before minus literals after: 14 - 10 = 4 in the example. Kernel extraction repeatedly picks the divisor of highest value among kernel intersections, extracts it, and updates the network, until no divisor saves literals. The search can be cast as rectangle covering on a co-kernel/cube matrix, where a rectangle is a set of co-kernels sharing a set of kernel cubes. Fast extract (fx), used in SIS and ABC, restricts candidates to double-cube divisors and single-cube divisors, which can be enumerated and updated incrementally on very large networks. Extraction increases fanout, so a later stage may need to buffer or duplicate the shared node.

## Going deeper

Rajski and Vasudevamurthy, "The testability-preserving concurrent decomposition and factorization of Boolean expressions", IEEE TCAD (1992), which introduced fast extraction; Brayton, Rudell, Sangiovanni-Vincentelli and Wang, "MIS", IEEE TCAD (1987).
