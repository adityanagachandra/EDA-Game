# The recurrence-constrained minimum II

RecMII is the lowest II that a loop's dependence cycles allow, however many resources are available. A value computed in one iteration and needed by a later one forms a cycle in the loop's dependence graph, and the operations on that cycle must complete within the time the iteration distance allows. Recurrences are why some loops refuse to reach II = 1 no matter how much hardware is added.

## One recurrence

Give every dependence edge $u \to v$ a latency $l_u$ and an iteration distance $d$, the number of iterations it spans: 0 within an iteration, 1 from one iteration to the next. In a modulo schedule the edge requires $s_v + d\,\mathrm{II} \ge s_u + l_u$. Summing around a cycle $C$, the start times cancel and leave $\mathrm{II} \cdot \sum_C d \ge \sum_C l$. So each cycle gives $\mathrm{II} \ge \lceil \sum l / \sum d \rceil$. An accumulation whose cycle holds a 12-cycle floating-point adder and has distance 1 needs $\mathrm{II} \ge 12$: iteration $i + 1$ cannot start its addition until iteration $i$'s result exists. With distance 2, as when each iteration updates a value used two iterations later, the bound halves to 6.

## RecMII over all cycles

A loop may have several elementary cycles, and RecMII is the maximum of their bounds. With cycles of latency 4 and distance 2, latency 7 and distance 3, latency 17 and distance 3, and latency 14 and distance 3, the bounds are 2, 3, 6 and 5, so RecMII = 6. Computing it exactly does not require enumerating cycles, which can be exponential in number. Instead, a tool tests a candidate II by checking whether the graph with edge weights $l_u - d\,\mathrm{II}$ has a positive cycle, using a longest-path or Floyd-Warshall pass, and searches for the smallest II without one. That is the same positive-cycle test that decides feasibility of an SDC system, applied to the modulo constraints. The minimum II is then $\max(\mathrm{ResMII}, \mathrm{RecMII})$.

## Reading cycles from a loop body

In a loop body written as a chain of dependent operations $o_1 \to o_2 \to \dots \to o_n$, a loop-carried edge from $o_j$ back to $o_i$, with $i \le j$, closes a cycle through $o_i, \dots, o_j$. Its latency is the sum of those operations' latencies, and its distance is the edge's distance. With latencies 1, 2, 3, 4 and 4, an edge $o_5 \to o_1$ of distance 2 gives $\lceil 14 / 2 \rceil = 7$, and a self-edge $o_3 \to o_3$ of distance 2 gives $\lceil 3 / 2 \rceil = 2$, so RecMII = 7. Memory dependences count too: a store to $a[i]$ followed by a load of $a[i-1]$ in the next iteration is a carried edge of distance 1 through memory, and if the tool cannot prove the addresses differ, it must assume such an edge, which can make RecMII large.

## Going deeper

Rau, "Iterative modulo scheduling: an algorithm for software pipelining loops" (MICRO-27, 1994); Lam, "Software pipelining: an effective scheduling technique for VLIW machines" (PLDI, 1988).
