# From SAT to 3-SAT

3-SAT, where every clause has at most three literals, is the standard starting point for NP-hardness proofs because its instances are uniform and small gadgets can simulate a clause. It is NP-complete because general CNF-SAT reduces to it by splitting long clauses with fresh variables. The construction is also a lesson in equisatisfiable encodings: the new formula is not equivalent to the old one, but it is satisfiable exactly when the old one is, and that is all a reduction needs.

## Splitting one long clause

A clause $(l_1 \lor l_2 \lor \dots \lor l_k)$ with $k > 3$ is replaced by a chain linked by fresh variables $y_1, \dots, y_{k-3}$:
$$(l_1 \lor l_2 \lor y_1)(\neg y_1 \lor l_3 \lor y_2)(\neg y_2 \lor l_4 \lor y_3) \cdots (\neg y_{k-3} \lor l_{k-1} \lor l_k).$$
The first and last clauses take two original literals each and every middle clause one, so there are $k - 2$ clauses and $k - 3$ new variables. For $k = 5$: $(l_1 \lor l_2 \lor y_1)(\neg y_1 \lor l_3 \lor y_2)(\neg y_2 \lor l_4 \lor l_5)$, three clauses and two new variables.

If some $l_j$ is true, set $y_i = 1$ for the clauses before it and $y_i = 0$ after it, and every chain clause is satisfied. If every $l_j$ is false, the chain forces $y_1 = 1$, then $y_2 = 1$, and so on, until the last clause fails. So the chain is satisfiable exactly when the original clause is.

## Fresh variables for a whole formula

Each long clause gets its own fresh variables, so the total is $\sum_j (k_j - 3)$ over clauses with $k_j > 3$, and clauses with three or fewer literals add none. A CNF with clauses of sizes 9, 5, 8, 6 and 9 gains $6 + 2 + 5 + 3 + 6 = 22$ variables.

Clauses shorter than three can be padded with fresh variables if exact 3-literal clauses are required, or left alone for at-most-3-SAT. Either way the growth is linear in the size of the input, so the reduction is polynomial.

## Size of the 3-CNF

The new formula has $\sum_j (k_j - 2)$ clauses over the long clauses plus the short clauses unchanged. Sizes 7, 7 and 8 give $5 + 5 + 6 = 16$ clauses. The total number of literal occurrences grows by a factor of at most three, again linear.

The same chaining idea reappears in practice: the sequential-counter encoding of at-most-one constraints, the Tseitin encoding of wide gates and the splitting of long XORs all trade one big constraint for a linear number of small ones linked by auxiliary variables. Fewer literals per clause help unit propagation, but more variables can mean longer search, so solvers and encoders choose case by case.

## Going deeper

Garey and Johnson, "Computers and Intractability" (1979), section 3.1.1. Sipser, "Introduction to the Theory of Computation" (3rd ed., 2012), section 7.4.
