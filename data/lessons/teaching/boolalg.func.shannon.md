# Shannon expansion

Shannon expansion writes any Boolean function as a choice between its two cofactors: $f = x\, f_x + \overline{x}\, f_{\overline{x}}$. It is the identity behind every recursive Boolean algorithm. A multiplexer implements it in hardware, a BDD node stores it, a tautology checker recurses on it, and a DPLL search tree is its unfolding. Knowing how to split, rebuild and count with it is the base of decision diagrams and of complexity arguments about Boolean procedures.

## Splitting a function on one variable

For any variable $x$, $f = x\, f_x + \overline{x}\, f_{\overline{x}}$. Check it on a row: when $x = 1$ the second term vanishes and the first gives $f_x$ at that row, which is $f$; when $x = 0$ the roles swap. The two terms are disjoint, since one needs $x = 1$ and the other $x = 0$, so the ON count of $f$ is the sum of the ON counts of its cofactors.

For $f(a, b, c) = ab + \overline{c}$, $f = a\,(b + \overline{c}) + \overline{a}\,(\overline{c})$. The positive cofactor has 3 ON rows over $b, c$ and the negative one has 2, so $f$ has $3 + 2 = 5$ ON rows. The dual form is $f = (x + f_{\overline{x}})(\overline{x} + f_x)$.

## Rebuilding a function from its cofactors

Expansion runs backwards: given $f_x$ and $f_{\overline{x}}$, $f$ is determined. Each minterm of the positive cofactor becomes a minterm of $f$ by inserting $x = 1$ at $x$'s bit position, and each minterm of the negative cofactor by inserting $x = 0$.

Over $a, b, c$, let $f_b = m(1, 3)$ and $f_{\overline{b}} = m(0)$, both over $a, c$. Minterm 1 of $f_b$ is $(a, c) = (0, 1)$; inserting $b = 1$ gives $011 = 3$. Minterm 3 is $(1, 1)$, giving $111 = 7$. Minterm 0 of $f_{\overline{b}}$ is $(0, 0)$, giving $000 = 0$. So $f = m(0, 3, 7)$.

## Recursive expansion and distinct subfunctions

Expanding on $x_1$, then on $x_2$ in each cofactor, and so on, builds a full binary tree of $2^{k}$ subfunctions at depth $k$. Many of them coincide. For $f = ab + \overline{c}$, depth 1 gives $\overline{c}$ and $b + \overline{c}$, two distinct functions; depth 2 gives $\overline{c}, \overline{c}, \overline{c}, 1$, again only two distinct.

Counting distinct subfunctions per level is how a tool sizes the shared structure: a reduced ordered BDD stores each distinct subfunction once, so the number of distinct cofactors at each depth bounds its width. A full Shannon tree without sharing has $2^{n+1} - 1$ nodes, which is why naive recursion is exponential and memoized recursion is not, for many functions.

## Going deeper

Shannon, "A Symbolic Analysis of Relay and Switching Circuits" (1938). Bryant, "Graph-Based Algorithms for Boolean Function Manipulation", IEEE Transactions on Computers (1986).
