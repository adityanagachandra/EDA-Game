# ASAP and ALAP scheduling

ASAP and ALAP are the two extreme schedules of a data-flow graph with unlimited resources: every operation as early as its inputs allow, or as late as the latency bound allows. They take linear time, they give the minimum schedule length, and the window between them, each operation's time frame, is the input to list scheduling, force-directed scheduling and every ILP or SDC formulation.

## As soon as possible

Number control steps from 1. An operation that starts in step $s$ with latency $l$ occupies steps $s$ to $s + l - 1$, and a consumer may start in step $s + l$. ASAP visits the operations in topological order: a source starts in step 1, and any other operation starts at $\max_p (s_p + l_p)$ over its predecessors $p$. With additions of 1 step and multiplications of 2, take $n_1 = a + d$, $n_2 = d \times d$, $n_3 = n_1 + n_2$. Then $n_1$ and $n_2$ start in step 1, and $n_3$ waits for the later of $1 + 1 = 2$ and $1 + 2 = 3$, so it starts in step 3. This is a longest-path computation on the DAG with edge weights equal to the producer latencies, done in one pass over the graph.

## The minimum schedule length

The ASAP schedule finishes when its last operation finishes, in step $\max_i (s_i + l_i - 1)$. No schedule can be shorter, because each operation already starts at the earliest moment its dependences allow; this number is the latency-weighted critical path of the DFG. Continuing the example with $n_4 = n_3 + g$, $n_4$ starts in step 4 and ends in step 4, so the block needs at least 4 steps. Unlimited hardware is assumed: ASAP may put five multiplications in the same step. The schedule length is therefore a lower bound that real, resource-limited schedules can only match or exceed, and it is the natural latency target to start from. A designer who needs fewer steps must change the graph, by tree-height reduction or chaining, not the scheduler.

## As late as possible

Given a latency bound of $L$ steps, ALAP schedules backwards. A sink starts at $L - l + 1$, so that it finishes exactly in step $L$; any other operation starts at $\min_c (s_c) - l_i$ over its consumers $c$. Visit the operations in reverse topological order. With $L$ equal to the critical path, operations on the critical path get the same start in ASAP and ALAP; with a looser bound, every operation can slide later. The pair gives each operation a time frame $[\mathrm{ASAP}_i, \mathrm{ALAP}_i]$ of legal start steps, assuming unlimited resources. ALAP also has a practical use: placing operations late shortens the lifetimes of values that feed them from far away, and it is the classic seed for list-scheduling priorities.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5, sections on unconstrained and latency-constrained scheduling.
