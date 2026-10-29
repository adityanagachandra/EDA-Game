# Shape functions and block orientation

Blocks often have more than one allowed shape: a hard macro may rotate, and a soft block may take any aspect ratio within limits. A shape function lists a block's achievable (width, height) pairs. Combining shape functions up a slicing tree, keeping only the non-dominated pairs, finds the best orientation of every block at once, which is Stockmeyer's classical result for slicing floorplans.

## Shapes of one block

A shape function of a block is the set of (width, height) pairs it can take. A rotatable hard block of 8 by 15 has two: (8, 15) and (15, 8). A soft block of area 120 with aspect ratio between 0.5 and 2 has a continuous curve $h = 120/w$ for $w$ between about 7.75 and 15.5. A pair is dominated if another pair is no wider and no taller, since the dominated shape can never help. The remaining non-dominated pairs form a staircase: as width grows, height falls. Listing a hard block's shapes is trivial, but this staircase view is what makes combination efficient.

## Combining two blocks under a cut

Under a vertical cut, each choice of shapes $(w_a, h_a)$ and $(w_b, h_b)$ gives a combined shape $(w_a + w_b, \max(h_a, h_b))$; under a horizontal cut, $(\max(w_a, w_b), h_a + h_b)$. With A = 8 by 15 rotatable and B = 11 by 10 fixed, a vertical cut gives (19, 15) or (26, 10), areas 285 and 260, so rotating A wins. With both blocks rotatable there are four combinations. Pruning dominated results after each combination keeps the list short: for staircases of sizes $m$ and $n$ the non-dominated combined list has at most $m + n - 1$ entries and can be built by a merge-like walk, not by trying all $mn$ pairs.

## Orientation optimisation on a slicing tree

Stockmeyer's algorithm computes the shape function of every node bottom-up, pruning dominated pairs, and picks the smallest-area shape at the root; a top-down pass then recovers each block's chosen orientation. For slicing floorplans this is polynomial, roughly linear in the number of blocks times the length of the shape lists. For general non-slicing floorplans the same orientation problem is NP-hard, which is one reason slicing representations remained popular. Fixed-outline floorplanners use the root's shape function differently: instead of minimum area they pick the shape that fits the outline, or report that none does. In an annealer, the shape computation runs inside every move evaluation, so its cost directly limits the number of moves per second.

## Going deeper

Stockmeyer, "Optimal orientations of cells in slicing floorplan designs" (Information and Control, 1983); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 3.6.
