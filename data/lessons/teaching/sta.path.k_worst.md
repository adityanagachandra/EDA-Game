# Enumerating the k worst paths

Block-based propagation gives each endpoint its worst arrival, but reports, path-based re-timing and optimizers need actual paths, and often more than one per endpoint. Enumerating the $k$ worst paths in delay order, without listing all paths, is a classic graph problem that every timer solves. OpenSTA's report_checks with -group_count and -endpoint_count is built on it.

## The worst path to an endpoint

The worst path ends with the latest arrival at the endpoint, so its delay is exactly the endpoint's arrival from the forward pass, measured from the source. To recover the path itself, walk backward: at each node choose the fanin $u$ with $AT(u) + d(u,v) = AT(v)$, the fanin that produced the maximum. For arcs $P \to W$ 23, $P \to Q$ 10, $Q \to W$ 29, $W \to B$ 11, $B \to F$ 26 and $Q \to F$ 27 ps, the arrival at $F$ is 76 ps, along $P \to Q \to W \to B \to F$. Storing which fanin won at each node makes this backtrace linear in the path length.

## The second-worst path

A second path to the same endpoint must deviate from the worst one somewhere: at some node it takes a fanin other than the winning one. The cost of a deviation at node $v$ through fanin $u$ is $AT(v) - (AT(u) + d(u,v))$, how much later the winner was. The second-worst path is the worst path with the cheapest single deviation, then completed by the worst path into the deviating fanin. In the example, deviating at $W$ to the fanin $P$ costs $39 - (0 + 23) = 16$, giving a 60 ps path, while deviating at $F$ to the fanin $Q$ costs $76 - (10 + 27) = 39$, giving 37 ps. So the delays in order are 76, 60 and 37 ps. Paths that tie count as separate paths.

## Priority-queue enumeration and reporting

Generalizing gives an efficient algorithm: keep a priority queue of partial paths ordered by their best possible completed delay, which is known exactly from the arrivals. Pop the best, extend it backward toward the source, and push its deviations. Each popped complete path is the next worst, so $k$ paths cost roughly $O(k)$ queue operations times the path length, after one propagation. Reports then ask for paths within some $\delta$ of the worst, or the worst few per endpoint: for delays 68, 58 and 46 ps, two paths lie within 15 ps of the worst. Because reconvergent logic yields many near-duplicates, tools limit paths per endpoint to expose distinct problems.

## Going deeper

Eppstein, "Finding the k shortest paths" (SIAM Journal on Computing, 1998); the OpenSTA source, PathEnum.
