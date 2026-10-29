# Layer assignment and via minimization

Many global routers route on a 2D projection and only afterwards assign each wire segment to a metal layer. Layer assignment decides via count, and vias matter: each adds resistance, blocks tracks on the layers it passes through and is a yield risk. Counting the vias of an assignment, and knowing when it is minimal, is the starting point for every assignment algorithm.

## Counting vias of a path

In a preferred-direction stack each segment of a route sits on a layer of its own direction. A via stack between layers $a$ and $b$ counts $|a - b|$ vias. For a two-pin route with pins on metal 1 and segments on layers $l_1, \dots, l_n$, the via count is
$$(l_1 - 1) + \sum_{i=2}^{n} |l_i - l_{i-1}| + (l_n - 1).$$

Example: with M2, M4 and M6 horizontal and M3, M5 and M7 vertical, a route H on M4, V on M5, H on M4, V on M7 needs $3 + 1 + 1 + 3 + 6 = 14$ vias: 3 down to the first pin, one at each of the first two bends, 3 at the last bend and 6 down to the second pin. Long via stacks to pins are why upper layers are expensive for short nets.

## The minimum for a fixed 2D route

For a two-pin route whose directions alternate, putting every horizontal segment on the lowest horizontal layer and every vertical segment on the lowest vertical layer minimises vias. Every bend needs at least one via (adjacent segments are on different layers), and each pin needs at least the stack to the lowest layer of its segment's direction, and the lowest assignment achieves exactly those bounds.

Example: an assignment using 12 vias for a five-segment route drops to $1 + 4 + 1 = 6$ on M2 and M3 (one via up from the first pin to M2, four bends, one down from M2 at the end), saving 6. Real assigners do not do this, because the lowest layers would overflow. The minimum is the baseline, and the question becomes how many vias to spend to fit capacity and meet timing.

## Vias at junctions

In a multipin tree, several segments meet at a Steiner point or a pin. The via stack at such a point must reach every layer used there, so it spans from the lowest to the highest: $\max(S) - \min(S)$ vias for the set $S$ of layers meeting at the point (including metal 1 if a pin is there). Two collinear segments on the same layer meeting at a junction need no via between them.

Example: a junction where horizontal segments on M6 and M6 meet a vertical branch on M5 needs one via. A tree's total via count sums these stacks over all junctions, bends and pins. This is why layer assignment works on whole trees, usually by dynamic programming from the leaves up: a choice on one branch changes the stack at the junction it shares with others.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); Y. Xu, Y. Zhang and C. Chu, "FastRoute 4.0" (2009) on via-aware layer assignment.
