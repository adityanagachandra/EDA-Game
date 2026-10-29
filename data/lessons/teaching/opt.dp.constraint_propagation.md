# Constraint satisfaction and propagation

Constraint programming states a problem as variables with finite domains and constraints between them, then alternates propagation, which removes values that cannot appear in any solution, with search over what is left. It handles scheduling, register and resource binding, channel and track assignment, and configuration problems where rules are irregular and objectives secondary. Arc consistency, the basic propagation, often solves precedence-heavy instances outright, and it is the same reasoning that unit propagation performs inside SAT solvers.

## Bounds propagation on precedence chains

Start times $s_1, s_2, s_3$ with domains $\{0..7\}$, $\{0..6\}$, $\{0..9\}$ and precedences $s_1 + 1 \le s_2$ and $s_2 + 3 \le s_3$. A constraint $s_i + d \le s_j$ removes small values of $s_j$ (below $\min s_i + d$) and large values of $s_i$ (above $\max s_j - d$). Forward: $s_2 \ge 1$ and $s_3 \ge 4$. Backward: $s_2 \le 9 - 3 = 6$, already true, and $s_1 \le 6 - 1 = 5$.

The final domains are $s_1 \in \{0..5\}$, $s_2 \in \{1..6\}$ and $s_3 \in \{4..9\}$, so $s_2$ keeps 6 values. Propagate until nothing changes, since a tightened bound can trigger further tightening along the chain. For interval domains and binary linear precedences this bounds reasoning is exactly arc consistency, and it computes the same ASAP and ALAP windows as scheduling.

## Counting solutions of a small CSP

Three operations a, b, c each take a step from 1 to 5. Constraint: b before c, and a and b share a unit, so $a \ne b$. Count: pairs with $b < c$ number $\binom{5}{2} = 10$, and for each, a can be any step except b, 4 choices, giving 40 solutions out of $5^3 = 125$ assignments.

Counting by hand mirrors how a solver reasons: fix the most constrained variables first, then count the free choices for the rest. Real CP solvers do not enumerate. They pick a variable (often the one with the smallest domain, the fail-first heuristic), try a value, propagate, and backtrack on a wiped-out domain. Global constraints carry stronger propagators: all-different reasons about matchings, and cumulative reasons about resource profiles over time, which is how CP schedulers handle functional-unit limits far better than pairwise $\ne$ constraints can.

## What propagation buys

Measure propagation by the values it removes. Domains $s_1 \in \{2..11\}$, $s_2, s_3 \in \{0..13\}$ with $s_1 + 4 \le s_2$ and $s_2 + 1 \le s_3$ propagate to $\{2..8\}$, $\{6..12\}$ and $\{7..13\}$: from 10, 14 and 14 values to 7, 7 and 7, removing 17 values without any search. Every removed value is a branch the search never explores.

The trade-off is propagation strength against cost per node. Arc consistency on binary constraints is cheap; stronger consistencies remove more values but cost more. CP shines when constraints are irregular, such as pin-access rules, cell-placement patterns or test-scheduling restrictions, and weakest on problems with large linear objectives, where LP bounds in a MILP solver prune better. Lazy clause generation combines CP propagation with SAT-style clause learning, and is one reason CP-SAT solvers have become a strong default for EDA scheduling and assignment problems.

## Going deeper

Rossi, van Beek and Walsh (eds.), "Handbook of Constraint Programming" (2006), chapters on constraint propagation and global constraints; Mackworth, "Consistency in networks of relations" (Artificial Intelligence, 1977).
