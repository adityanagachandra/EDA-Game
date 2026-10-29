# Convergence of an agentic ECO loop

Does an ECO loop finish, and how many iterations will it take? A simple linear model answers surprisingly well. Each iteration fixes a fraction of the open violations and creates some new ones as side effects. The model predicts the iteration count, reveals a floor below which the loop cannot go, and tells an agent when to change strategy instead of iterating.

## A linear model of the loop

Let $v_k$ be the number of violations after iteration $k$. If each iteration fixes a fraction $f$ of the open violations and the fixes create $n$ new ones,
$$v_{k+1} = (1 - f)\,v_k + n.$$
Starting from $v_0 = 400$ with $f = 0.6$ and $n = 20$: $v_1 = 180$, $v_2 = 92$, $v_3 = 56.8$, $v_4 = 42.7$. The new violations might be hold violations created by setup fixes, or DRC violations from inserted buffers. The model is crude, but $f$ and $n$ can be measured from the loop's own history after two or three iterations.

## Iterations to reach a target

The recurrence has the closed form
$$v_k = v^* + (v_0 - v^*)(1 - f)^k, \qquad v^* = \frac{n}{f}.$$
To reach a target $T > v^*$, solve $(1 - f)^k \le (T - v^*)/(v_0 - v^*)$ for the smallest integer $k$. With the numbers above, $v^* = 33.3$; for $T = 50$, $(0.4)^k \le 16.7/366.7 = 0.0455$, so $k \ge \ln 0.0455 / \ln 0.4 = 3.37$, that is 4 iterations, which matches $v_4 = 42.7$. Each iteration shrinks the distance to the floor by the factor $1 - f$, so convergence is geometric.

## The floor and when to change strategy

The fixed point $v^* = n / f$ is a floor: the loop approaches it and never goes below. If the target is below the floor, here any $T < 33.3$, more iterations cannot reach it. The agent must change the parameters of the process: fix more per iteration (raise $f$, for example by grouping violations by cause), or create fewer side effects (lower $n$, for example by checking hold before committing setup fixes). Halving $n$ to 10 lowers the floor to 16.7. Recognizing a floor from two or three iterations of data, instead of iterating until a budget runs out, is what distinguishes a good closure agent.

## Going deeper

Any text on linear first-order recurrences, for example Graham, Knuth and Patashnik, "Concrete Mathematics" (2nd ed., 1994), chapter 1.
