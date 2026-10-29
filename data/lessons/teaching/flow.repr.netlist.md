# Netlists as the shared currency of the flow

A netlist is the list of component instances in a design and the nets that connect their pins. It is the representation almost every stage reads or writes: synthesis produces it, placement and routing attach geometry to it, timing and power analyze it, LVS compares against it. Mathematically a netlist is a hypergraph, and the way a tool turns that hypergraph into a graph decides what its algorithms can do.

## Instances, pins and nets

Each instance references a master cell from a library, for example U17 of type NAND2. Each master has pins with directions: inputs, outputs, inouts. A net connects a set of pins: usually one driver (an output pin or a primary input) and one or more sinks. The fanout of a net is its number of sink pins. A netlist file says nothing about where cells are or how they behave inside; the master names point into Liberty and LEF for that. Netlists are usually hierarchical, with module instances inside modules, and tools flatten them on demand.

## The netlist as a hypergraph

Model each instance as a vertex and each net as a hyperedge, the set of instances whose pins it connects. A net with $p$ pins is one hyperedge of size $p$. Counting is the first computation: the total pin count is the sum of hyperedge sizes, and the number of nets is the number of hyperedges. If a netlist has nets with 2, 3, 2, 5 and 4 pins, it has 5 nets and 16 pin connections, and the 5-pin net has a fanout of 4. Partitioners such as hMETIS and the Fiduccia-Mattheyses algorithm work on the hypergraph directly, counting how many hyperedges a cut splits.

## Graph models of a hypergraph

Many algorithms need an ordinary graph. The clique model replaces a $p$-pin net with all $p(p - 1)/2$ pairwise edges, usually weighted by $1/(p - 1)$ or similar so that large nets do not dominate. The star model adds one new vertex per net and $p$ edges to it. A 5-pin net becomes 10 edges as a clique or 5 as a star. Quadratic placers use these models to write wirelength as a quadratic form, and the choice changes both the result and the problem size: a single 1000-pin net becomes almost half a million clique edges, which is why large nets are given star models or ignored in early placement.

## Going deeper

Alpert and Kahng, "Recent Directions in Netlist Partitioning: A Survey", Integration, the VLSI Journal (1995); Kahng, Lienig, Markov and Hu, "VLSI Physical Design", chapter 2.
