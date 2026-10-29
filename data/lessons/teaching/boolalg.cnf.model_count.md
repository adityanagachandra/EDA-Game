# Counting satisfying assignments

Model counting, written #SAT, asks how many assignments satisfy a formula rather than whether one exists. It is harder than SAT (#P-complete, even for monotone 2-CNF), yet it is useful in EDA: signal probabilities for power estimation, coverage metrics, the fraction of input space a constraint leaves for constrained-random simulation, and the probability that random simulation misses a bug all are model counts. On small formulas three techniques suffice: enumerate, use inclusion and exclusion, and multiply over independent components.

## Model counting by enumeration

With $n$ variables, test all $2^n$ assignments against every clause and count the survivors. For $(x_1 \lor x_2) \land (\neg x_1 \lor x_3)$ over $x_1, x_2, x_3$: with $x_1 = 0$ the second clause holds and the first needs $x_2 = 1$, giving 2 models ($x_3$ free); with $x_1 = 1$ the first holds and the second needs $x_3 = 1$, giving 2 models ($x_2$ free). Total 4 of 8.

Enumeration is what the count means, and it scales to perhaps 30 variables with bit-parallel evaluation. Exact counters beyond that use DPLL-style search with component caching, or compile the formula into a d-DNNF or BDD on which counting is linear in the size of the compiled form.

## Counting by inclusion and exclusion

Let $A_j$ be the set of assignments falsifying clause $j$. A clause with $k$ literals has $|A_j| = 2^{n-k}$. Several clauses can be falsified together only if they contain no complementary pair; then the intersection fixes every variable they mention, so $|A_i \cap A_j| = 2^{n - |vars(C_i) \cup vars(C_j)|}$, and otherwise it is 0. The model count is $2^n - |A_1 \cup \dots \cup A_m|$.

For the example: each two-literal clause is falsified by $2^{3-2} = 2$ assignments, and the pair contains $x_1$ and $\neg x_1$, so they are never falsified together. The union has $2 + 2 = 4$ assignments and the count is $8 - 4 = 4$, as before. With many clauses the number of terms is exponential, so this is a tool for small cases and for reasoning, not a general algorithm.

## Independent components multiply

If the clauses split into groups that share no variables, the groups constrain independent coordinates, so the count is the product of the per-group counts times 2 for each variable that appears in no clause. For $(x_1 \lor x_2) \land (x_3 \lor x_4)$ over $x_1, \dots, x_5$: each clause has 3 models over its own two variables, and $x_5$ is free, so $3 \cdot 3 \cdot 2 = 18$.

Component decomposition is the main trick of practical exact counters: after a few decisions the residual formula often falls apart into independent pieces, each counted once and cached. The same idea, independence, is why signal probability estimates are exact on fanout-free circuits and only approximate when reconvergent fanout correlates the inputs of a gate.

## Going deeper

Gomes, Sabharwal and Selman, "Model counting", in Biere et al., "Handbook of Satisfiability" (2nd ed.). Valiant, "The complexity of enumeration and reliability problems", SIAM Journal on Computing (1979).
