# Heterogeneous circuit graphs

A netlist has different kinds of objects (cells, nets, pins, macros, I/O ports) and different kinds of relations (drives, feeds, belongs to, is near). A homogeneous GNN forces all of them into one node type and one edge type. A heterogeneous GNN gives each type its own features and each relation its own weights, at a parameter cost the designer must budget.

## Relation-specific weights

A relational layer computes $h_v' = \sigma(W_0 h_v + \sum_r \sum_{u \in N_r(v)} W_r h_u / c_{v,r} + b)$, with one matrix $W_r$ per relation $r$ and a normalizing constant $c_{v,r}$ such as the number of $r$-neighbors. With 4 relations (cell-drives-net, net-feeds-cell and their reverses) and width 64, one layer has $(4 + 1) \cdot 64^2 + 64 = 20{,}544$ parameters. Parameters grow linearly with the number of relations, so schemas with many relation types share weights through a basis decomposition, writing each $W_r$ as a combination of a few shared matrices.

## Per-type input encoders

Node types carry different raw features. A cell has a one-hot cell type, drive strength and area; a net has fanout and bounding-box size; a pin has capacitance and direction. A separate linear encoder per type maps each into a shared width $d$, costing $f_t d + d$ parameters for a type with $f_t$ features. With 24 cell features, 6 net features, 4 pin features and $d = 64$, the encoders cost $(24 + 6 + 4) \cdot 64 + 3 \cdot 64 = 2368$ parameters. One-hot cell types tie the model to one library; features computed from Liberty data, such as input capacitance and drive resistance, transfer better across libraries.

## Typed edges and direction

In a heterogeneous netlist graph, a net with one driver and $s$ sinks contributes one drives edge and $s$ feeds edges. Seven nets with fanouts 6, 3, 2, 6, 4, 2 and 3 give 7 drives edges and 26 feeds edges. Adding a reverse relation for each lets messages flow from sinks back to drivers, doubling the total to 66. Direction carries meaning: timing propagates driver to sink and required times flow back, so forward and reverse relations should have separate weights. A homogeneous undirected graph would merge the two and lose that distinction.

## Going deeper

Schlichtkrull et al., "Modeling Relational Data with Graph Convolutional Networks" (ESWC 2018); Hamilton, "Graph Representation Learning" (2020), chapter 5.
