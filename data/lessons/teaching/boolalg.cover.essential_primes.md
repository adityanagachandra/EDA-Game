# Essential prime implicants

A prime implicant is essential when it is the only prime covering some ON minterm. Every sum-of-products cover built from primes must contain it. Exact minimizers select essentials first because they are free decisions, and heuristic minimizers such as Espresso compute them early to shrink the problem. Counting primes and essentials on small functions is the quickest way to feel why two-level minimization is easy for some functions and combinatorially hard for others.

## Generating all primes

The Quine-McCluskey method starts from the ON minterms and repeatedly merges pairs of cubes that are at distance 1 and have their dashes in the same positions, replacing the differing variable by a dash. A cube that never merges with anything is prime. Equivalently, a prime is an implicant contained in no larger implicant.

For $f(a, b, c) = \sum m(0, 1, 5, 7)$: rows 0 and 1 merge into 00-, rows 1 and 5 into -01, rows 5 and 7 into 1-1. None of these three merges further, so the primes are $\overline{a}\,\overline{b}$, $\overline{b}c$ and $ac$. The number of primes can be exponential in $n$, so practical tools avoid listing all of them.

## Essential primes

For each ON minterm, list the primes that cover it. A minterm covered by exactly one prime makes that prime essential. In the example, minterm 0 is covered only by $\overline{a}\,\overline{b}$ and minterm 7 only by $ac$, so both are essential. Minterms 1 and 5 are each covered by two primes, and $\overline{b}c$ is not essential.

Essentiality depends on the whole function, not on the cube alone: adding one ON minterm can create a new prime that steals the uniqueness. For the function $\sum m(0, 1, 2, 5, 6, 7)$ every minterm is covered by two of its six primes, and no prime is essential at all.

## What the essentials leave behind

Select the essentials and remove the minterms they cover. If nothing remains, the essentials alone are the unique minimum cover. Here $\overline{a}\,\overline{b}$ covers 0 and 1, and $ac$ covers 5 and 7, so all four ON minterms are covered and $f = \overline{a}\,\overline{b} + ac$; the non-essential prime $\overline{b}c$ is redundant.

When minterms remain, the leftover covering table is the cyclic core, and choosing primes for it is a set-cover problem. Further reductions (row and column dominance) may shrink it, but in general the core needs search. A tool reports the size of the remaining core because it predicts runtime: an empty core is solved instantly, a large cyclic core is where exact minimization explodes.

## Going deeper

McCluskey, "Minimization of Boolean functions", Bell System Technical Journal (1956). De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 7.
