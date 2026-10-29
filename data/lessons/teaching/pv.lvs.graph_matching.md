# LVS as graph matching

LVS asks whether the layout netlist and the schematic netlist describe the same circuit. Both are graphs of devices and nets, so the question is graph isomorphism with labels: device types, terminal kinds and properties. General isomorphism has no known polynomial algorithm, but circuit graphs are strongly labelled, and practical matchers solve them in near-linear time by partition refinement.

## The device-net bipartite graph

Model a netlist as a bipartite graph: one vertex per device, one vertex per net, and one edge per device terminal (gate, source, drain; bulk is often a fourth). An inverter has 2 devices and 4 nets (A, Y, VDD, GND), so 6 vertices and 6 edges. A NOR2 has 4 devices and 6 nets (A, B, Y, the internal node, VDD, GND): 10 vertices and 12 edges. Edges carry the terminal kind, and source and drain share one kind because they are interchangeable. This graph is the same structure as the netlist hypergraph, with each net expanded into a vertex.

## Isomorphism and invariants

Two netlists match if a bijection maps devices to devices of the same type and nets to nets so that every terminal edge is preserved with its label. Any isomorphism preserves invariants: device counts per type, the number of nets of each degree, and so on. If the layout has four nets with three terminals and the schematic five, they cannot match. Invariants prove mismatches cheaply but never prove a match. Refinement uses increasingly fine invariants, local neighbourhood signatures, to split vertices into classes that must map onto each other.

## From mismatch to diagnosis

A matcher refines both netlists together with shared colours, so each class can be compared across graphs. Classes with equal sizes are matched or ambiguous; a class with three layout nets but two schematic nets proves an error and points to where it is. Because colours depend on neighbourhoods, one short perturbs everything within a few hops, and the report can list hundreds of unmatched nets. Tools seed the partition with matched port names to anchor it, compare hierarchically by black-boxing verified cells, and report the first conflicting classes so a designer can find the root cause.

## Going deeper

Ebeling, "GeminiII: A Second Generation Layout Validation Program" (ICCAD, 1988); the Netgen documentation.
