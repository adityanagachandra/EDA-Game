# The ite and apply operations

Every operation a BDD package offers, from AND and XOR to existential quantification, is built on one recursive procedure: split both arguments on the top variable, recurse on the cofactors, and rebuild through the unique table. The if-then-else operator ite packages this into a single primitive. Knowing how apply recurses, and why a computed table bounds its cost, is how you predict the run time of symbolic algorithms such as image computation.

## ite as the universal operator

Define $\text{ite}(f, g, h) = f g + \neg f h$: if $f$ then $g$ else $h$. Every binary connective is an ite with constant or negated arguments: $f \wedge g = \text{ite}(f, g, 0)$, $f \vee g = \text{ite}(f, 1, g)$, $f \oplus g = \text{ite}(f, \neg g, g)$ and $\neg f = \text{ite}(f, 0, 1)$. A BDD node testing $x$ with children $L$ and $H$ is itself $\text{ite}(x, H, L)$, so the package needs only one core algorithm. Its result is always the canonical ROBDD of the combined function, so its size can be predicted by building that function directly. Take $f = x_1 \wedge x_2$ and $g = x_3$ under $x_1 < x_2 < x_3$. Then $f \vee g$ has a root on $x_1$ whose low child is the $x_3$ node and whose high child is an $x_2$ node with children $x_3$ and 1: 3 nonterminal nodes.

## Apply for XOR

Apply computes $u \text{ op } v$ recursively. If both arguments are terminals, return the terminal for the operation's value. Otherwise let $x$ be the earliest variable in the order at the top of $u$ or $v$, take the cofactors of each argument with respect to $x$ (an argument that does not test $x$ is its own cofactor), recurse on the two pairs, and return the unique node $(x, \text{low}, \text{high})$, or the child itself if both are equal. For $f = x_1 \vee x_2$ and $g = x_2 \vee x_3$ under $x_1 < x_2 < x_3$, $f \oplus g$ is 1 exactly on 001 and 100, which is $\neg x_2 \wedge (x_1 \oplus x_3)$. Its ROBDD has a root on $x_1$, two $x_2$ nodes (one leading to $x_3$, one to $\neg x_3$) and two $x_3$ nodes: 5 nonterminal nodes. XOR rarely simplifies at terminals, which is why it tends to cost more than AND or OR.

## The computed table and the product bound

Without memoization apply would revisit the same pair of nodes many times along different paths and could take exponential time. The computed table caches the result of each call keyed by (op, $u$, $v$), so each pair of nodes is processed at most once. Every call's arguments are one node of $F$ and one node of $G$, so the number of distinct calls is at most $|F| \cdot |G|$, with sizes counting terminals: for 19 and 12 nodes, at most 228 calls. The result can also be as large as the product, as with some XORs. Three-argument ite is bounded by $|F| \cdot |G| \cdot |H|$. Real packages make the computed table a lossy hash cache rather than a complete memo, trading occasional recomputation for bounded memory; the bound is still the right model for planning.

## Going deeper

Brace, Rudell and Bryant, "Efficient Implementation of a BDD Package" (DAC 1990); Bryant, "Graph-Based Algorithms for Boolean Function Manipulation" (1986).
