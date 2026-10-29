# Message passing and aggregation

A graph neural network layer updates each node from its own state and an aggregate of its neighbors' states. The aggregation is the part that handles graphs of any shape: it must accept any number of neighbors in any order. Sum, mean, max and degree-normalized sums each make different assumptions about what a netlist node's neighborhood means.

## Sum aggregation

One message-passing layer computes $h_v' = \sigma(W_s h_v + W_n \sum_{u \in N(v)} h_u)$. Take scalar features and node $v$ with $h_v = 1$ and neighbors with features 2, -1 and 5. The sum is 6. With $W_s = 1$, $W_n = 0.5$ and ReLU, $h_v' = \mathrm{ReLU}(1 + 3) = 4$. Sum aggregation is permutation-invariant, so the order in which neighbors are listed does not matter, and it is the most expressive of the simple aggregators because it can count. It also scales with degree, so a cell on a 1000-pin net receives a message a thousand times larger than a cell on a 2-pin net.

## Mean and max aggregation

Mean aggregation divides the sum by the degree: the same node has mean $6/3 = 2$, and $h_v' = \mathrm{ReLU}(1 + 0.5 \cdot 2) = 2$. Mean is insensitive to degree, which stabilizes training on netlists with skewed fanout, but it cannot tell a node with one neighbor of feature 2 from one with three such neighbors. Max takes the largest neighbor feature, 5 here, and suits questions such as the latest arriving input of a gate, which is a max in static timing. A negative pre-activation is clipped by ReLU to zero, so the sign of the weights matters when checking a hand computation.

## Degree-normalized GCN layers

The graph convolutional network adds self-loops and normalizes each message by $\sqrt{\tilde d_u \tilde d_v}$, where $\tilde d$ is the degree plus one: $h_v' = W \sum_{u \in N(v) \cup \{v\}} h_u/\sqrt{\tilde d_u \tilde d_v}$. If $v$ has degree 3 ($\tilde d_v = 4$) and its neighbors have degrees 1, 2 and 3 with features 2, -1 and 5, the sum is $1/4 + 2/\sqrt{8} - 1/\sqrt{12} + 5/\sqrt{16} \approx 1.918$ for $h_v = 1$. Messages from high-degree neighbors are damped, which softens the effect of large nets. Stacking such layers repeatedly averages features, and after many layers node states become nearly indistinguishable, a failure called over-smoothing.

## Going deeper

Kipf and Welling, "Semi-Supervised Classification with Graph Convolutional Networks" (ICLR 2017); Hamilton, Ying and Leskovec, "Inductive Representation Learning on Large Graphs" (NeurIPS 2017).
