# The timing graph

A static timer never simulates a circuit. It turns the netlist into a directed graph and propagates numbers along it. Every later idea in timing analysis, arrival times, required times, slack, incremental updates, path enumeration, is an operation on this graph. Its shape decides what the timer can compute in one pass and what it must approximate, so a tool designer has to know exactly which pieces of the netlist become nodes and edges.

## Pins as nodes, arcs as edges

Each node is a pin of a cell instance or a top-level port. Edges come in two kinds. A cell arc joins an input pin to an output pin of one instance; its delay and output slew come from that cell's Liberty timing arc, evaluated at the input slew and output load. A net arc joins a net's driver pin to one of its sink pins; its delay comes from the wire's parasitics. Paths begin at startpoints, the input ports and the clock pins of sequential cells, where a clock-to-Q arc launches data. They end at endpoints, the data pins of flip-flops and latches and the output ports. A flip-flop has no propagation arc from D to Q: D carries a setup and hold check against the clock pin, and Q is launched from the clock. Sequential cells therefore cut the graph into combinational pieces.

## Counting cell arcs and net arcs

The graph's size is easy to compute, and it is linear in the netlist. A single-output combinational cell with $k$ inputs contributes $k$ cell arcs, one per input pin. A net with $f$ sinks contributes $f$ net arcs. Take a NAND2 driving three pins, an AOI21 driving two and an inverter driving four. The cells add $2 + 3 + 1 = 6$ cell arcs and the nets add $3 + 2 + 4 = 9$ net arcs, 15 edges in all. Real timers also store separate rise and fall data on each arc and may split one Liberty arc by timing sense, but the count of graph edges still grows with pins, not with paths. That is the property that makes a full-chip timer possible: a design with billions of paths has only millions of arcs.

## Loops, disabled arcs and why a DAG matters

Propagation visits nodes in topological order, so every fanin is final before its fanout is evaluated, and the whole graph is timed in time proportional to its arcs. That needs a directed acyclic graph. Flip-flops already break the sequential feedback loops. A combinational loop, two cross-coupled NAND gates for instance, would make arrival times depend on themselves, so the timer breaks it by disabling one arc and warns. Designers can also remove arcs deliberately with set_disable_timing, for example through a mode-control pin. Every disabled arc deletes paths from the reports, so disabled arcs are reviewed as carefully as timing exceptions. Rise and fall are tracked separately at every node, because the two transitions have different delays and an inverting arc maps one to the other.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapters 1 and 7; the OpenSTA source, whose Graph and Levelize classes build and order this graph.
