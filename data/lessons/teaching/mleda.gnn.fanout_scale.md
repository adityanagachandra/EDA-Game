# Netlist size and fanout skew in GNNs

Netlists are hard for GNNs for two reasons that most benchmark graphs do not share: they are large, with millions of cells, and their degree distribution is extremely skewed, with clock, reset and scan nets touching thousands of pins. Both shape how a netlist GNN is trained and what it can compute in one pass.

## Messages per layer

In a full-graph layer every edge carries a message in each direction it is used. A bipartite cell-net graph has one edge per pin, so a block with 10,000 signal nets averaging 3 pins plus one 5000-pin clock net has $30{,}000 + 5000 = 35{,}000$ edges and sends 70,000 messages per layer if messages flow both ways. Memory for a layer scales with edges times width, plus stored activations for backpropagation. A design with tens of millions of pins exceeds one GPU's memory for full-graph training, which forces mini-batching. The clock net alone carries one seventh of those messages, from a single net among ten thousand.

## Neighbor sampling

Mini-batch training computes embeddings only for a batch of target nodes, and for each one samples at most $s$ neighbors per layer. The computation tree of one target then has at most $1 + s + s^2 + \dots + s^L$ nodes; with $s = 10$ and $L = 2$ that is 111. Without the cap, a target adjacent to the clock net would pull in all 5000 of its pins at the first hop. Sampling makes training feasible but adds variance, and it changes what the model sees: an aggregate over 10 of 5000 sinks is a noisy estimate of the full aggregate.

## When one net dominates

Under the clique model, a 5000-pin net contributes $\binom{5000}{2} = 12{,}497{,}500$ edges. If the other 10,000 nets add 30,000 edges, the clock net is $12{,}497{,}500/12{,}527{,}500 \approx 99.8$ percent of the graph. Its messages swamp every cell it touches, and every cell becomes two hops from every other. Practical netlist GNNs drop such nets, cap their fanout, star-expand them, or model them separately; clock nets are usually handled by clock-tree synthesis anyway, so their logical connectivity carries little information for tasks such as congestion. The decision is part of the representation, so it should be recorded with the dataset and applied identically at training and inference time.

## Going deeper

Hamilton, Ying and Leskovec, "Inductive Representation Learning on Large Graphs" (NeurIPS 2017), for neighbor sampling; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), on high-fanout nets.
