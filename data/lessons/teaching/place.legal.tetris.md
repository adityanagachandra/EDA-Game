# Tetris legalization

Tetris legalization is the simplest practical legalizer: sort cells by x and drop each into the row where it moves least, at the leftmost free position. It runs in $O(n \log n)$ for sorting plus a cheap scan per cell, and it is the baseline every better legalizer is measured against. Its greedy order explains both its speed and its typical failure.

## Sort, then place greedily

Sort the cells by their global x coordinate. Keep, for each row, the right edge of the cells already placed in it, initially the row's left boundary. For each cell in sorted order, consider every row (or the rows near the cell): the cell would go at $x = \max(x_{desired}, \text{right edge of that row})$ and its displacement would be $|x - x_{desired}| + |y_{row} - y_{desired}|$. Put it in the row with the smallest displacement, ties broken by a fixed rule, and advance that row's right edge by the cell's width. In a single row with cells wanting x = 2, 6, 10 and 19 and widths 6, 4, 2 and 5, the cells land at 2, 8, 12 and 19.

## Choosing the row

The row choice trades vertical against horizontal movement. A cell that wants (1, 2) with rows at y = 0 and y = 2 may prefer the upper row even though another row is emptier, because vertical moves cost displacement too. With two rows and many cells, a cell late in the order often finds its preferred row filled up to the right of its desired x, so it either shifts right in that row or moves to a neighbouring row. Restricting the search to a few rows above and below keeps each step cheap; scanning every row is $O(n \cdot \text{rows})$. Implementations track free segments per row so that macros and blockages, which split rows, are handled by the same scan.

## Strengths and failure modes

Tetris is fast, simple and always legal if there is room. Its weakness is the greedy order: early cells never move to make room for later ones, so cells at the right end of crowded rows can be pushed far right or into distant rows, giving a long tail of displacement. It also ignores the original relative order across rows and minimises only each cell's own displacement. Total displacement in a two-row example might be small while one late cell moves several units. Improvements include processing from both ends, limiting per-row fill, and replacing the greedy choice with Abacus-style clustering, which lets earlier cells shift left to share the displacement. Tetris remains useful as a quick legalizer inside fast loops.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.4; Spindler, Schlichtmann and Johannes, "Abacus" (ISPD 2008), which compares against Tetris.
