# Unate covering

Choosing a minimum set of primes is a unate covering problem: select rows of a 0-1 matrix so that every column is covered, at minimum cost. It is NP-hard, it appears all over EDA (test compaction, technology mapping, routing track assignment), and exact minimizers solve it with branch and bound guided by lower bounds. Writing it as an integer program makes its structure plain.

## The covering problem as an integer program

Give each row $i$ a binary variable $x_i$ (1 if selected) and a cost $w_i$. For each column $j$ add the constraint $\sum_{i \in R(j)} x_i \ge 1$: at least one row marking column $j$ is chosen. Minimize $\sum_i w_i x_i$. The constraints are **unate**: every variable appears uncomplemented, so turning a row on can never violate a constraint. A selection is feasible exactly when no column constraint is violated. For a candidate selection, counting violated constraints means counting columns none of the selected rows marks; that count is zero for a cover. The binate covering problem, which allows complemented variables, arises in state minimization and technology mapping with choices and is harder still.

## Solving by branch and bound

Exact solvers reduce the matrix (essentials, dominance), then branch: pick a column, and for each row that covers it, try including that row and recurse. A **lower bound** prunes branches: if the cost so far plus a bound on the remaining cost is no better than the best cover found, stop. A standard bound is a maximal set of pairwise disjoint columns (no row covers two of them): each needs its own row, so their count bounds the number of rows still needed. The LP relaxation, with $0 \le x_i \le 1$, gives another bound. On small tables an exhaustive search over subsets in order of size gives the minimum directly, which is how the drills here are checked.

## Costs and weighted covers

In logic minimization the row cost is usually the literal count of the prime, or one per cube when the goal is the fewest products. The two objectives can disagree: a cover with fewer cubes can have more literals. With weights the minimum is no longer the smallest set of rows. If rows cost 3, 3, 3 and 4, a two-row cover that uses the 4-cost row costs 7 while another two-row cover costs 6, and a three-row cover of cheap rows can beat a two-row cover of expensive ones. Weighted covering also changes dominance: a row can only be dropped in favor of a row that covers more and costs no more. Tools such as Espresso-exact and the covering solvers inside SIS use these weighted forms.

## Going deeper

Coudert and Madre, "New ideas for solving covering problems", DAC 1995; Hachtel and Somenzi (1996), on covering; any integer programming text on set covering.
