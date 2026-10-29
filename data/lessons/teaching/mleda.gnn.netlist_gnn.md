# Graph neural networks on netlists

Stacking message-passing layers gives a GNN, which maps a netlist to embeddings per cell, net or pin, and from them to predictions such as slack, congestion or a partition label. Three properties decide whether a GNN fits an EDA task: how far information travels, whether the path it must travel is short enough, and whether the model applies to netlists it has never seen.

## The receptive field

After $k$ layers, a node's embedding depends exactly on the nodes within $k$ hops: its $k$-hop neighborhood. A one-layer GNN sees direct neighbors, a two-layer GNN their neighbors as well. Computing this is a breadth-first search to depth $k$. If a gate has three neighbors in the cell graph and each of them has two further neighbors, a two-layer embedding of the gate depends on $1 + 3 + 6 = 10$ cells. The receptive field is a hard limit: if the label depends on a node outside it, no amount of training can find that dependence, and the model learns a correlation instead.

## Depth versus logic depth

Timing is a long-range property. A path from a launching flip-flop through 9 gates to a capturing flip-flop has 11 cells and 10 nets. In a bipartite cell-net graph each net costs two hops, cell to net and net to cell, so the capturing flop is 20 hops from the launching one, and a GNN needs 20 layers before one can influence the other. Deep GNNs over-smooth and are slow, so timing models usually do something else: they propagate along the timing graph in topological order, one level per step, so that arrival-like features accumulate along paths as they do in STA.

## Inductive models and parameter count

A GNN's weights are shared by every node, so its parameter count does not depend on the netlist size. A model with 16 input features, width 64, three layers each with a self matrix, a neighbor matrix and a bias, and a linear readout has $(2 \cdot 16 \cdot 64 + 64) + 2(2 \cdot 64^2 + 64) + 65 = 18{,}689$ parameters, whether it runs on a 10-thousand-cell or a 10-million-cell netlist. That makes the model inductive: it applies to unseen netlists. It does not make it accurate on them. Generalization depends on whether the new netlist's neighborhoods look like the training ones.

## Going deeper

Hamilton, "Graph Representation Learning" (2020); Guo et al., "A Timing Engine Inspired Graph Neural Network Model for Pre-Routing Slack Prediction" (DAC 2022).
