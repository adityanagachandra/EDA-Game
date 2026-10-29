# Deferred-merge embedding

Deferred-merge embedding (DME) is the algorithm that turns a clock topology into a zero-skew tree with minimum wire. Instead of fixing each merge point as it goes, it records every valid location as a merging segment on the way up, then chooses actual points on the way down, where each choice can be made relative to an already placed parent. Under the linear delay model DME is optimal for the given topology, and with Elmore delay it is the basis of most zero-skew routers.

## The bottom-up phase

Process the topology in post-order. A sink's segment is its own location. For an internal node with children segments $ms_a$ and $ms_b$, compute their Manhattan distance $d$, then the edge lengths from the zero-skew condition, $e_a = (d + t_b - t_a)/2$ and $e_b = d - e_a$ (or Tsay's formula under Elmore delay). The node's merging segment is the set of points within $e_a$ of $ms_a$ and $e_b$ of $ms_b$, the intersection of two tilted rectangular regions, which is again a Manhattan arc. For sinks at $(50, 0)$ with $t_a = 25$ and $(20, 20)$ with $t_b = 35$, $d = 50$ and $e_a = 30$ µm. If $|t_a - t_b| > d$, one edge needs detour wire.

## The top-down phase

Process the topology in pre-order. Place the root at any point of its segment, or at the point nearest the clock source. For every child, place it at the point of its own segment nearest its placed parent. Because the parent's segment lies at exactly distance $e$ from the child's segment, that nearest point is at distance $e$, so the planned edge length is realised with no detour. For a child segment from $(45, 20)$ to $(35, 30)$ and a parent at $(-10, 30)$, the nearest point is $(35, 30)$, 45 µm away. Each step is a constant-size geometric computation, so both phases take linear time.

## What deferral buys

Choosing a tap point early, as a greedy bottom-up router does, commits to one location before knowing where the parent will be, and later edges pay for the mismatch. DME keeps all equally good locations until the parent is fixed, which is why it achieves the minimum total wirelength among zero-skew trees for that topology under path-length delay. It does not choose the topology, which remains a heuristic problem, and under Elmore delay or with buffers its optimality is no longer exact. Bounded-skew DME generalises segments to regions, trading a skew budget for less wire.

## Going deeper

Chao, Hsu, Ho, Boese and Kahng, "Zero skew clock routing with minimum wirelength", IEEE Transactions on Circuits and Systems II (1992); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
