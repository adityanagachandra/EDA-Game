# Maximum-slack clock scheduling as a linear program

Feasibility is not enough: a schedule that meets every constraint with zero margin breaks under the first bit of variation. Clock scheduling is therefore posed as a linear program that maximises the minimum slack over all setup and hold checks. The decision variables are the clock arrival times and the margin; the constraints are the difference constraints with the margin subtracted. The LP's structure makes it solvable at scale, and its solution becomes the latency targets that CTS must implement.

## The formulation

Variables: an arrival time $t_i$ for each flop and a margin $M$. Objective: maximise $M$. Constraints: for each path $i \to j$, $t_i - t_j \le T - S_{ij} - M$ for setup, with $S_{ij} = t_{cq} + d_{max} + t_{su}$, and $t_j - t_i \le H_{ij} - M$ for hold, with $H_{ij} = t_{cq} + d_{min} - t_h$. For fixed $M$ these are difference constraints, so the LP is a parametric shortest-path problem; for a ring of $n$ flops with setup only, summing gives $M^\ast = T - \sum S/n$. With stage requirements summing to 1314 ps over two flops and $T = 700$ ps, $M^\ast = 700 - 657 = 43$ ps.

## A two-flop loop with hold

Let A and B launch into each other and set $u = t_A - t_B$. Setup A to B and hold B to A both bound $u$ from above: $u \le P - M$ with $P = \min(T - S_{AB}, H_{BA})$. Setup B to A and hold A to B bound it from below: $-u \le Q - M$ with $Q = \min(T - S_{BA}, H_{AB})$. Adding the two gives $M \le (P + Q)/2$, reached at $u = (P - Q)/2$. With $P = -18$ and $Q = 96$ ps, the best margin is 39 ps at $t_A - t_B = -57$ ps: B's clock should arrive 57 ps after A's.

## From schedule to tree

The LP's output is a target arrival for every flop, typically expressed as offsets from a common latency. CTS then builds a tree whose sink latencies hit those targets, for example by giving sinks different required delays in the merge, or by inserting delay on chosen branches. Two cautions apply. The targets are only as good as the path delays used to compute them, and those delays change after CTS and routing, so scheduling is iterated. And large intentional skews add latency on some branches, which costs power and adds variation. Production tools therefore limit how far a schedule may move any flop.

## Going deeper

Fishburn, "Clock skew optimization", IEEE Transactions on Computers (1990); Bertsimas and Tsitsiklis, "Introduction to Linear Optimization" (1997), for the LP background.
