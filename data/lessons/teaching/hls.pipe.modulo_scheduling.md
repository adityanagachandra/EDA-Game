# Modulo scheduling

Modulo scheduling builds one schedule for a single iteration such that the same schedule, started every II cycles, never overuses a resource or violates a dependence. It is how HLS tools pipeline loops, inherited from software pipelining for VLIW processors. The key structure is the modulo reservation table, which folds time modulo II so that resource conflicts between overlapping iterations become conflicts in one table.

## The modulo reservation table

When iterations start every II cycles, an operation scheduled at time $t$ of its iteration executes at $t$, $t + \mathrm{II}$, $t + 2\,\mathrm{II}$ and so on in absolute time, once per iteration. In steady state, all operations whose times share the same residue $t \bmod \mathrm{II}$ execute in the same cycle, from different iterations. The modulo reservation table has II rows, and each operation reserves its resource in row $t \bmod \mathrm{II}$. A schedule is resource-legal when no row of the table uses more copies of a resource than exist. With II = 5 and multiplications at times 8, 9, 9, 10 and 12, the rows are 3, 4, 4, 0 and 2, so row 4 holds two multiplications and the schedule needs two multipliers, or one multiplication must move.

## Iterative modulo scheduling

Rau's iterative modulo scheduling starts with $\mathrm{II} = \max(\mathrm{ResMII}, \mathrm{RecMII})$. It schedules operations in priority order. For each one it computes the earliest start allowed by already-scheduled predecessors, including carried edges through the term $-d\,\mathrm{II}$, and then tries times $t, t + 1, \dots, t + \mathrm{II} - 1$ until it finds a free row in the table; trying more than II consecutive times is pointless, because the rows repeat. With II = 5, one multiplier, rows 0 occupied and an earliest start of 10, time 10 falls in row 0 and is blocked, and time 11 falls in free row 1, so the operation starts at 11. If no row is free, the algorithm evicts a conflicting operation and reschedules it later, within a budget; if the budget runs out, it increases II and starts again.

## Stages, prologue and epilogue

The finished schedule of one iteration has some length $SL$; cutting it into blocks of II cycles gives $S = \lceil SL / \mathrm{II} \rceil$ stages. In the steady-state kernel, stage $k$ of iteration $i$ runs alongside stage $k - 1$ of iteration $i + 1$, and so on, so all $S$ stages are busy. Before that, a prologue of $(S - 1)\,\mathrm{II}$ cycles fills the pipeline one stage at a time, and after the last iteration starts, an epilogue of the same length drains it. With II = 3 and $SL = 10$, there are 4 stages, and the prologue lasts 9 cycles. In HLS the prologue and epilogue are usually realized with stage-valid bits rather than separate code. Short loops spend a large fraction of their time in fill and drain, which is why tools offer to flatten nested loops into one long pipelined loop.

## Going deeper

Rau, "Iterative modulo scheduling: an algorithm for software pipelining loops" (MICRO-27, 1994); Lam, "Software pipelining: an effective scheduling technique for VLIW machines" (PLDI, 1988); Zhang and Liu, "SDC-based modulo scheduling for pipeline synthesis" (ICCAD, 2013).
