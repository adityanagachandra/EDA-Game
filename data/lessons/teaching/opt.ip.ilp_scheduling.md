# ILP formulation of operation scheduling

High-level synthesis must assign operations of a dataflow graph to control steps so that dependences hold and the number of functional units stays within a limit. The textbook exact method is a time-indexed integer program. Its size, and therefore its solve time, is governed by each operation's mobility window, and its resource constraints decide how many adders and multipliers the datapath needs. The same template schedules tests, tasks on an accelerator, and cuts in a lithography flow.

## Time-indexed binaries and their count

Let $x_{i,t} = 1$ when operation i starts in step t. Each operation starts exactly once: $\sum_t x_{i,t} = 1$. Its start time is $\sum_t t\, x_{i,t}$, so a dependence $i \to j$ with delay $d_i$ is linear: $\sum_t t\, x_{j,t} \ge \sum_t t\, x_{i,t} + d_i$.

Variables only need to exist where an operation can actually start, between its ASAP and ALAP times. With windows op1 [3, 5], op2 [3, 8], op3 [7, 8], op4 [7, 8], op5 [2, 3] and op6 [3, 4], the counts are 3, 6, 2, 2, 2 and 2, total 17 binaries, against $6 \times 8 = 48$ without pruning. Window pruning is the cheapest and most effective model reduction in ILP scheduling, and tighter latency bounds shrink windows further.

## ASAP, ALAP and mobility

ASAP times come from a forward longest-path pass: sources start in step 1, and each operation starts one step (its delay) after its latest predecessor. ALAP times come from a backward pass under the latency bound L: sinks start at L, and each operation must start before its earliest successor. Mobility is ALAP minus ASAP.

Example with single-cycle operations and dependences op1 to op2, op1 to op3, op2 to op4, op1 to op5, op4 to op5, op5 to op6, and $L = 5$. ASAP: op1 = 1, op2 = 2, op3 = 2, op4 = 3, op5 = 4, op6 = 5. ALAP: op6 = 5, op5 = 4, op4 = 3, op2 = 2, op1 = 1, and op3 has no successor so it may wait until 5. Op4 has mobility 0: it is on the critical path. Op3 has mobility 3 and four candidate steps.

## Resource constraints and unit counts

Resource limits are linear in the binaries. An operation of type r with delay $d_i$ is busy in step t if it started in any of the $d_i$ steps up to t, so the constraint for each step and type is $\sum_{i \in r} \sum_{\tau = t - d_i + 1}^{t} x_{i,\tau} \le a_r$.

Given a schedule, the smallest $a_r$ it satisfies is the peak number of busy units. Multiplier mul1 starts in step 3 for 3 steps (busy in 3, 4, 5), mul2 starts in 2 for 2 steps (2, 3), mul3 starts in 2 for 1 step (2), and mul4 starts in 4 for 3 steps (4, 5, 6). Busy counts: step 2 has 2, step 3 has 2, step 4 has 2, step 5 has 2, step 6 has 1, so 2 multipliers suffice. Making $a_r$ a variable with a cost turns the ILP into the area-latency trade-off that HLS tools explore, and list scheduling is the fast heuristic that approximates it.

## Going deeper

Hwang, Lee and Hsu, "A formal approach to the scheduling problem in high level synthesis" (IEEE TCAD, 1991); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5.
