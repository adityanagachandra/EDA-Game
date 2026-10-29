# Bounding-box cost with the q(n) correction

VPR's wirelength cost for placement is a bounding-box estimate of each net's routing, corrected for nets with many terminals by the factor $q(n)$ and normalized by channel capacity. It is cheap to update incrementally, which matters because an annealer evaluates millions of moves, and it correlates well with the channel tracks the router will use.

## The cost of one net

For a net with $n$ terminals, VPR uses $q(n)\left[\frac{b_x}{C_{av,x}} + \frac{b_y}{C_{av,y}}\right]$, where $b_x$ and $b_y$ are the numbers of channels the bounding box spans and $C_{av}$ the average channel capacity in each direction. With uniform channels the capacity is a constant factor, so the cost is $q(n)\left[(x_{max} - x_{min} + 1) + (y_{max} - y_{min} + 1)\right]$; the $+1$ counts channels, since a net confined to one column still uses a channel. For a 7-terminal net whose box spans 9 columns and 9 rows, the cost is $q(7)(9 + 9) = 1.2823 \times 18 = 23.08$. The capacity terms matter in architectures whose channels differ in width, steering wiring toward wider channels.

## The q(n) correction

Half-perimeter wirelength equals the rectilinear Steiner tree length for nets of up to 3 terminals, but it underestimates larger nets, whose trees must cross the box more than once. $q(n)$ is the expected ratio, computed by Cheng from random nets: $q(n) = 1$ for $n \le 3$, $q(4) = 1.0828$, $q(5) = 1.1536$, $q(6) = 1.2206$, $q(7) = 1.2823$, $q(8) = 1.3385$, $q(9) = 1.3991$, $q(10) = 1.4493$, rising slowly beyond. A placement's total cost is the sum over nets. For nets of 4, 5 and 3 terminals with spans 9, 12 and 6, the cost is $1.0828 \times 9 + 1.1536 \times 12 + 6 = 29.59$. Without $q(n)$, the placer would undervalue shrinking high-fanout nets.

## Incremental cost of a move

When a block moves, only the nets connected to it change, so the annealer recomputes just those. Moving block B from $(4, 6)$ to $(3, 3)$ on a 2-terminal net with the other terminal at $(4, 5)$ changes the spans from $1 + 2 = 3$ to $2 + 3 = 5$; on a 3-terminal net with fixed terminals at $(3, 5)$ and $(5, 4)$, the box spans $3 + 3 = 6$ before and $3 + 3 = 6$ after. The change is $+2$. VPR goes further and updates each bounding box from edge counts (how many terminals lie on each edge) in constant time when possible, rescanning only when a moving block was the sole terminal on an edge. This is what makes millions of evaluated moves affordable.

## Going deeper

Betz and Rose, "VPR: A New Packing, Placement and Routing Tool for FPGA Research", FPL (1997); Cheng, "RISA: Accurate and Efficient Placement Routability Modeling", ICCAD (1994); Betz, Rose and Marquardt (1999), chapter 5.
