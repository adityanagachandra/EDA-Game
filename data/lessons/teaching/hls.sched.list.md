# List scheduling

List scheduling is the workhorse heuristic for resource-constrained scheduling: walk the control steps in order and, in each, fill the free units with the most urgent ready operations. Resource-constrained scheduling is NP-hard in general, but list scheduling runs in near-linear time, is easy to extend, and is close to optimal on most real blocks, which is why nearly every HLS tool contains a version of it.

## Priority by distance to the sink

A list scheduler needs a priority that says which ready operation matters most. The classic choice is the height of an operation: the length of the longest path from it to any sink, counted in steps and including the operation itself. With unit latencies, compute it in reverse topological order: a sink has height 1, and every other operation has one plus the largest height among its consumers. An operation with a large height starts a long chain, so delaying it is likely to delay the whole schedule. Height is closely related to mobility: under a latency bound equal to the critical path $L$, the ALAP start of an operation is $L - h + 1$ when latencies are one, so sorting by largest height is the same as sorting by earliest ALAP. Other priorities, such as smallest mobility or most successors, are common variants.

## Stepping through the schedule

Start at step 1 with nothing scheduled. In each step, the ready list holds the unscheduled operations whose predecessors have all finished. Sort it by priority, ties to the lower-numbered operation, and for each unit type assign ready operations of that type to free units until the units or the candidates run out. Then advance one step. Take $n_1 = a + b$, $n_2 = c + d$, $n_3 = n_1 + n_2$, $n_4 = n_1 \times n_2$, $n_5 = n_4 + e$ with one adder and one multiplier, unit latencies. Heights are 3, 3, 1, 2, 1. Step 1 schedules $n_1$; step 2 schedules $n_2$; step 3 has $n_3$ and $n_4$ ready, and they use different units, so both go; step 4 runs $n_5$. The schedule takes 4 steps, against a critical path of 3, because the single adder serialized $n_1$ and $n_2$.

## Multicycle units and the quality of the result

With multicycle operations, the scheduler tracks which units are still busy. A non-pipelined multiplier with a latency of 2 started in step 3 is unavailable in step 4, and its consumers become ready in step 5. The algorithm is otherwise unchanged. List scheduling is greedy: it never delays a ready operation to leave room for a more important one that becomes ready a step later, so it can miss the optimum. Hu's algorithm, list scheduling by height, is optimal for unit-latency trees on identical units, but general graphs with several unit types have no such guarantee. In practice tools compare the result against the lower bound from the critical path and the resource counts; if they match, the schedule is provably optimal, and if not, the gap tells how much an exact method could win.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5, section on list scheduling; Hu, "Parallel sequencing and assembly line problems" (Operations Research, 1961).
