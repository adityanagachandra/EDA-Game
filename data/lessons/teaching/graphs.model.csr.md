# Compressed incidence storage for netlists

A million-cell netlist has several million pins, and every pass of a partitioner or placer walks them. The data structure that holds incidence decides cache behaviour, memory use and the cost of each traversal. Engines almost universally store a netlist in compressed sparse form, an offset array plus an index array in each direction, so a tool designer must be able to compute offsets and sizes and know why this layout wins.

## Offsets and pin arrays

Compressed storage of net-to-pin incidence uses two arrays. The pin array lists the cells of net 0, then net 1, and so on, back to back. The offset array netStart has length $m + 1$ for $m$ nets: netStart[j] is where net $j$'s pins begin, and netStart[m] equals the total pin count $P$. Net $j$'s pins are pin[netStart[j]] through pin[netStart[j+1] - 1], and its size is the difference of consecutive offsets.

With nets of sizes 2, 4, 4 and 2, the offsets are 0, 2, 6, 10, 12. Net 2 (the third) starts at index 6 and has $10 - 6 = 4$ pins. Each offset is a prefix sum of net sizes, so the array is built in one pass after counting.

## The transpose: cell-to-net incidence

Algorithms need both directions. An FM move starts from a cell and must reach its nets, then from each net reach its other cells. So engines also store cell-to-net incidence: an offset array cellStart of length $n + 1$ for $n$ cells and an index array of net ids. cellStart[c] is the sum of the degrees of the cells before $c$.

Building the transpose is a counting sort: count each cell's degree in one pass over the pins, take prefix sums to get cellStart, then make a second pass that drops each net id into its cell's slot. If degrees are 2, 2, 2, 4, 1, then cellStart is 0, 2, 4, 6, 10, 11. Both directions together are what make one FM pass linear in $P$.

## Why compressed storage wins

Storing both directions takes $(m + 1) + P + (n + 1) + P$ integers, about $2P$ for large designs. With 4-byte integers, 6 cells, 6 nets and 17 pins need $4 \times (7 + 17 + 7 + 17) = 192$ bytes. A pointer-based structure (a list object per net and per cell) costs several times more and scatters data across memory.

The contiguous layout means a traversal of a net's pins is a sequential read, which caches and prefetches well and vectorizes. The price is that the structure is static: inserting a pin means shifting arrays. Engines therefore build the compressed form once per run, or once per level in a multilevel partitioner, and handle incremental netlist edits (ECOs, buffer insertion) with a separate dynamic structure or a periodic rebuild.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003), chapter 3, for the compressed sparse row format; the KaHyPar and hMETIS source code for hypergraph storage.
