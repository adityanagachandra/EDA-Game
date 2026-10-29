# Polynomial-time reductions

A reduction from problem A to problem B is a polynomial-time transformation of every instance of A into an instance of B with the same yes or no answer. Reductions are how complexity transfers: they prove that a new problem is hard, and they are also how tools solve problems, by encoding them into SAT, ILP or graph problems that have good solvers. The two uses run in opposite directions, and confusing them is the most common error in arguments about EDA problem complexity.

## What a reduction must preserve

A many-one (Karp) reduction is a function $r$ computable in polynomial time such that $x \in A$ exactly when $r(x) \in B$. Write $A \le_p B$. It must map yes-instances to yes-instances and no-instances to no-instances; it need not preserve the number of solutions, although parsimonious reductions do, which matters for counting.

Two consequences follow. If $A \le_p B$ and B is in P, then A is in P: run $r$, then B's algorithm, polynomial composed with polynomial. If $A \le_p B$ and A is NP-hard, then B is NP-hard. Reductions compose, which is how chains such as SAT to 3-SAT to independent set to vertex cover are built.

## 3-SAT to independent set

Karp's reduction maps a 3-CNF with $m$ clauses to a graph: one vertex per literal occurrence, a triangle joining the three literals of each clause, and an edge between every pair of complementary literals in different clauses. Ask for an independent set of size $k = m$. An independent set can take at most one vertex per triangle, so size $m$ means one literal per clause, and the complementary edges forbid choosing both $x$ and $\neg x$: such a set exists exactly when the formula is satisfiable.

For $(x_1 \lor x_2 \lor x_3)(\neg x_1 \lor \neg x_2 \lor x_3)$: 6 vertices, 6 triangle edges, and 2 complementary edges ($x_1$ with $\neg x_1$, $x_2$ with $\neg x_2$), so 8 edges and $k = 2$. Picking $x_3$ in both clauses is an independent set of size 2, matching the model $x_3 = 1$.

## Choosing the direction

To show a new problem X is hard, reduce a known hard problem to X: every instance of the hard problem becomes an instance of X. To solve X with existing tools, reduce X to a solved problem, for example encode X as SAT. The second direction gives an algorithm and an upper bound and says nothing about hardness.

So encoding clock-gating legality or test generation into CNF shows those problems are in NP, provided the encoding is polynomial, while Ibarra and Sahni's construction of a circuit whose fault is detectable exactly when a given circuit is satisfiable shows test generation is NP-hard. An argument that reduces X to SAT and concludes "X is NP-hard" has the arrow backwards.

## Going deeper

Karp, "Reducibility among combinatorial problems" (1972). Garey and Johnson, "Computers and Intractability: A Guide to the Theory of NP-Completeness" (1979), chapter 3. Ibarra and Sahni, "Polynomially complete fault detection problems", IEEE Transactions on Computers (1975).
