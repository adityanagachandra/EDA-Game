# Netlists as graphs for learning

A netlist is a hypergraph: a net joins any number of pins. Graph neural networks operate on ordinary graphs, so the first modeling decision is how to turn nets into edges. The choice decides how large the graph is, whether high-fanout nets swamp it, and whether the net itself survives as an object the model can attach features to.

## Clique, star and bipartite encodings

The clique model joins every pair of a net's pins, so a $p$-pin net becomes $\binom{p}{2}$ edges: a 5-pin net gives 10. The star model adds one node for the net and joins it to each pin, $p$ edges. The bipartite cell-net graph is the star model applied everywhere: one node per cell, one node per net, one edge per pin. Each preserves different information. The clique model loses the net as an entity and treats a 2-pin and a 50-pin connection between the same cells alike. The star and bipartite models keep the net, so net features such as fanout or bounding box have somewhere to live.

## Counting edges in a netlist

Edge counts add over nets. A netlist with nets of 3, 2, 3, 2 and 16 pins has clique edges $3 + 1 + 3 + 1 + 120 = 128$; the one 16-pin net contributes 120 of them. Under the bipartite model the same netlist has $3 + 2 + 3 + 2 + 16 = 26$ edges. Real designs have clock, reset and scan-enable nets with thousands of pins. A 5000-pin clock net alone has $\binom{5000}{2} \approx 12.5$ million clique edges, more than all the signal nets of a large block. Counting each net's edges separately gives a multigraph, since two cells can share several nets.

## Choosing an encoding

The ratio of clique edges to bipartite edges, $\sum \binom{p}{2} / \sum p$, measures how much the clique model inflates the graph; for nets of 2, 10, 2, 5 and 2 pins it is $58/21 \approx 2.76$, and it grows without bound with fanout. Most netlist GNNs therefore use a bipartite or star encoding, or keep cliques only for small nets and drop or star-expand large ones. Direction matters too: signals flow from driver to sinks, and timing tasks usually keep that direction, sometimes with separate driver and sink edge types. The encoding is a representation decision like any other in EDA: it fixes what the model can see and what it costs to compute.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 2, on net models; Hamilton, "Graph Representation Learning" (2020), chapters 5 and 7.
