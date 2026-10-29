# Formulating an optimization problem

Every EDA engine is an optimizer wearing a domain costume. A placer, a sizer, a router and a scheduler each choose values for some unknowns to make a cost small while respecting rules. Writing that down precisely, as decision variables, an objective and constraints, is the first act of tool design. A sloppy formulation produces a fast solver for the wrong problem, and most QoR surprises trace back to what was, or was not, written into the model.

## Decision variables versus data

Decision variables are the quantities the optimizer is free to set. Everything else is data. In gate sizing the variables are the size of each cell; the library, the netlist and the clock period are data. In global placement the variables are the $(x, y)$ coordinates of movable cells; pad positions and the die outline are data. Derived quantities such as arrival times are neither: they are functions of the variables computed by the timer.

The encoding matters as much as the choice. An HLS scheduler can use one integer start time per operation, or binaries $x_{i,t} = 1$ when operation $i$ starts in cycle $t$. The binary, time-indexed form makes per-cycle resource limits linear, $\sum_i x_{i,t} \le a$, which an ILP solver can use. A macro that may rotate gets one binary that swaps its width and height. Pick variables so that the rules become simple functions of them.

## Objective versus constraints

The objective is what we trade; the constraints are what we must meet. A buffering tool must meet required times and wants little buffer area, so it minimizes area subject to timing at every endpoint. Swapping the roles, minimizing slack subject to an area cap, is a different problem with a different answer.

A constraint can be made soft by moving it into the objective with a weight: a placer that minimizes $W(x) + \lambda D(x)$ treats density $D$ as a penalty rather than a hard bin capacity. Then the optimizer may buy wirelength with overflow, and a router scoring wirelength plus 1000 times overflow can legitimately return 3 units of overflow if that saves enough wire. A worst-case objective such as worst negative slack, $\max_p \ell_p$, is written with an auxiliary variable: minimize $t$ subject to $t \ge \ell_p$ for every path $p$. That epigraph trick keeps a linear model linear.

## Proxies, scaling and infeasibility

Most EDA objectives are proxies. Placers minimize half-perimeter wirelength (HPWL) although the real metric is routed length; sizers use Elmore or table delays although signoff uses a full timer. The gap between proxy and signoff is a design decision to be measured, not ignored.

Scaling is the second trap. If delay is about 200 ps and area about 10,000 square microns, the objective $\text{delay} + \text{area}$ is almost all area, and timing is ignored. Normalize each term, for example by its baseline value, before choosing weights.

Infeasibility is a property of the constraints alone. If a clock-tree model has no feasible point, look for conflicting constraints, such as a skew bound tighter than wire delays allow; the objective plays no part. Finally, two formulations with the same optimum can have very different relaxations, and the tighter one solves faster.

## Going deeper

Boyd and Vandenberghe, "Convex Optimization" (2004), chapter 4 on problem forms and equivalent reformulations; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), which states each physical-design step as an optimization problem.
