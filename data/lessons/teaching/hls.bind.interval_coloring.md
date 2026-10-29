# Register binding as interval-graph coloring

Register binding is graph coloring: values are vertices, an edge joins two values that are live at the same time, and registers are colors. For straight-line schedules this interference graph is an interval graph, a perfect graph, so coloring is easy and the minimum number of registers equals the largest clique. Knowing why tells a tool designer when binding stays polynomial and when it becomes NP-hard.

## The interference graph

Draw one vertex per value and an edge between two values whose lifetimes share a step: $[a, b]$ and $[c, d]$ overlap when $a \le d$ and $c \le b$. Two adjacent values can never share a register, and two non-adjacent values can. Count edges by testing every pair: for seven values that is 21 tests. For lifetimes $[4, 5]$, $[3, 3]$, $[2, 4]$, $[4, 6]$, $[3, 3]$, $[5, 6]$ and $[3, 3]$, the three values living only in step 3 all overlap each other and $[2, 4]$, and so on. An interference graph of this kind, built from intervals on a line, is an interval graph. The same structure appears in register allocation for software compilers, where the program's control flow can make it far more general.

## Chromatic number equals the largest clique

A set of intervals that pairwise overlap has a common point, so every clique of an interval graph corresponds to a step where all its values are live, and the largest clique is the maximum overlap. Interval graphs are perfect: for them and every induced subgraph, the chromatic number equals the clique number. So the fewest registers is exactly the maximum number of values live in one step, and a greedy coloring in order of left endpoints, the left-edge algorithm, achieves it in polynomial time. For general graphs, deciding whether three colors suffice is NP-complete, which is why compiler register allocators fall back on heuristics such as Chaitin-style simplification. HLS binding stays easy only as long as lifetimes stay intervals.

## Checking a binding and leaving the interval world

A binding is a coloring, and it is legal when no edge joins two values with the same register. To check a proposed binding, list for each register the values assigned to it and test those pairs for overlap; each overlapping pair is a conflict that would corrupt one value. Conflicts creep in when a binder changes assignments to reduce multiplexers and forgets a lifetime. Lifetimes stop being intervals in loops: a value defined late in one iteration and used early in the next wraps around the loop boundary, and the graph becomes a circular-arc graph, whose chromatic number can exceed its maximum overlap and whose coloring is NP-hard in general. Pipelined loops add modulo lifetimes. Tools then use heuristics, or unroll the kernel so lifetimes become intervals again, at the cost of a larger controller.

## Going deeper

Golumbic, "Algorithmic Graph Theory and Perfect Graphs" (2nd ed., 2004), chapter on interval graphs; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 6.
