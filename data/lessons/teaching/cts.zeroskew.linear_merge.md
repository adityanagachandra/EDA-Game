# Zero-skew merging under linear delay

Zero-skew clock routing builds a tree bottom-up by merging pairs of subtrees, each already balanced, so that the merged tree is balanced too. The basic operation is choosing the tapping point on the wire between two subtree roots. Under the linear (path-length) delay model, where delay equals wire length, the choice is a one-line formula. It is the model behind DME's geometry and the right place to learn the merge before adding Elmore loading.

## Balancing the two sides

Let subtree 1 have delay $t_1$ from its root to every one of its sinks and subtree 2 have $t_2$, with the roots a wire of length $L$ apart. A tap at distance $x$ from root 1 sees delay $x + t_1$ to subtree 1's sinks and $(L - x) + t_2$ to subtree 2's. Zero skew requires these to be equal:
$$x = \frac{L + t_2 - t_1}{2}.$$
With $t_1 = 150$, $t_2 = 0$ and $L = 480$ µm, the tap is $(480 - 150)/2 = 165$ µm from subtree 1. The tap moves toward the slower subtree, so less new wire is added to the path that is already long.

## The merged delay

After the merge, the new subtree is again balanced, with delay $x + t_1 = (L + t_1 + t_2)/2$ from the tap to every sink. For $t_1 = 60$, $t_2 = 100$ and $L = 120$ µm, $x = 80$ µm and the merged delay is 140 µm. This value becomes the $t$ of the merged node when it is merged again higher up, so a whole tree is built by repeating the operation along a topology. The total wire added by a merge is exactly $L$ when the tap falls on the segment, which is why linear-delay zero-skew trees need no wire beyond the merge distances.

## When the tap leaves the segment

If $|t_1 - t_2| > L$, the formula gives $x < 0$ or $x > L$: even a tap at the faster subtree's root cannot equalise. The tap is then placed at the root of the slower subtree, and the wire to the faster subtree is elongated to length $|t_1 - t_2|$, by snaking. With $t_1 = 330$, $t_2 = 80$ and $L = 160$ µm, the wire must be 250 µm, so 90 µm of detour is added. Detour wire is pure waste in wirelength and power, so good topologies pair subtrees of similar delay, and bounded-skew methods relax the balance to avoid snaking.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter; Chao, Hsu, Ho, Boese and Kahng, "Zero skew clock routing with minimum wirelength", IEEE Transactions on Circuits and Systems II (1992).
