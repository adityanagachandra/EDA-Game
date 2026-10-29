# Resource lower bounds on schedule length

Before running a scheduler, a tool can bound how short any schedule could be from two simple counts: the critical path and the work each unit type must do. These bounds are instant to compute, they tell the designer whether a latency target is even possible with a given allocation, and branch-and-bound schedulers use them to prune.

## Counting work per unit type

Suppose a block has $n_k$ operations of type $k$, each occupying a unit for one step, and the allocation provides $a_k$ units of that type. Each step serves at most $a_k$ of them, so the operations of type $k$ alone need at least $\lceil n_k / a_k \rceil$ steps, whatever the dependences. Twenty-three multiplications on four multipliers need at least $\lceil 23 / 4 \rceil = 6$ steps. This bound ignores the graph completely, and that is its strength: it holds for every schedule, so a latency target below it is infeasible with this allocation, and the tool can say so before scheduling. Reading it the other way, a target of $L$ steps needs at least $\lceil n_k / L \rceil$ units of type $k$, which is the starting point of allocation.

## Multicycle operations occupy units longer

If an operation of type $k$ takes $l_k$ steps on a non-pipelined unit, the unit is busy for all of them, so the type needs $n_k l_k$ unit-steps of work. With $a_k$ units supplying $a_k$ unit-steps per step, the bound becomes $\lceil n_k l_k / a_k \rceil$. Ten two-step multiplications on three multipliers need at least $\lceil 20 / 3 \rceil = 7$ steps. A fully pipelined unit accepts a new operation every step, so its occupancy is one step and the original bound returns, although the last results still arrive $l_k - 1$ steps later. This is the first place where the difference between latency and occupancy matters; it reappears in modulo scheduling, where occupancy rather than latency determines how many operations fit in one initiation interval.

## Combining the bounds

A schedule must respect every bound at once, so the best simple lower bound is the maximum of the critical-path length and the per-type bounds: $L \ge \max(L_{cp}, \lceil n_{add} / a_{add} \rceil, \lceil n_{mul} / a_{mul} \rceil)$. For a block with a critical path of 4 steps, 7 additions on 2 adders and 5 multiplications on 1 multiplier, the bound is $\max(4, 4, 5) = 5$. A real scheduler may need more, because dependences and resources interact: operations may be ready only in steps where the units are already busy. Exact schedulers such as branch and bound use sharper versions of these bounds on partial schedules to prune; list schedulers simply compare their result against this bound to judge how far from optimal they might be.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5, section on resource-constrained scheduling; Hu, "Parallel sequencing and assembly line problems" (Operations Research, 1961).
