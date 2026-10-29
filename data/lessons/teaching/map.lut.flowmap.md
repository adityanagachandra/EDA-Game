# FlowMap

FlowMap computes, for every node, the minimum depth at which it can be implemented in a network of k-input LUTs, in polynomial time. It labels nodes in topological order and decides each label with one max-flow computation. Its theorem, that delay-optimal LUT mapping of a DAG is polynomial, underlies every modern FPGA mapper, even those that use cut enumeration instead of flow.

## Labels as optimal depth

With unit LUT delay, the **label** l(v) is the minimum depth of v over all k-LUT mappings of its fanin cone; primary inputs have label 0. Equivalently, l(v) is the minimum over k-feasible cuts C of v of 1 + max of the leaves' labels. Labels never decrease along edges. Example with k = 3: n1 has fanins b, c and n2 has fanins a, d, so both get label 1. n3 with fanins d, n2 can absorb n2 using cut {a, d}, label 1. n4 with fanins n1, n3 needs leaves {b, c, a, d} to absorb both, which is 4 > 3, so its best cut is {n1, n3} and its label is 2.

## The flow test: p or p plus one

Let p be the largest label among v's fanins. Then l(v) is either p or p + 1. To decide, collapse v together with every node of label p in its fanin cone into a single sink, give every other node capacity 1, connect a source to the primary inputs, and compute a maximum flow. By max-flow min-cut the flow equals the size of the smallest node cut separating the inputs from the sink. If that size is at most k, the cut's nodes all have label below p, so a LUT over them gives v label p; otherwise l(v) = p + 1, achieved by a LUT whose inputs are v's own fanins, which is allowed because the network is k-bounded. In the example above, p = 1 for n4 and the minimum cut after collapsing has 4 nodes, so l(n4) = 2.

## From labels to a mapping

After labeling, FlowMap generates LUTs backward from the outputs: each needed node is implemented over its minimum-height cut, and that cut's nodes become needed. Every output then sits at depth equal to its label, the optimum. Node splitting turns the node-capacity problem into an ordinary edge flow, and since the flow value is at most k + 1 when it matters, each augmenting-path search can stop early; the total is polynomial, $O(k m n)$ for n nodes and m edges. The mapping is depth-optimal but not area-aware and may duplicate logic heavily, so area recovery passes follow, as in later depth-relaxing and cut-based mappers.

## Going deeper

Cong and Ding, "FlowMap: an optimal technology mapping algorithm for delay optimization in lookup-table based FPGA designs", IEEE TCAD 1994; Ahuja, Magnanti and Orlin, "Network Flows", for max-flow with node capacities.
