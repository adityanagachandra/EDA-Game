# Difference constraints and longest paths

A system of constraints $x_j - x_i \le w_{ij}$ is an LP with a special structure: it is the dual of a shortest-path problem, so it is solved by graph algorithms in place of a general LP solver. EDA uses it constantly: one-dimensional layout compaction, retiming feasibility, clock-skew scheduling and timing-budget distribution. Seeing the constraint graph behind such a system turns an LP into a Bellman-Ford or a longest-path pass.

## Compaction as a longest path

In one-dimensional compaction each layout edge has a position $x_i$, and design rules give minimum spacings $x_j \ge x_i + d_{ij}$. Build a graph with an edge $i \to j$ of weight $d_{ij}$ for each rule. The smallest feasible position of every element is its longest-path distance from the left boundary $x_0 = 0$, and on an acyclic graph that is one pass in topological order.

Example: $x_1 \ge x_0 + 6$, $x_2 \ge x_1 + 1$, $x_3 \ge x_2 + 2$, $x_3 \ge x_1 + 11$, $x_4 \ge x_3 + 3$, $x_4 \ge x_2 + 8$, $x_5 \ge x_4 + 4$. In order: $x_1 = 6$, $x_2 = 7$, $x_3 = \max(9, 17) = 17$, $x_4 = \max(20, 15) = 20$, $x_5 = 24$. The constraint achieving each maximum is the critical one; the chain $x_0, x_1, x_3, x_4, x_5$ fixes the width.

## Feasible solutions by shortest paths

For a general system $x_j - x_i \le w_{ij}$, weights may be negative and the graph may have cycles. Draw an edge $i \to j$ of weight $w_{ij}$ for each constraint and add a source $s$ with a zero-weight edge to every variable. Shortest-path distances $d$ satisfy $d_j \le d_i + w_{ij}$ for every edge, which is exactly the constraint, so $x = d$ is a feasible solution. Bellman-Ford computes it in $O(VE)$.

Example: $x_2 - x_4 \le -3$, $x_1 - x_2 \le 4$, $x_1 - x_4 \le -1$, $x_2 - x_1 \le -3$. Start every $x$ at 0. Relaxing gives $x_2 = -3$ (from $x_4 = 0$), $x_1 = -1$ (from $x_4$), then $x_2 = \min(-3, -1 - 3) = -4$, and nothing improves further: $x = (-1, -4, 0, 0)$. Any constant can be added to every variable and the system still holds.

## Negative cycles and the clock-skew bound

The system is infeasible exactly when the constraint graph has a negative cycle: summing the constraints around the cycle gives $0 \le$ a negative number. This is how clock-skew scheduling finds the best period.

With useful skew, flop $i$ gets clock arrival $s_i$, and a stage of delay $d_i$ into flop $i+1$ needs $s_i + d_i + t_{setup} \le s_{i+1} + T$, that is $s_i - s_{i+1} \le T - d_i - t_{setup}$. Around a loop of $k$ flops the skews cancel and feasibility needs $kT \ge \sum d_i + k\, t_{setup}$. For two flops with stage delays 890 ps and 630 ps and setup 10 ps, $T \ge 1520/2 + 10 = 770$ ps, though the slower stage alone would need 900 ps. Skew can balance a loop but never beat its average delay, so the optimal period is set by the maximum mean cycle of the timing graph.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), section on difference constraints and shortest paths; Fishburn, "Clock skew optimization" (IEEE Transactions on Computers, 1990); Leiserson and Saxe, "Retiming synchronous circuitry" (Algorithmica, 1991).
