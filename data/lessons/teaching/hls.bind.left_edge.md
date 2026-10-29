# Register binding by the left-edge algorithm

The left-edge algorithm assigns values with interval lifetimes to registers using the minimum possible number of registers. It was invented for channel routing, assigning wire segments to tracks, and Kurdahi and Parker carried it into register allocation for HLS. It is simple, runs in $O(n \log n)$, and is optimal for intervals, which makes it the standard register binder in straight-line code.

## Sort by left edge and fill registers

Sort the values by the start of their lifetime, the left edge, breaking ties by a fixed order such as value number. Process them in that order. Give each value the lowest-numbered register whose last value ended before this one starts; if no register is free, open a new one. For lifetimes $v_1: [1, 1]$, $v_6: [1, 2]$, $v_5: [1, 2]$, $v_7: [2, 3]$, $v_2: [3, 5]$, $v_3: [4, 4]$ and $v_4: [5, 5]$, sorted order is $v_1, v_5, v_6, v_7, v_2, v_3, v_4$. Then $v_1 \to R_1$, $v_5 \to R_2$, $v_6 \to R_3$ (both earlier registers are busy in step 1), $v_7 \to R_1$, $v_2 \to R_2$, $v_3 \to R_1$ and $v_4 \to R_1$. Three registers are used, which matches the three values live together in step 1. The only data structure needed is the list of each register's last end time, so the pass is a sort followed by a linear scan.

## Why it is optimal

Suppose the algorithm opens register $R_k$ for a value $v$ starting in step $t$. Every earlier register is busy at $t$, holding a value that started at or before $t$ (sorting guarantees that) and ends at or after $t$. So $k$ values, including $v$, are live in step $t$, and at least $k$ registers are necessary. The number of registers used therefore equals the maximum number of overlapping lifetimes, the lower bound. Equivalently, left edge is a greedy coloring of an interval graph in order of left endpoints, and such a coloring uses exactly the clique number of colors. Filling one register at a time, taking from the sorted list every value that fits after the previous one, yields the same assignment, which is the form in which the channel-routing literature states it.

## Why the sort matters

The sort is the whole trick. First-fit in an arbitrary order can use more registers than the maximum overlap. Take $a: [1, 2]$, $b: [4, 5]$, $c: [3, 6]$ and $d: [2, 3]$, processed in that listed order. First-fit puts $a$ and then $b$ in $R_1$; $c$ overlaps $b$, so it opens $R_2$; $d$ overlaps $a$ in $R_1$ and $c$ in $R_2$, so it opens $R_3$. Yet no step has more than two live values. Sorted by left edge, the order is $a, d, c, b$, giving $a \to R_1$, $d \to R_2$, $c \to R_1$ and $b \to R_2$: two registers. The sort prevents the failure because no value processed later can start before one already placed, so whenever every register is busy at a value's start, each is busy with a value live in that very step. Left edge minimizes registers only; it ignores which values share a register, which decides multiplexer cost, so practical binders weight the choice of register to reduce interconnect.

## Going deeper

Hashimoto and Stevens, "Wire routing by optimizing channel assignment within large apertures" (DAC, 1971); Kurdahi and Parker, "REAL: a program for register allocation" (DAC, 1987).
