# Range limiting

Range limiting restricts proposed moves to a window around the moving block that shrinks as the anneal cools. Early on, moves span the chip; late, they are local, where small adjustments are likely to be accepted. VPR updates the window so the acceptance rate stays near 0.44, a value Lam and Delosme associated with efficient annealing.

## The range-limit update

VPR keeps a range limit $R$ and updates it after each temperature as $R_{new} = R_{old}(1 - 0.44 + R_{accept})$, clamped to $[1, \max(n_x, n_y)]$. If more than 44% of moves were accepted, the window grows; if fewer, it shrinks. From $R = 23.5$ on a $30 \times 30$ grid with acceptance 0.97, the raw update is $23.5 \times 1.53 = 35.96$, clamped to 30. From $R = 9.25$ with acceptance 0.35, the new limit is $9.25 \times 0.91 = 8.42$. At high temperature almost every move is accepted, so the window stays at the full chip; as the anneal cools, the window contracts toward single-site moves.

## Successive updates

Because the factor depends on each temperature's own acceptance rate, the range limit follows the anneal's state, and two updates compound. From $R = 9.25$, acceptance 0.35 gives 8.4175, and then acceptance 0.82 gives $8.4175 \times 1.38 = 11.62$. The window and the temperature interact: a smaller window proposes nearby moves with small $\Delta C$, which raises the acceptance rate at a given temperature. Range limiting is therefore what keeps the acceptance rate near its target late in the anneal, when global moves would almost always be rejected. Clamping matters at both ends: the window never exceeds the chip, and it never shrinks below a single neighbouring site, so the anneal can always still move.

## Counting candidate targets

A move picks a target within $\lfloor R \rfloor$ sites of the block in both $x$ and $y$, inside the grid and not the block's own site. For a block at $(x, y)$ on an $n_x \times n_y$ grid, the window covers $[\max(1, x - r), \min(n_x, x + r)]$ by the same in $y$, with $r = \lfloor R \rfloor$. For a block at $(7, 14)$ on an $11 \times 18$ grid with $R = 9$, $x$ runs from 1 to 11 and $y$ from 5 to 18, giving $11 \times 14 - 1 = 153$ targets. Blocks near the edge get smaller windows. In heterogeneous FPGAs the target must also be a site of the same block type, so a DSP block's window contains only DSP columns.

## Going deeper

Betz and Rose, FPL (1997); Lam and Delosme, DAC (1988); the VPR placer documentation.
