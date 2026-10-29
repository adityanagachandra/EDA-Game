# Simple disjoint decomposition

A decomposition writes a function as a composition of smaller ones. The simple disjoint form is $f(A, B) = h(g(A), B)$: a sub-function $g$ reads only the bound set $A$, and the rest of the logic sees $A$ only through that single wire. Finding such structure is how FPGA mappers fit logic into lookup tables, how synthesis discovers intermediate signals, and how bi-decomposition and functional decomposition tools restructure a circuit. The test for it, due to Ashenhurst and extended by Curtis, is a counting argument on cofactors.

## The decomposition chart

Split the inputs into a bound set $A$ and a free set $B$. Draw the decomposition chart: one column per assignment of $A$, one row per assignment of $B$. Column $\alpha$ is the cofactor $f_{A = \alpha}$, a function of $B$. The column multiplicity $\mu$ is the number of distinct columns. Ashenhurst's theorem: $f = h(g(A), B)$ exists exactly when $\mu \le 2$; then $g$ marks which of the two column types each $\alpha$ has.

Let $f(a, b, c, d) = (a \oplus b)\,c + d$ with $A = \{a, b\}$. The columns are $d$ for $ab = 00$ and 11, and $c + d$ for 01 and 10: two distinct columns, so $f$ decomposes with $g = a \oplus b$ and $h = g\,c + d$.

## How many intermediate signals

When $\mu > 2$, $A$ can still be compressed through $k$ wires: $f = h(g_1(A), \dots, g_k(A), B)$ with $k = \lceil \log_2 \mu \rceil$, because the $g_i$ must give each distinct column its own code. This is Curtis decomposition, and it pays off whenever $k < |A|$.

For the same $f$ with bound set $\{c, d\}$, the columns are 0 for $cd = 00$, 1 for 01 and 11, and $a \oplus b$ for 10: $\mu = 3$, so $k = \lceil \log_2 3 \rceil = 2$, no better than passing $c$ and $d$ through. If $\mu = 1$, $f$ does not depend on $A$ at all and $k = 0$.

## Searching for a good bound set

The decomposition depends entirely on the choice of $A$, and there are $\binom{n}{|A|}$ candidates for each size. Tools enumerate small bound sets, build each chart from cofactors (on truth tables or BDDs, where $\mu$ is the number of distinct nodes below a cut in an order that puts $A$ on top), and keep the ones with small multiplicity.

For $f = (a \oplus b)\,c + d$, only $\{a, b\}$ among the six two-variable bound sets gives $\mu \le 2$. The search is exponential in the worst case, so practical decomposers restrict bound-set size to the target LUT width and prune with symmetries: symmetric variables are interchangeable in the bound set.

## Going deeper

Ashenhurst, "The decomposition of switching functions" (1957). Curtis, "A New Approach to the Design of Switching Circuits" (1962). Perkowski and Grygiel, "A survey of literature on function decomposition" (1995).
