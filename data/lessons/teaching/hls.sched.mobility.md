# Mobility and the critical path

Mobility, also called slack, is how far an operation can move in time without stretching the schedule: its ALAP start minus its ASAP start. It is the single most used number in HLS scheduling: zero-mobility operations form the critical path, list schedulers use it for priority, force-directed scheduling turns it into probabilities, and ILP models use it to decide how many variables to create.

## ALAP minus ASAP

Compute the ASAP schedule, choose a latency bound $L$, compute ALAP under $L$, and subtract: $\mu_i = \mathrm{ALAP}_i - \mathrm{ASAP}_i$. The time frame $[\mathrm{ASAP}_i, \mathrm{ALAP}_i]$ has $\mu_i + 1$ legal start steps. With additions of 1 step and multiplications of 2, consider $n_1 = a + b$, $n_2 = c \times d$, $n_3 = n_1 + n_2$, with $L$ equal to the critical path of 3 steps. ASAP gives 1, 1, 3. ALAP gives $n_3 = 3$, $n_2 = 3 - 2 = 1$ and $n_1 = 3 - 1 = 2$. So $n_1$ has mobility 1, while $n_2$ and $n_3$ have mobility 0. A negative mobility for any operation would mean the bound $L$ is shorter than the critical path, which is infeasible. Mobility is computed in linear time from the two passes, so tools recompute it after every scheduling decision as frames shrink.

## Zero mobility marks the critical path

When $L$ equals the critical-path length, the operations with zero mobility are exactly those on some longest path. They cannot move at all: delaying any of them by one step delays the end of the schedule. Every other operation has room. Counting zero-mobility operations tells you how constrained a block is. A long chain where every operation is critical leaves the scheduler no freedom, so resource sharing must come from elsewhere; a wide, shallow graph with many mobile operations can be squeezed onto few units. List schedulers exploit this by scheduling low-mobility operations first, since postponing them is what lengthens schedules. The critical path also tells a designer where to apply tree-height reduction or chaining to shorten the latency.

## Relaxing the latency bound

If the bound grows from the critical path $L_0$ to $L_0 + r$, ASAP does not change, because it ignores the bound, and every ALAP start shifts later by exactly $r$, because the backward recursion starts $r$ steps later at every sink and subtracts the same latencies. So every mobility grows by exactly $r$: an operation with mobility 0 at $L_0$ has mobility $r$, and one with mobility 2 has $2 + r$. Relaxing the bound is how a designer trades latency for area: wider frames let the scheduler spread operations of one type across more steps and use fewer units. It also grows every formulation built on frames. A time-indexed ILP has one variable per operation per frame step, so its size grows by $r$ variables per operation.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5; Gajski, Dutt, Wu and Lin, "High-Level Synthesis" (1992), chapter on scheduling.
