# The W and D matrices

To ask whether a retiming can reach clock period $c$, Leiserson and Saxe summarize every pair of nodes by two numbers: the fewest registers on any path between them, W, and the longest delay among the paths that achieve that minimum, D. The period constraints are then simple difference constraints on the lags, which a shortest-path algorithm can solve.

## W: the fewest registers between two nodes

$W(u, v)$ is the minimum, over all paths from $u$ to $v$, of the number of registers on the path. Because retiming changes a path's register count by $r(v) - r(u)$ regardless of the route, the path with the fewest registers before retiming is still the one with the fewest after. For a ring $v_1 \to v_2 \to v_3 \to v_4 \to v_1$ with weights 2, 2, 1, 0 plus a chord $v_2 \to v_1$ of weight 0, $W(v_4, v_1) = 0$ through the direct edge. $W$ is an all-pairs shortest-path computation on the register weights, so Floyd-Warshall or repeated Dijkstra computes it.

## D: the critical delay at that register count

$D(u, v)$ is the maximum total delay, including $d(u)$ and $d(v)$, over the paths from $u$ to $v$ that have exactly $W(u, v)$ registers. Those paths are the ones that will become register-free together if a retiming removes all their registers. Computing $W$ and $D$ together uses lexicographic weights $(w(e), -d(u))$ in a shortest-path algorithm. For a five-node example with $W(v_1, v_2) = 2$, the paths from $v_1$ to $v_2$ with two registers have delays up to 11 ns, so $D(v_1, v_2) = 11$. $D(v, v)$ is just $d(v)$, a reminder that the period can never be below the largest node delay.

## Period constraints

Leiserson and Saxe proved that a retiming $r$ achieves period at most $c$ if and only if it is legal, $r(u) - r(v) \le w(e)$ for every edge $e: u \to v$, and $r(u) - r(v) \le W(u, v) - 1$ for every pair with $D(u, v) > c$. The second family forces at least one register onto every slow path. Both families are difference constraints, solvable (or shown infeasible) by Bellman-Ford on a constraint graph. Testing $c = 6$ ns in a five-node circuit means counting the ordered pairs with $D(u, v) > 6$, here 19 of the 20, each adding one constraint. Minimum-period retiming binary-searches $c$ over the sorted distinct values of $D$.

## Going deeper

Leiserson and Saxe, "Retiming synchronous circuitry", Algorithmica (1991); De Micheli (1994), the chapter on sequential logic optimization.
