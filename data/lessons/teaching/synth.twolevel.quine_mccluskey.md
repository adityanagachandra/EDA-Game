# The Quine-McCluskey method

Quine-McCluskey is the tabular, exact algorithm for two-level minimization. Its first phase generates every prime implicant by repeatedly merging cubes that differ in one variable; its second phase chooses a minimum subset of primes, a covering problem. It is the textbook baseline that defines what "exact" means, and its first phase shows why exact methods blow up.

## Merging adjacent minterms

List the ON minterms (and don't cares) in binary and group them by the number of 1s. Two entries in neighboring groups that differ in exactly one bit merge into a cube with a dash in that position: 0101 and 0111 give 01-1. Only adjacent groups need to be compared, since a one-bit difference changes the count of 1s by exactly one. The number of one-dash cubes produced in this first column equals the number of pairs of adjacent ON minterms. For $\sum m(0, 1, 4, 5, 6, 7, 8, 9, 12)$ there are 12 such pairs, among them 0-1, 0-4, 0-8, 1-5, 1-9 and 4-5. Each entry that takes part in a merge is checked off.

## Iterating to the prime implicants

Repeat on the new column: two cubes merge if they have dashes in the same positions and differ in exactly one specified bit, producing a cube with one more dash. Duplicates (the same cube reached in two ways) are kept once. When no more merges are possible, the entries never checked off are the **prime implicants**: implicants that no larger implicant contains. For a function of $n$ variables the number of primes can grow as roughly $3^n/n$ in the worst case, and the merging tables grow with it. That growth is the reason exact two-level minimization is reserved for small functions and why Espresso never enumerates all primes.

## Don't cares and the hand-off to covering

Don't cares join the merging phase as if they were 1s, which can only enlarge primes. They do not join the covering phase: a prime that covers only don't cares is dropped, and the table has one column per ON minterm only. For $\sum m(4, 6, 7, 10, 15)$ plus don't cares, the primes that touch the ON-set become the rows of the prime implicant table. The second phase then picks a minimum set of rows covering every column. That phase is unate covering, NP-hard in general, attacked with essential-prime selection, dominance and branch and bound. Quine-McCluskey thus splits minimization into prime generation and set covering.

## Going deeper

McCluskey, "Minimization of Boolean functions", Bell System Technical Journal (1956); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), the chapter on two-level optimization.
