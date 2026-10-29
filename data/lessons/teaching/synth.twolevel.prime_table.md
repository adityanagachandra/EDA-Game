# Prime implicant tables and essential primes

The second phase of exact two-level minimization works on a table: one row per prime implicant, one column per ON minterm, a mark where the prime covers the minterm. Some primes are forced into every solution because they alone cover some minterm. Finding these essential primes is the first, cheapest step of every covering algorithm, and it often solves most of the problem.

## Building the table

Take the primes from prime generation as rows and the ON minterms as columns, and mark row $p$, column $m$ when cube $p$ contains $m$. Don't-care minterms get no column, because they need not be covered. For $f = \sum m(3, 6, 8, 10, 11, 14)$ the primes are $\bar{b}cd$, $bc\bar{d}$, $ac\bar{d}$, $a\bar{b}\bar{d}$ and $a\bar{b}c$. The column for minterm 10 (1010) is marked by $ac\bar{d}$, $a\bar{b}\bar{d}$ and $a\bar{b}c$, three rows. The number of marks in a column is the number of ways that minterm can be covered, and a selection of rows is a valid cover exactly when every column has a mark in some selected row.

## Essential primes

A column with a single mark identifies an **essential prime**: that row must be in every cover. For $f = \sum m(2, 3, 4, 5, 6, 10, 11, 12)$ the primes are $\bar{b}c$, $b\bar{c}\bar{d}$, $\bar{a}c\bar{d}$, $\bar{a}b\bar{d}$ and $\bar{a}b\bar{c}$. Minterm 11 is covered only by $\bar{b}c$, minterm 12 only by $b\bar{c}\bar{d}$, and minterm 5 only by $\bar{a}b\bar{c}$, so those three are essential. Selecting an essential prime removes every column it covers. Here the three essentials cover every ON minterm except 6, which either $\bar{a}c\bar{d}$ or $\bar{a}b\bar{d}$ can cover, so the minimum has four products and two equally good completions.

## The residual problem

After removing essential rows and their columns, what remains is a smaller covering problem. If no columns remain, the essentials are the unique minimum cover. Otherwise the residue may contain new essentials (secondary essentials after dominance reductions) or may be **cyclic**: every column marked at least twice, no dominance applicable. The classic cyclic example is $f = \sum m(0, 1, 2, 5, 6, 7)$ over three variables, which has six primes, none essential, and two different minimum covers of three cubes each. Cyclic cores are where branch and bound or Petrick's method are needed, and where the exponential cost of exact minimization lives.

## Going deeper

Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996), on exact two-level minimization; De Micheli (1994), the chapter on two-level optimization.
