# Clock scheduling as difference constraints

Every setup and hold check between two flops bounds the difference of their clock arrival times. A whole design's clock schedule is therefore a system of difference constraints, which has a beautiful graph structure: feasibility is the absence of a negative cycle, and a feasible schedule comes from shortest paths. This is the formulation behind clock-scheduling engines and the reason useful skew can be computed exactly on large designs.

## Constraints as a graph

For a path from flop $i$ to flop $j$, setup gives $t_i - t_j \le T - t_{cq} - d_{max} - t_{su}$ and hold gives $t_j - t_i \le t_{cq} + d_{min} - t_h$. Each constraint $t_v - t_u \le w$ becomes an edge from $u$ to $v$ of weight $w$. Summing constraints around any cycle cancels every $t$, so a cycle of negative total weight is a contradiction. For three flops in a ring with longest delays 1070, 720 and 810 ps, $t_{cq} = 70$ and $t_{su} = 34$, the cycle weight at $T = 930$ ps is $3 \cdot 930 - 2912 = -122$ ps: no schedule meets setup.

## Shortest paths give a schedule

Add a virtual source with a zero-weight edge to every node and run Bellman-Ford. If no negative cycle exists, the distances $t_v = \mathrm{dist}(v)$ satisfy every constraint, because a shortest path obeys $\mathrm{dist}(v) \le \mathrm{dist}(u) + w$, which is exactly $t_v - t_u \le w$. For constraints $t_D - t_C \le 5$, $t_C - t_D \le 30$, $t_D - t_A \le -15$ and $t_C - t_A \le 20$, relaxation gives $t = (0, 0, 0, -15)$ for A to D. If an edge can still be relaxed after $n - 1$ rounds, there is a negative cycle and the period is infeasible.

## The minimum period

With setup constraints only, every cycle $C$ of flops needs $|C| \cdot T \ge \sum_{C}(t_{cq} + d_{max} + t_{su})$, so the minimum period is the maximum over cycles of the mean stage delay. For a ring with delays 780, 640 and 1100 ps, $t_{cq} = 45$ and $t_{su} = 36$, it is $2763/3 = 921$ ps, below the 1181 ps the slowest stage would need with zero skew. Skew redistributes time around a loop but cannot create it. Hold constraints add edges of the opposite direction, so in general the minimum period is found by searching on $T$ with a negative-cycle test, or as a maximum cycle ratio.

## Going deeper

Fishburn, "Clock skew optimization", IEEE Transactions on Computers (1990); Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms", the section on difference constraints and shortest paths.
