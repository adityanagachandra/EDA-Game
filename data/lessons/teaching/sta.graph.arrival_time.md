# Forward propagation of arrival times

The arrival time at a pin is when the signal there settles, measured from the launching clock edge. Block-based static timing computes it for every pin in one forward pass over the timing graph, with a maximum at every merge. This is the core numerical kernel of a timer: it must be exact, fast and incremental, and everything else, required times, slack and optimization, is built on top of it.

## The max-plus recurrence

A startpoint gets its arrival from the clock: a flip-flop's Q pin arrives at the clock arrival plus clock-to-Q, and an input port at its input delay. Every other pin takes the latest of its fanins: $AT(v) = \max_{u \to v} \big(AT(u) + d(u,v)\big)$. The maximum is the right operator for setup, because a gate's output is not settled until its last input has settled. Take a source $S$ at 5 ps with arcs $S \to A$ (17 ps), $S \to B$ (16 ps) and $A \to B$ (8 ps). Then $AT(A) = 22$ ps, and $B$ sees $5 + 16 = 21$ and $22 + 8 = 30$, so $AT(B) = 30$ ps. The arrival is a property of the pin, not of one path, which is why the recurrence touches each arc once.

## Propagating to the endpoints

Visit the nodes in topological or level order, so that every fanin is final first, and apply the recurrence at each node. The latest arrival at an endpoint is then the delay of the longest path into it, found without listing a single path. This is longest path in a DAG, which is easy even though longest path in a general graph is hard. The cost is linear in the number of arcs. For a chain with a side branch, $S$ at 0 ps, $S \to A$ 20 ps, $A \to C$ 15 ps, $S \to C$ 30 ps and $C \to E$ 10 ps, the arrivals are $A = 20$, $C = \max(35, 30) = 35$ and $E = 45$ ps. A real timer stores this per transition (rise and fall) and per clock domain, but the recurrence is the same.

## Early arrivals for hold

Hold checks ask the opposite question: how soon can new data arrive? The same pass with a minimum computes the earliest arrival, $AT_{min}(v) = \min_{u \to v} \big(AT_{min}(u) + d_{min}(u,v)\big)$, using minimum delays from the fast corner or early derates. In the example above the earliest arrival at $C$ is $\min(35, 30) = 30$ ps, so $E$ is at 40 ps early and 45 ps late. Timers propagate both values together in one traversal, storing a min and a max arrival per pin and transition. Using a max where a min is needed would make hold analysis optimistic, which is the most dangerous kind of error in a timer: it reports a passing check on silicon that fails.

## Going deeper

Sapatnekar, "Timing" (2004), chapter on static timing analysis; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
