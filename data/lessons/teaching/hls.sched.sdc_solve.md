# Solving SDC schedules by longest paths

A system of difference constraints is a graph problem in disguise. Each constraint $s_j - s_i \ge w$ is a weighted edge, the least feasible schedule is the longest-path distance from a source, and infeasibility is a positive cycle. Seeing SDC this way lets a tool designer solve small systems by hand, explain why a latency constraint fails, and understand why SDC scheduling is fast.

## Constraints as a weighted graph

Make one vertex per start variable. A constraint $s_j - s_i \ge w$ becomes an edge $i \to j$ with weight $w$. With all starts at least 0, the smallest feasible value of each $s_j$ is the longest-path distance to $j$, starting every vertex at 0: $s_j = \max(0, \max_{i \to j} (s_i + w_{ij}))$. For forward constraints only, the graph is a DAG and one pass in topological order solves it. Take $s_2 - s_1 \ge 2$, $s_3 - s_1 \ge 1$, $s_4 - s_2 \ge 1$ and $s_4 - s_3 \ge 3$. Then $s_1 = 0$, $s_2 = 2$, $s_3 = 1$ and $s_4 = \max(2 + 1, 1 + 3) = 4$. These are exactly the ASAP starts when the weights are latencies, so ASAP is the simplest SDC solution. The same relaxation, run backwards from a latency bound, gives the ALAP starts as the latest feasible values.

## Upper bounds become backward edges

A maximum-latency constraint $s_j - s_i \le c$ is rewritten as $s_i - s_j \ge -c$: an edge from $j$ back to $i$ with weight $-c$. It says $s_i \ge s_j - c$, so if other constraints push $s_j$ late, $s_i$ is dragged later too, to stay within $c$ steps of it. In the example, adding $s_4 - s_2 \le 1$ gives $s_2 \ge s_4 - 1 = 3$, so $s_2$ rises from 2 to 3 and the system is re-relaxed. The graph now has cycles, so a single topological pass no longer suffices; Bellman-Ford relaxation, repeated until nothing changes, finds the least solution in $O(VE)$ time. Such constraints model throughput requirements, I/O timing relative to a start signal, and the maximum separation between two operations that share a register.

## Feasibility and positive cycles

A difference system is feasible exactly when its constraint graph has no cycle of positive total weight. Around a cycle the constraints add up to $0 \ge \sum w$, which fails if the sum is positive; Bellman-Ford detects this as relaxation that never stops. For a latency constraint $s_j - s_i \le c$, the cycle is the longest forward path from $i$ to $j$, of weight $P$, closed by the backward edge of weight $-c$, so the system is feasible if and only if $c \ge P$. With $P = 4$ from $s_1$ to $s_4$, the smallest allowed $c$ is 4: the designer's latency target is infeasible below the critical path between those two operations. This is how an SDC scheduler explains an infeasible constraint: it reports the positive cycle, which names the chain of operations that cannot be squeezed.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms", the section on difference constraints and shortest paths; Cong and Zhang, "An efficient and versatile scheduling algorithm based on SDC formulation" (DAC, 2006).
