# Latency-constrained scheduling and unit counts

The dual of resource-constrained scheduling fixes the latency and asks for the fewest units. This is the usual situation when a throughput or deadline requirement is given and area is to be minimized. Every schedule implies a unit count, the busiest step's demand, and different legal schedules within the same latency imply very different areas.

## The unit count a schedule implies

Given a fixed schedule with single-cycle operations, the datapath needs as many units of type $k$ as the largest number of type-$k$ operations sharing one step: $a_k = \max_s |\{ i : \mathrm{type}(i) = k, s_i = s \}|$. In an ASAP schedule that places three multiplications in step 1 and one in each later step, three multipliers are needed even though most of them idle after step 1. ASAP tends to bunch operations early, because everything ready starts at once, so it usually implies more units than necessary. For multicycle non-pipelined operations, count the operations whose busy interval covers the step instead of only those starting in it. This count is what binding later realizes: it is the minimum number of cliques in the compatibility graph of that type.

## The counting lower bound under a latency target

Whatever the dependences, $n_k$ single-cycle operations of type $k$ in $L$ steps need at least $\lceil n_k / L \rceil$ units, since $a$ units serve at most $aL$ operations. Fifteen multiplications in 9 steps need at least $\lceil 15 / 9 \rceil = 2$ multipliers. The bound is achieved only if the operations can be spread evenly, which dependences may prevent: if all fifteen sit in two dependent layers, at least eight must share a step. Still, the bound is the right first estimate of area for a latency target, and it is what a designer uses to size the allocation before running a scheduler. A larger $L$ lowers it, which is the basic latency-area trade-off of HLS.

## Where operations sit inside their frames

Moving operations within their time frames changes the implied unit count without changing the latency. Placing every operation at its ALAP step bunches work at the end, just as ASAP bunches it at the start; for a block whose multiplications have frames spread over several steps, both extremes can need two multipliers where a balanced placement needs one. Latency-constrained schedulers therefore look for a placement that flattens each type's per-step demand. Force-directed scheduling does it with distribution graphs, ILP formulations minimize the unit count directly subject to the latency bound, and simple iterative methods repeatedly move an operation from the busiest step to a less busy step inside its frame. The question in every case is the same: which step within its frame should each operation take so that the peak demand is lowest.

## Going deeper

Paulin and Knight, "Force-directed scheduling for the behavioral synthesis of ASICs" (IEEE TCAD, 1989); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5.
