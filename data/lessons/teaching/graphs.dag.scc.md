# Strongly connected components and combinational loops

A combinational loop, a cycle that passes through no register, breaks static timing, makes simulation oscillate or depend on order, and is almost always a design error. Finding loops is a strongly-connected-components problem: every vertex on a cycle lies in a component with at least two vertices. Lint tools, synthesis front ends and timers all run this check, and they must report the loop in a form a designer can act on.

## Components by mutual reachability

Two vertices of a directed graph are strongly connected if each can reach the other. This is an equivalence relation, and its classes are the **strongly connected components**. A vertex on no cycle is a component by itself. Contracting each component to a single vertex gives the condensation, which is always a DAG.

Arcs $v_1 \to v_2$, $v_2 \to v_1$, $v_1 \to v_6$, $v_6 \to v_3$, $v_6 \to v_4$, $v_4 \to v_3$, $v_6 \to v_5$, $v_7 \to v_2$ give components $\{v_1, v_2\}$, $\{v_3\}$, $\{v_4\}$, $\{v_5\}$, $\{v_6\}$, $\{v_7\}$: six components, one of them a 2-cycle. The condensation is acyclic, so a timer could still levelize everything except the loop.

## Linear-time algorithms

Tarjan's algorithm finds all components in one depth-first search. It gives each vertex a discovery index and a lowlink, the smallest index reachable from its subtree through at most one back arc, and keeps visited vertices on a stack. When a vertex's lowlink equals its own index, it is the root of a component, and the component is popped from the stack.

Kosaraju's algorithm uses two searches: one on the graph to get finishing times, then one on the reversed graph in decreasing finishing time, where each search tree is a component. Both run in $O(V + E)$. On a netlist with millions of pins this is a fraction of a second, which is why loop checks run on every read of a design.

## From components to loop reports

Every vertex in a component of size two or more lies on a cycle, given no self-loops. So the number of vertices to flag is the total size of all nontrivial components. A self-loop, an output fed directly to its own input, makes a single vertex a loop and must be checked separately.

A component is not yet a useful report: it may hold hundreds of gates. Tools extract one short cycle inside it (a breadth-first search from one member back to itself) and report that path. To proceed with timing anyway, a timer breaks loops by disabling chosen arcs, which is a feedback arc set problem; minimum feedback arc set is NP-hard, so timers break arcs heuristically and warn. Some loops are intentional, such as latch-based or asynchronous structures, and designers mark them so the tool does not break them blindly.

## Going deeper

Tarjan, "Depth-first search and linear graph algorithms" (SIAM J. Computing, 1972); the OpenSTA documentation on combinational loops.
