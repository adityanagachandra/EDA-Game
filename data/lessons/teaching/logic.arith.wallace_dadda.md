# Wallace and Dadda trees

A Wallace or Dadda tree reduces a matrix of partial-product bits to two rows using full adders (3:2 counters) and half adders, with no carry propagation, so the reduction takes logarithmically many levels. Wallace reduces as aggressively as possible at every level; Dadda reduces as little as possible while keeping to the same number of levels, which saves adders. Counting levels and adders is how multiplier architectures are compared.

## Levels of 3:2 reduction

A full adder takes three bits of one column and emits a sum bit in that column and a carry bit in the next. Applied across the matrix, each level turns every group of three rows into two, so a height of $r$ becomes $2\lfloor r/3 \rfloor + (r \bmod 3)$, which equals $\lceil 2r/3 \rceil$.

Dadda states the same schedule from the bottom: the maximum heights that can be reduced to 2 in $j$ levels are $d_1 = 2$, $d_{j+1} = \lfloor 1.5\, d_j \rfloor$, giving $2, 3, 4, 6, 9, 13, 19, 28, 42, 63, \dots$. The number of levels for height $r$ is the number of terms of this sequence below $r$. An $8 \times 8$ multiplier has maximum column height 8, and $2, 3, 4, 6 < 8$, so 4 levels; a $32 \times 32$ multiplier needs 8 levels.

## Counting Dadda's adders

Every full adder takes in three bits and emits two, so it removes exactly one bit from the matrix; a half adder takes two and emits two, removing none (it only moves a bit to the next column). An $n \times n$ multiplier starts with $n^2$ bits. Dadda's minimal-reduction rule yields closed forms (Townsend, Swartzlander and Abraham): an $n \times n$ Dadda tree uses

$$n^2 - 4n + 3 \text{ full adders and } n - 1 \text{ half adders},$$

followed by a $(2n - 2)$-bit carry-propagate adder. For $n = 8$: 35 full adders and 7 half adders. A Wallace tree for the same multiplier has the same number of levels but uses more half adders, because it reduces every column eagerly.

## Booth recoding meets the tree

Radix-4 Booth recoding halves the number of rows, so a $32 \times 32$ multiplier presents about 16 rows (plus the bits that complete the negations) instead of 32. The Dadda count drops from 8 levels (for height 32) to 6 (for height 16: the terms 2, 3, 4, 6, 9, 13 lie below 16). Each level costs one full-adder delay, about two gate delays, so Booth saves about four gate delays here, against the cost of the Booth encoders and partial-product muxes.

4:2 compressors, built from two full adders but with a carry path designed to stay inside the cell, reduce four rows to two per level, so they halve the height each level and give a more regular tree. Whether Booth plus a Dadda tree, or a plain array, or 4:2 compressors win depends on the target's timing, area and wiring, a classic space for automated design exploration.

## Going deeper

Wallace, "A Suggestion for a Fast Multiplier" (IEEE Transactions on Electronic Computers, 1964); Dadda, "Some Schemes for Parallel Multipliers" (Alta Frequenza, 1965); Townsend, Swartzlander and Abraham, "A Comparison of Dadda and Wallace Multiplier Delays" (SPIE, 2003).
