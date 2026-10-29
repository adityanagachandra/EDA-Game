# GPU computing for EDA kernels

GPUs offer far more arithmetic and memory bandwidth than CPUs, but EDA kernels benefit only when they fit the hardware. The roofline model, built on arithmetic intensity, predicts which kernels speed up and by how much. It tells a tool designer that sparse solvers are bandwidth-bound and where GPU effort pays.

## Arithmetic intensity of sparse kernels

Arithmetic intensity is flops performed per byte moved from memory. For CSR SpMV, flops $= 2\,\text{nnz}$ and, with 8-byte values, 4-byte indices, and $x$ read and $y$ written once,
$$\text{bytes} = 12\,\text{nnz} + 4(n + 1) + 16n.$$
For $n = 10^6$ and 7 nonzeros per row: $14 \times 10^6$ flops over $84 + 4 + 16 = 104$ MB, an intensity of 0.135 flop/byte.

By contrast, dense matrix multiply on $b \times b$ blocks performs $2b^3$ flops on $3 \cdot 8b^2$ bytes, an intensity of $b/12$ that grows with block size. Iterative solvers, explicit stencils and most sparse kernels sit below 0.25 flop/byte; dense factorization and GEMM sit far above.

## The roofline bound

A processor with peak $P$ flop/s and memory bandwidth $B$ byte/s can sustain at most
$$\text{attainable} = \min(P,\ I \cdot B)$$
for a kernel of intensity $I$. The crossover $I^* = P/B$ is the machine balance. For a GPU with $P = 20$ TFLOP/s and $B = 2$ TB/s, $I^* = 10$ flop/byte. SpMV at $I = 0.135$ attains at most $0.135 \cdot 2000 = 270$ GFLOP/s, about 1.4% of peak: it is memory-bound, and its speed is set by bandwidth alone.

So the speedup of a GPU over a CPU on SpMV-dominated solvers is roughly the bandwidth ratio, perhaps 5 to 10 times, not the flop ratio. Tuning efforts reduce bytes (narrower indices, FP32 values where accuracy allows, blocked formats) rather than flops.

## Mapping solver kernels to GPU threads

Simple CSR kernels assign a thread per row; rows of very different lengths then leave threads idle, so better kernels assign a warp per row or use formats such as ELL or SELL-C-sigma that pad rows into regular slices. Power grids and thermal meshes, with nearly uniform rows, map well; netlists with dense supply rows need special handling.

Krylov solvers map well overall: SpMV, dot products (reductions) and vector updates, all bandwidth-bound and data-parallel. Sparse direct factorization of circuit matrices does not: small, irregular supernodes and sequential dependencies limit parallelism, though batched factorization of many small matrices of the same pattern works well. Libraries such as cuSPARSE provide the sparse kernels. The time of one SpMV is about bytes divided by bandwidth: 104 MB at 2 TB/s is 52 µs, so a 500-iteration CG solve takes tens of milliseconds.

## Going deeper

Williams, Waterman and Patterson, "Roofline: an insightful visual performance model for multicore architectures", Communications of the ACM (2009); Bell and Garland, "Implementing sparse matrix-vector multiplication on throughput-oriented processors", Supercomputing (2009).
