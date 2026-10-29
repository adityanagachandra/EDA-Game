# Spreading with anchors and pseudo-nets

A quadratic placement solution is wirelength-optimal and hopelessly overlapped. Spreading turns it into a distributed placement. Modern quadratic placers such as SimPL compute a spread-out target for every cell and pull cells toward it with anchor pseudo-nets, solving the linear system again, so that wirelength and density are traded gradually over many iterations.

## Anchors as extra fixed pins

An anchor is a pseudo-net, a two-pin connection of weight $\alpha$ from a movable cell to a fixed point at its target position $t$. To the linear system it is just another fixed connection: it adds $\alpha$ to the cell's diagonal and $\alpha t$ to the right-hand side. For a single cell whose real connections give $\sum w_j x_j = S$ and $\sum w_j = W$, the new optimum is
$$x = \frac{S + \alpha t}{W + \alpha}.$$
With pins at 79 (weight 1) and 29 (weight 2), the plain solution is 45.67; an anchor of weight 4 at 96 moves the cell to $(137 + 384)/7 = 74.43$. The anchor does not force the cell to its target; it balances the target against the cell's connections.

## Where the targets come from

SimPL alternates a lower-bound step and an upper-bound step. The lower bound is the quadratic solve, wirelength-optimal and overlapped. The upper bound is look-ahead legalization: over-full regions are found, expanded until they have enough capacity, and their cells are spread inside them by recursive geometric partitioning that preserves the cells' relative order. The spread positions become the anchor targets for the next solve. As iterations proceed, the anchor weight grows, often linearly with the iteration count, and in SimPL it is also divided by the cell-to-target distance, linearised like a B2B edge, so that the pull behaves like a linear rather than quadratic penalty. The gap between the two bounds shrinks until they meet.

## Choosing the anchor weight

The weight decides how far the cell moves toward its target. Solving $x = p$ for the weight gives $\alpha = (S - pW)/(p - t)$. If the cell should land a fraction $f$ of the way from its plain solution to its target, this simplifies to $\alpha = Wf/(1-f)$, independent of where the pins are: halfway needs $\alpha = W$. Too large a weight early freezes cells at spread positions computed from a poor solution; too small a weight late lets cells collapse back. That schedule is a tunable placer parameter, a Tier 1 change, while the spreading algorithm that produces targets is strategy and the linear solver is core.

## Going deeper

Kim, Lee and Markov, "SimPL: an effective placement algorithm" (ICCAD 2010); Viswanathan, Pan and Chu, "FastPlace 3.0" (ASP-DAC 2007).
