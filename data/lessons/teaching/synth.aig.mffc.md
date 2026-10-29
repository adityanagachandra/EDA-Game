# Maximum fanout-free cones

When a node is replaced, how much logic actually disappears? Only the nodes used exclusively by it: its maximum fanout-free cone. The MFFC is the gain side of every AIG transformation, computed in a few microseconds by reference counting, and it separates logic that can be freed from logic shared with the rest of the graph.

## Definition

The **maximum fanout-free cone** of node $n$ is the largest set of nodes in $n$'s transitive fanin such that every path from a node of the set to a primary output passes through $n$. Equivalently, it is $n$ plus every node all of whose fanouts lie inside the cone. If $n$ is removed, exactly the MFFC becomes dead; every other node in its fanin cone is still used elsewhere. In $n_1 = AND(a, b)$, $n_4 = AND(b, d)$, $n_5 = AND(n_1, n_4)$, $n_6 = AND(a, n_5)$, $n_7 = AND(n_1, n_6)$ with output $n_7$, all five AND nodes belong to the MFFC of $n_7$, since $n_1$'s two fanouts are both inside.

## Computing it by dereferencing

Keep a reference count for every node: the number of fanout edges plus one for each primary output that points to it. To measure the MFFC of $n$, dereference it: decrement the count of each fanin, and when a fanin's count reaches zero, recurse into it. The nodes visited are the MFFC. Afterwards, reference the nodes back to restore the counts. The cost is proportional to the MFFC size, so it can be done for every node of a large AIG in a single pass. ABC uses exactly this in rewriting and refactoring to compute how many nodes a replacement would free.

## Shared logic in the cone

The rest of the fanin cone, cone size minus MFFC size, is logic shared with other outputs; it survives any replacement of $n$. A node with a large cone but a small MFFC is cheap to rewrite only if the replacement reuses the shared nodes. In one example the cone of $n_7$ has 4 AND nodes and its MFFC has 3, so 1 node is shared and survives. MFFC also defines the area-flow intuition used in mapping: logic duplicated across fanouts costs more than logic used once. And a transformation that duplicates part of a shared region can increase total area even if the local count improves.

## Going deeper

Cong and Ding, "On area/depth trade-off in LUT-based FPGA technology mapping", IEEE Transactions on VLSI Systems (1994), where MFFCs were used in mapping; Mishchenko, Chatterjee and Brayton, DAC (2006).
