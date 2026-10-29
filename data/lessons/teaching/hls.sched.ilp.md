# ILP and SDC scheduling

Exact scheduling is written as an integer linear program: binary variables say which step each operation starts in, and linear constraints encode uniqueness, dependences and resources. ILP gives optimal schedules for small blocks but scales poorly. The system-of-difference-constraints (SDC) formulation, from Cong and Zhang, keeps only constraints of the form $s_j - s_i \ge c$, whose LP relaxation is integral, and became the scheduling core of several commercial and open HLS tools.

## The time-indexed formulation

Introduce $x_{i,s} \in \{0, 1\}$, equal to 1 when operation $i$ starts in step $s$. Only steps inside the frame $[\mathrm{ASAP}_i, \mathrm{ALAP}_i]$ are allowed, so operation $i$ contributes $\mu_i + 1$ variables. Six operations with frames of widths 2, 1, 3, 1, 2 and 2 under a 5-step bound need $2 + 1 + 3 + 1 + 2 + 2 = 11$ variables, instead of $6 \times 5 = 30$ for the naive model over all steps. This is why mobility matters for exact methods: the tighter the frames, the smaller the integer program. The start time of operation $i$ is the linear expression $s_i = \sum_s s\, x_{i,s}$, which lets dependences be written on start times while the decision itself stays binary. Tight frames also strengthen the LP relaxation that branch and bound solves at each node, so good bounds and small models go together.

## Constraints and objective

Three families of constraints define a legal schedule. Unique start: $\sum_s x_{i,s} = 1$ for every operation, one constraint each. Precedence: $s_j \ge s_i + l_i$ for every edge $i \to j$, one constraint per edge, written on the start-time expressions. Resources: for every step $s$ and unit type $k$, the operations of type $k$ busy in step $s$ number at most $a_k$, which for single-cycle operations is $\sum_{i \in k} x_{i,s} \le a_k$; that is one constraint per step per type. A block of 7 operations, 11 edges, a 4-step bound and 3 types has $7 + 11 + 12 = 30$ constraints. The objective is either the latency, with the $a_k$ fixed, or the weighted unit count $\sum_k c_k a_k$ with the $a_k$ as integer variables, with the latency fixed. Branch and bound solves it exactly, at exponential worst-case cost.

## The SDC formulation

SDC drops the binary variables and keeps one integer start time $s_i$ per operation. Every constraint is a difference: dependences become $s_j - s_i \ge l_i$, a latency bound becomes $s_{sink} - s_{source} \le L$, and timing constraints for chaining become $s_j - s_i \ge 1$ whenever the combinational delay from the start of $i$ through $j$ exceeds the clock period. On a chain with delays 3, 3, 3 and 1 ns and a 6 ns clock, the pairs (o1, o3), (o1, o4) and (o2, o4) exceed 6 ns and get such constraints. The constraint matrix of a difference system is totally unimodular, so the LP relaxation has integral optimal vertices and an LP solver returns an integer schedule in polynomial time. Resource constraints are not differences; SDC handles them heuristically, by adding ordering constraints between operations that compete for a unit.

## Going deeper

Hwang, Lee and Hsu, "A formal approach to the scheduling problem in high level synthesis" (IEEE TCAD, 1991); Cong and Zhang, "An efficient and versatile scheduling algorithm based on SDC formulation" (DAC, 2006).
