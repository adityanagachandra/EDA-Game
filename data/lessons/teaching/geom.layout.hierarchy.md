# Hierarchical layout and flat shape counts

A chip layout is a hierarchy: standard cells placed in blocks, blocks in the top level, memories as arrays of bit cells. Stored hierarchically it is compact; flattened, a modern design has billions of shapes. Hierarchical processing exploits repetition by doing work once per cell, and the question that decides whether it pays is which instances interact with their surroundings. Counting shapes and interactions makes the trade-off concrete.

## Flattening by multiplication

The hierarchy is a directed acyclic graph: each cell holds its own shapes and instances of other cells. The flat shape count of a cell is its own shapes plus, for each child placement, the child's flat count. Computing bottom-up, each cell once, gives the total by the product rule. Example: cell C has 40 shapes; B has 10 shapes and 6 instances of C, so $10 + 240 = 250$; A has 20 shapes, 4 instances of B and 2 of C, so $20 + 1000 + 80 = 1100$; the top has 5 shapes and 100 instances of A, so $110{,}005$ flat shapes. Counts multiply down the hierarchy, which is how a few thousand distinct cells become billions of shapes, and why flattening is the expensive step to avoid.

## What hierarchy saves

A hierarchical engine stores and checks each cell's own shapes once, not once per placement. The saving is the ratio of flat shapes to shapes stored once per cell. Example: in the hierarchy above, the cells store $40 + 10 + 20 + 5 = 75$ shapes, against 110,005 flat, a ratio of about 1466.7. The saving is real only for work that does not depend on context. A cell's internal width and spacing results are the same everywhere, but anything near its boundary depends on what is placed next to it. Hierarchical DRC engines therefore check each cell in isolation, then check interaction regions per context, and promote or flatten geometry where contexts differ too much for reuse to pay.

## Interaction between instances

Two instances interact when geometry from one comes within the rule's interaction distance of the other; then their shared border must be checked together. For a row of $k$ instances of width $w$ at pitch $p$, instances $i$ and $j$ have bounding boxes $(j - i)p - w$ apart, and they must be checked together when that gap is below the interaction distance $s$. Example: $k = 8$, $w = 100$, $p = 120$, $s = 50$: adjacent instances are 20 apart and interact, instances two apart are 140 apart and do not, so 7 pairs. Abutted standard cells interact with every neighbour, so the interaction regions are a thin band along each row. That band, not the cell interiors, is where hierarchical checking spends its effort.

## Going deeper

The KLayout documentation on deep (hierarchical) mode; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification.
