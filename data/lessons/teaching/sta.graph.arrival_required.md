# Arrival and required times

Arrival times say when data settles; required times say when it must settle. A timer computes the first forward from the startpoints and the second backward from the endpoints, and their difference is slack. The backward pass mirrors the forward one, with a minimum instead of a maximum, and it is what turns a single endpoint constraint into a deadline at every internal pin, the deadline an optimizer reads when it decides which gate to resize.

## Required time at an endpoint

At a flip-flop's data pin, setup requires the data to settle a setup time before the capturing clock edge reaches the flop. With period $T$, capture clock latency $L_{cap}$ and setup time $t_{su}$, the required time is $RAT = T + L_{cap} - t_{su}$, measured from the launching source edge at 0. For a 1000 ps clock, a capture latency of 60 ps and a 40 ps setup, the data must arrive by $1000 + 60 - 40 = 1020$ ps. At an output port, the required time is the capturing edge minus the output delay from the SDC. With ideal clocks the latency term is zero, and the endpoint required time is simply $T - t_{su}$. Uncertainty, when present, is subtracted as well.

## Backward propagation with a minimum

For an internal pin $u$, every fanout arc $u \to v$ imposes a deadline: the signal must leave $u$ by $RAT(v) - d(u,v)$ to reach $v$ in time. All fanouts must be satisfied, so $RAT(u) = \min_{u \to v} \big(RAT(v) - d(u,v)\big)$. Visit nodes in reverse topological order. Suppose $u$ feeds endpoint $X$ through 20 ps with $RAT(X) = 150$ ps and endpoint $Y$ through 35 ps with $RAT(Y) = 120$ ps. The deadlines are 130 and 85 ps, so $RAT(u) = 85$ ps. One backward pass, linear in the arcs, gives every pin its deadline without listing paths.

## Why the minimum governs

The minimum is the mirror of the forward maximum. A pin can feed many endpoints, and the tightest one, the endpoint with the least room after its own path delay, sets the deadline for everything upstream. A startpoint deep in the graph inherits the strictest constraint among all endpoints it reaches. That is also why a single badly constrained endpoint, say an output delay set far too large, drags down the required times of every pin in its fanin cone and makes the optimizer chase phantom violations. For hold, the backward pass flips again: the hold requirement propagates with a maximum, since data must not arrive earlier than the latest hold deadline among the fanouts.

## Going deeper

Sapatnekar, "Timing" (2004), chapter on static timing analysis; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8.
