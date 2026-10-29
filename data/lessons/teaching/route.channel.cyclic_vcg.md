# Cyclic vertical constraints

If net a must be above net b in one column and b above a in another, the vertical constraint graph has a cycle and no assignment of one trunk per net can satisfy both. Detecting cycles, and knowing how many nets must be split to break them, is the decision a channel router makes before it can place a single trunk.

## Detecting cycles

A net lies on a directed cycle of the VCG exactly when it can reach itself: run a search from each net, or compute strongly connected components and take every component with more than one net. If any net lies on a cycle, a dogleg-free channel routing is impossible, because the nets of the cycle would each have to be strictly above the next.

Example: top pins 3, 1, 5, 5 and bottom pins 5, 3, 4, 4 give edges $3 \to 5$, $1 \to 3$ and $5 \to 4$; if some column also gives $5 \to 1$ or $4 \to 3$, a cycle appears. Counting the nets on cycles tells you how much of the channel is affected.

## How many nets must be split

Splitting a net into subnets (a dogleg) can break a cycle, but only cycles that pass through a split net. So the set of nets a router splits must meet every cycle: it is a feedback vertex set of the VCG. The size of the smallest feedback vertex set is therefore a lower bound on the number of nets to split.

Example: a VCG whose cycles all pass through net 5 needs only one split. Finding a minimum feedback vertex set is NP-hard in general, but channel VCGs are small enough to search. The bound is not a guarantee: a split helps only if the net has a pin column between the two conflicting columns, where the dogleg can jump tracks.

## Shortest cycles

The shortest directed cycle tells you how tight the conflict is. A 2-cycle arises when two nets swap top and bottom positions in two columns, the most common case. Breadth-first search from each net for a path back to itself finds the shortest cycle in $O(V(V + E))$.

A short cycle between two nets whose pin columns interleave is easy to break with one dogleg. A long cycle through many nets suggests a pin arrangement that will cost several doglegs and extra tracks. Historically, channels with unbreakable cycles required changing the pin order or adding a column; in modern area routing the same situation shows up as a local knot of pins that only a reroute through another layer resolves.

## Going deeper

Sherwani, "Algorithms for VLSI Physical Design Automation", on channel routing with cyclic constraints; D. N. Deutsch, "A Dogleg Channel Router" (1976).
