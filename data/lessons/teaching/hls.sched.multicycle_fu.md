# Multicycle and pipelined functional units

Many operators take several cycles: multipliers, dividers, floating-point units, memory reads. Whether such a unit is pipelined, accepting a new operation every cycle, or not, blocking until it finishes, changes how many operations it serves and therefore the whole schedule. The distinction between latency, the cycles to a result, and occupancy, the cycles the unit is busy, is the key to resource bounds and modulo scheduling.

## Non-pipelined units

A non-pipelined unit with latency $L$ accepts an operation, computes for $L$ cycles and only then accepts the next. With $u$ such units and $n$ independent operations all ready in step 1, the units work in rounds of $L$ cycles, each round finishing $u$ operations, so the last result is available at the end of step $\lceil n / u \rceil L$. Thirteen multiplications on three 6-cycle multipliers take 5 rounds, finishing in step 30. Occupancy equals latency here, so the resource bound on schedule length is $\lceil n L / u \rceil$, and every operation ties up its unit for the full duration. Non-pipelined iterative units, such as a sequential divider that produces one quotient bit per cycle, are small but slow; their area advantage is what justifies the blocking.

## Fully pipelined units

A pipelined unit splits its logic into $L$ stages separated by registers and accepts a new operation every cycle, so its occupancy is one cycle while its latency stays $L$. With $u$ units and $n$ independent operations, the last issue happens in step $\lceil n / u \rceil$ and its result appears $L - 1$ steps later, in step $\lceil n / u \rceil + L - 1$. Ten multiplications on two 6-cycle pipelined multipliers issue in steps 1 to 5 and finish in step 10, against step 30 for the non-pipelined case. For independent work the pipeline hides latency almost completely. For a dependent chain it hides nothing: each operation still waits $L$ cycles for its input. That asymmetry is why recurrences, not resources, often limit pipelined loops.

## Partially pipelined units

Between the two extremes sits a unit with an initiation interval $p$: it accepts a new operation every $p$ cycles, so occupancy is $p$ and latency is $L$, with $1 \le p \le L$. A divider with $L = 4$ and $p = 2$ issuing nine divisions on three units issues in steps 1, 3 and 5 and finishes the last in step $5 + 4 - 1 = 8$. In general the completion step is $(\lceil n / u \rceil - 1) p + L$. Floating-point units and memory controllers often look like this, and HLS operator libraries record both numbers per operator. Schedulers must use occupancy for resource conflicts and latency for dependences; confusing the two either overbooks a unit or wastes cycles waiting for results that are already available.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5; Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", the appendix on pipelining with multicycle operations.
