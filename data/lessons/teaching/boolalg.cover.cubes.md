# Cubes and their minterms

A cube is a product of literals, such as $a\overline{c}d$, and also the set of minterms it covers, a subcube of the Boolean hypercube. Covers of cubes are the data structure of two-level logic: Espresso, the PLA format, and the cube-based reasoning inside many SAT preprocessors all store functions as lists of cubes. The basic operations, size, intersection, union, distance and consensus, are bit operations on a positional encoding.

## Cubes as products of literals

In positional notation each variable gets one symbol: 1 for a true literal, 0 for a complemented one, and a dash for an absent variable. Over $a, b, c, d$ the cube $a\overline{c}$ is 1-0-. A cube with $k$ literals over $n$ variables has $n - k$ dashes and covers $2^{n-k}$ minterms; 1-0- covers 4. Implementations encode each position in two bits, one per allowed value, so 0, 1 and dash become 10, 01 and 11, and intersecting two cubes is a bitwise AND over machine words.

Counting cubes is also easy: each position takes one of three symbols, so there are $3^n$ cubes over $n$ variables, and $\binom{n}{k} 2^k$ of them have exactly $k$ literals.

## Intersection and union of cubes

The intersection of two cubes is computed position by position: a dash meets anything and yields the other symbol, equal symbols stay, and 0 against 1 makes the intersection empty. So 1-0- and -10- intersect in 110-, while 1-0- and 0--- are disjoint.

The union of several cubes is counted by inclusion and exclusion. Over four variables, 1--- and -1-- each cover 8 minterms and intersect in 11--, which covers 4, so the union covers $8 + 8 - 4 = 12$. With three cubes, add back the triple intersection. This is how a minimizer measures how much of the ON-set a candidate cover still misses.

## Distance and consensus

The distance between two cubes is the number of positions where one has 0 and the other 1. Distance 0 means they intersect. At distance 1 the cubes are adjacent and have a consensus: drop the conflicting variable and intersect the rest. For $a\overline{b}$ and $\overline{a}c$ over $a, b, c$ (10- and 0-1), the conflict is at $a$ and the consensus is -01, that is $\overline{b}c$, covering 2 minterms.

The consensus theorem $xy + \overline{x}z = xy + \overline{x}z + yz$ says adding the consensus never changes the function. Iterating consensus and removing contained cubes generates all prime implicants, and the same operation on clauses is resolution, the engine of SAT solving.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984), chapter 2. Rudell, "Logic Synthesis for VLSI Design", PhD thesis, UC Berkeley (1989).
