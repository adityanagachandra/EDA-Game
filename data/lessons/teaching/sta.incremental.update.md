# Incremental timing analysis

An optimizer changes a design thousands of times, resizing a cell, inserting a buffer, moving an instance, and after each change it needs new slacks. Re-timing the whole design each time would be far too slow. Incremental timing recomputes only what a change can affect, and knowing exactly what that is, a forward cone for arrivals and a backward cone for required times, is the core of making a timer usable inside optimization loops.

## Arrivals change in the forward cone

Suppose only the delay of arc $u \to v$ changes. Arrival times flow forward, so the arrival at $v$ may change, and so may the arrival at every node reachable from $v$; nothing else can, since no other node's arrival depends on that arc. The set to re-evaluate is the fanout cone of $v$, including $v$. In a graph with arcs $M \to N$, $N \to B$, $M \to H$, $B \to U$, $E \to N$, $H \to P$, $B \to P$, $H \to B$ and $E \to M$, a change on $H \to B$ affects arrivals at $B$, $U$ and $P$: three nodes. A resize changes several arcs at once, the driven arcs of the cell and the arcs into it through its input capacitance, so the cone starts from all of them.

## Required times change in the backward cone

Required times flow backward, so the change on $u \to v$ can alter the required time at $u$ and at every node that reaches $u$: the fanin cone of $u$, including $u$. For the change on $H \to B$ above, that is $H$, $M$ and $E$. The two cones are independent: an arrival never depends on downstream arcs and a required time never depends on upstream ones. Endpoint required times change only if the clocks or constraints change, which is why a constraint edit is far more expensive to update than a netlist edit: it can touch the required times of the entire design.

## Slack changes in the union, lazily

Slack is required minus arrival, so it can change wherever either does: in the union of the forward cone of $v$ and the backward cone of $u$. In the example, that is $E$, $M$, $H$, $B$, $P$ and $U$, six of seven nodes. Timers like OpenSTA avoid even this work until it is needed: an edit marks the affected arrivals and required times invalid, and the next query recomputes only what that query needs, processing invalid nodes in level order. If the optimizer then asks for the slack of one endpoint, only its fanin cone is brought up to date. Many edits between queries are batched into a single update.

## Going deeper

The OpenSTA source and documentation on incremental timing; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8.
