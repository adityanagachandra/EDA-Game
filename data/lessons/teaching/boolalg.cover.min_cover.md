# Minimum covers and the covering problem

Once the primes are known, two-level minimization becomes a covering problem: choose the cheapest set of primes whose union is the ON-set. This is unate covering, a form of set cover, and it is NP-hard. Exact minimizers solve it by branch and bound with reductions; Espresso avoids building it by iterating expand, irredundant and reduce. The covering formulation reappears across EDA, in test compaction, technology mapping and resource binding, so it is worth learning on the smallest instances.

## Covering the ON-set with primes

Build a table with one column per ON minterm and one row per prime, marking which primes cover which minterms. A cover is a set of rows that marks every column. With cost one per cube, the goal is the fewest rows.

For $f(a, b, c) = \sum m(1, 2, 4, 5, 6, 7)$ the primes are $a$, $b\overline{c}$ and $\overline{b}c$. Minterm 1 needs $\overline{b}c$, minterm 2 needs $b\overline{c}$, and minterms 4 to 7 are covered by $a$. All three primes are needed, and $f = a + b\overline{c} + \overline{b}c$ is minimum with 3 cubes. As an integer program: minimize $\sum_p x_p$ subject to $\sum_{p \ni m} x_p \ge 1$ for every ON minterm $m$, with $x_p \in \{0, 1\}$.

## The cyclic core

Essential primes are forced, so select them and delete their minterms. What remains, after dominance reductions, is the cyclic core, where every remaining minterm has at least two candidate primes. The minimum cover is the essentials plus a minimum cover of the core.

For $f = \sum m(0, 1, 2, 5, 6, 7)$ there are six primes, $\overline{a}\,\overline{b}$, $\overline{b}c$, $ac$, $ab$, $b\overline{c}$ and $\overline{a}\,\overline{c}$, none essential, each covering two minterms around a cycle. Six minterms with two per prime need at least 3 primes, and $\overline{a}\,\overline{b} + b\overline{c} + ac$ achieves it. So the core needs 3 primes beyond the 0 essentials. Petrick's method writes the core as a product of sums and expands it, which is exact but exponential.

## Cost beyond cube count

Cubes are not the only cost. A PLA's area depends on cube count, but a multi-level netlist cares about literals, and the usual tie-breaker is the total literal count among covers with the fewest cubes. For $f = \sum m(1, 3, 5, 6, 7) = c + ab$ the minimum has 2 cubes and 3 literals; the non-prime cover $c + ab\overline{c}$ also has 2 cubes but 4 literals, which is one reason exact minimizers draw rows only from primes.

Exact solvers (the classic one is in Espresso-exact, and later ones use implicit BDD-based representations of the table) prune with lower bounds from independent sets of columns. In the example above, minterms 0, 5 and 6 share no prime, a lower bound of 3 that proves the cover optimal without search.

## Going deeper

Coudert, "Two-level logic minimization: an overview", Integration, the VLSI Journal (1994). De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 7.
