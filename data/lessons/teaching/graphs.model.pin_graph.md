# Directed pin graphs of drivers and sinks

Timing analysis, levelization and loop detection need to know which way signals flow. A net has one driver pin and some sink pins, so the netlist becomes a directed graph with an arc from each driver to each of its sinks. This representation is the substrate of static timing, logic simulation and incremental analysis, and its in-degrees, out-degrees and reachable cones are the quantities those engines manipulate.

## From nets to driver-to-sink arcs

Each net $e$ with $k$ pins has exactly one driver (an output pin, or a primary input) and $k-1$ sinks (input pins, or primary outputs). The **pin graph** has a vertex per pin, or per cell when each cell has one output, and an arc from the driver to each sink. So a $k$-pin net contributes $k - 1$ arcs: this is the star model with the driver as centre and a direction attached.

Nets $n_1$ driven by $c_1$ with sinks $c_3, c_4$, and $n_2$ driven by $c_2$ with sink $c_4$ give arcs $c_1 \to c_3$, $c_1 \to c_4$, $c_2 \to c_4$: three arcs from nets of sizes 3 and 2. Inside a cell, timing tools add internal arcs from each input pin to each output pin it affects, which is how a cell-level graph becomes a pin-level timing graph.

## Fan-in, fanout and cones

The out-degree of a driver is its **fanout**, the number of sinks it must charge. The in-degree of a cell is its **fan-in**, the number of its inputs that are driven. The **transitive fanout cone** of a vertex is everything reachable from it along arcs; the transitive fan-in cone is everything that reaches it.

Cones bound the effect of a change: if a cell's delay changes, only arrival times in its fanout cone can change, and only required times in its fan-in cone. For arcs $c_2 \to c_3$, $c_2 \to c_6$, $c_3 \to c_8$, the fanout cone of $c_2$ is $\{c_3, c_6, c_8\}$: three cells, found by depth-first search in time linear in the arcs visited. Logic synthesis uses fan-in cones the same way, as the region a rewrite of one output can touch.

## Direction as a source of structure

Undirected models cannot express that a signal flows one way. The directed pin graph of a combinational netlist is acyclic, which is what lets timers and simulators process vertices in topological order. Sequential elements break the cycles: a timer cuts the graph at flip-flop data pins and starts new paths at their clock-to-output arcs.

Direction also changes the net models a tool should use. A partitioner treats all pins of a net symmetrically, but a timing-driven placer weights the driver-to-critical-sink connections more than others. And a bidirectional pin, or a tristate bus with several drivers, does not fit the one-driver assumption; tools either model it specially or report it, which is one reason designs avoid internal tristates.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8 on timing graphs; the OpenSTA documentation on its timing graph.
