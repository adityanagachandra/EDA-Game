# Sparse matrix storage

A sparse matrix is stored as lists of its nonzeros, never as a dense array. The format decides what is cheap: appending entries while stamping, walking a row for a matrix-vector product, or walking a column during factorization. Every solver in an EDA flow, from KLU inside ngspice to the CG solver in an IR-drop tool, is built around one of these layouts.

## Coordinate and compressed sparse row formats

Coordinate (COO) format stores one triple (row, column, value) per nonzero: $3\,\text{nnz}$ numbers. It is the natural target while stamping a netlist, because each element simply appends triples, and duplicates are summed later.

Compressed sparse row (CSR) stores three arrays: val (the nonzero values row by row), col_idx (the column of each value) and row_ptr, of length $n + 1$, where row $i$ occupies positions row_ptr[$i$] to row_ptr[$i+1$] $- 1$. Total array length is $2\,\text{nnz} + n + 1$. For a 4-by-4 matrix with 9 nonzeros: COO stores 27 numbers, CSR stores $9 + 9 + 5 = 23$.

Compressed sparse column (CSC) is the same idea by columns. Left-looking LU codes such as KLU and the SuperLU family work on columns, so they take CSC; SpMV kernels prefer CSR.

## Reading the row pointer

With 0-based indexing, row_ptr[0] = 0 and row_ptr[$i$] is the number of nonzeros in rows 0 to $i - 1$. The difference row_ptr[$i+1$] $-$ row_ptr[$i$] is the length of row $i$, and row_ptr[$n$] = nnz.

Example: a matrix whose rows hold 3, 2, 4 and 1 nonzeros has row_ptr = [0, 3, 5, 9, 10]. Row 2 occupies val positions 5 to 8. To build row_ptr from COO, count entries per row, then take a running (prefix) sum, an $O(n + \text{nnz})$ pass. The same prefix-sum trick builds the transpose, which is how a tool gets CSC from CSR.

Off-by-one errors here are the classic sparse-code bug: the half-open range [row_ptr[$i$], row_ptr[$i+1$]) is the convention, and the extra final entry exists so that the last row has an end.

## Memory cost against dense storage

With 8-byte double values and 4-byte integer indices, CSR costs $12\,\text{nnz} + 4(n + 1)$ bytes, dense storage costs $8n^2$. A power grid with $n = 10^6$ nodes and 5 nonzeros per row needs $12 \cdot 5 \times 10^6 + 4 \times 10^6 \approx 64$ MB in CSR and $8 \times 10^{12}$ bytes (8 TB) dense: a ratio of about 125,000.

Index width matters at scale. Once nnz exceeds about $2^{31}$, 32-bit indices overflow and the code must switch to 64-bit, which raises the per-nonzero cost from 12 to 16 bytes. Blocked formats (BSR) store small dense blocks with one index each, which saves index memory when unknowns come in groups, for example the several variables of a device or a multi-layer grid.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 3; Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapter 2.
