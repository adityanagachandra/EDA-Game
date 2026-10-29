# Sparse matrix-vector product in CSR

The sparse matrix-vector product $y = Ax$ is the inner kernel of every Krylov solver, every power-grid iteration and every explicit time step. In CSR it is a short double loop, but its cost is set by memory traffic, not arithmetic. Knowing exactly what it touches tells a tool designer how fast an IR-drop or thermal solve can possibly run.

## The row loop over the row pointer

In CSR the product is
$$y_i = \sum_{k = \text{row\_ptr}[i]}^{\text{row\_ptr}[i+1] - 1} \text{val}[k]\; x[\text{col\_idx}[k]].$$
The outer loop runs over rows, the inner loop over the stored entries of a row. Rows are independent, so they parallelize directly.

Locating a single entry $(i, j)$ means scanning row $i$: its values start at position row_ptr[$i$], and when columns are kept sorted, the entry sits after the nonzeros of row $i$ with smaller column index. For row 1 of a matrix with row_ptr = [0, 4, 6, ...] and nonzeros of row 1 at columns 1 and 3, the entry in column 3 is at position $4 + 1 = 5$. Sorted columns also allow binary search and make the access to $x$ more regular.

## Computing one output entry

Take val = [4, -1, -1, 3, 2, 5], col_idx = [0, 2, 0, 1, 1, 2], row_ptr = [0, 2, 4, 6] and $x = (1, 2, -1)$.

Row 0 uses positions 0 and 1: $y_0 = 4 \cdot x_0 + (-1) \cdot x_2 = 4 + 1 = 5$.
Row 1 uses positions 2 and 3: $y_1 = -1 \cdot x_0 + 3 \cdot x_1 = -1 + 6 = 5$.
Row 2 uses positions 4 and 5: $y_2 = 2 \cdot x_1 + 5 \cdot x_2 = 4 - 5 = -1$.

The access pattern is the point: val and col_idx stream through memory in order, while $x$ is read at scattered positions col_idx[$k$]. Orderings that cluster nonzeros near the diagonal (reverse Cuthill-McKee, for instance) make those reads cache-friendly.

## Cost and memory traffic of SpMV

Each stored entry costs one multiply and one add: $2\,\text{nnz}$ flops. Memory traffic, under the usual model, is 12 bytes per nonzero (value and index), $4(n + 1)$ bytes of row pointers, $8n$ bytes to read $x$ once and $8n$ to write $y$:
$$\text{bytes} = 12\,\text{nnz} + 4(n + 1) + 16n.$$
For $n = 10^6$ and 5 nonzeros per row that is about 80 MB moved for $10^7$ flops, an intensity near 0.13 flop per byte.

Every modern processor can do several flops per byte of memory bandwidth, so SpMV runs at the memory system's speed. Optimizations therefore cut bytes: smaller indices, blocked formats, single-precision values in a mixed-precision solver, or reusing $x$ from cache. Floating-point tricks do not help a memory-bound kernel.

## Going deeper

Williams et al., "Optimization of sparse matrix-vector multiplication on emerging multicore platforms" (Supercomputing 2007); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 3.
