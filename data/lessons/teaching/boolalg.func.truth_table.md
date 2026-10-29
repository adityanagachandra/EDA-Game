# Boolean functions and truth tables

A Boolean function $f : \{0,1\}^n \to \{0,1\}$ is the object every logic tool manipulates: synthesis rewrites it, a SAT solver searches its inputs, an equivalence checker compares two of them, and a test generator looks for inputs that expose a difference. The truth table is its most explicit representation, one output bit per input row. It is exponential in $n$, which is exactly why the rest of this domain exists, but it is the ground truth every compact representation must agree with.

## Functions from n bits to one bit

Number the rows of the table by minterms. Row $i$ writes $i$ in binary over the variables in a fixed order, the first variable as the most significant bit, so over $a, b, c$ row 5 is $a = 1, b = 0, c = 1$. The ON-set is the set of rows where $f = 1$, the OFF-set the rest, and a function is completely described by its ON-set: $f = \sum m(\dots)$.

Take $f(a, b, c) = ab + \overline{c}$. Row by row: rows 0, 2, 4 and 6 have $c = 0$, so $f = 1$; row 7 has $a = b = 1$, so $f = 1$; rows 1, 3 and 5 give 0. The ON-set is $m(0, 2, 4, 6, 7)$, five rows of eight. Evaluating an expression on every row is the slow but certain way to learn what a function is.

## Counting Boolean functions

A table over $n$ inputs has $2^n$ rows and each output bit is free, so there are $2^{2^n}$ functions of $n$ variables: 16 for $n = 2$, 256 for $n = 3$, 65,536 for $n = 4$ and about $4.3 \times 10^9$ for $n = 5$. A specification that pins the output on $r$ rows leaves $2^{2^n - r}$ functions. The number of functions with exactly $k$ ON rows is $\binom{2^n}{k}$; for $n = 4$ and $k = 2$ that is $\binom{16}{2} = 120$.

These counts explain a design choice in tools such as ABC: all 65,536 functions of four inputs can be enumerated and precomputed, so four-input cuts are rewritten from a table, while five or six inputs already need canonical forms and hashing.

## Comparing two functions row by row

Two functions are equal when their tables agree on every row. The rows where they differ are the ON-set of $f \oplus g$, and their number is the Hamming distance between the two tables. Let $f = m(0, 2, 4, 6, 7)$ as above and $g = \overline{c} = m(0, 2, 4, 6)$. They differ only on row 7, distance 1.

That count is the number of counterexamples an equivalence checker could return. A distance of 0 means equivalent; a small positive distance is the dangerous case, because random simulation is unlikely to hit one bad row among millions.

## Going deeper

Knuth, "The Art of Computer Programming", volume 4A, section 7.1.1 (Boolean basics). Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", chapter 2.
