# AIG rewriting

AIG rewriting replaces the logic of a small cut with a better precomputed structure for the same function. DAG-aware rewriting, introduced in ABC in 2006, counts only nodes that a replacement truly frees and truly adds, so it improves the whole graph rather than one tree at a time. It is fast enough to run over millions of nodes many times, and it is the workhorse of ABC scripts.

## Gain of a replacement

For node $n$, enumerate its 4-feasible cuts. For each cut, compute the 4-input truth table of $n$ over the leaves. Look up precomputed AIG structures implementing that function. A replacement frees the MFFC of $n$ (nodes used only by $n$) and adds the nodes of the new structure that do not already exist in the graph; existing ones are found by structural hashing. The **gain** is MFFC size minus newly added nodes. With an MFFC of 5 and a candidate of 5 nodes of which 3 already exist, the gain is $5 - (5 - 3) = 3$. Positive gain means the graph shrinks.

## Precomputed structures and NPN classes

There are $2^{16} = 65536$ functions of four inputs, but they fall into 222 classes under negation and permutation of inputs and negation of the output (NPN equivalence). Rewriting precomputes good AIG structures for the classes that matter, transforms a cut's function to its class representative, and maps the stored structure back by applying the same permutation and negations. For each cut, several candidate structures are tried, and the best gain over all cuts and candidates is chosen. With an MFFC of 7 and candidates adding 3, 2 and 2 new nodes, the best gain is 5. Lookup is a table access, so the whole step costs microseconds per node.

## Zero-cost moves and the overall pass

The pass visits nodes in topological order and accepts a replacement when the gain is positive. The variant rwz also accepts gain 0. Zero-cost rewriting does not shrink the graph, but it changes its structure, which can expose positive gains to later passes; scripts such as resyn2 use rw, then rwz, then other passes. Rewriting is local and greedy, so results depend on order and are not optimal, but it preserves function exactly, because each replacement implements the same truth table over the same leaves. Level-preserving options reject replacements that would increase depth, trading area gain for delay.

## Going deeper

Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting: a fresh look at combinational logic synthesis", DAC (2006); Bjesse and Boralv, "DAG-aware circuit compression for formal verification", ICCAD (2004).
