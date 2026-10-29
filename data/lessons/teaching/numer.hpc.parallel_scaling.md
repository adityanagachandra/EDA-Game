# Parallel speedup and its limits

Parallel hardware speeds up EDA engines only as far as their serial parts allow. Amdahl's and Gustafson's laws quantify this for fixed and growing problems, and the elimination tree of a sparse factorization shows where a direct solver's parallelism actually comes from. Together they let a tool designer predict scaling before buying hardware.

## Amdahl's law for a fixed problem

If a fraction $p$ of a program's single-thread run time parallelizes perfectly over $N$ processors and the rest is serial,
$$S(N) = \frac{1}{(1 - p) + p/N}, \qquad \lim_{N \to \infty} S = \frac{1}{1 - p}.$$
A transient simulation spending 90% of its time in device evaluation and solves, and 10% in serial bookkeeping, gets $S(16) = 1/(0.1 + 0.05625) = 6.4$ and can never exceed 10.

The lesson for engine design is that the serial fraction must shrink as hardware grows: timestep control, convergence checks, output writing and any single-threaded phase become the bottleneck. Profiling the serial fraction is the first step in any parallelization effort.

## Gustafson scaling for growing problems

Often the problem grows with the machine: a larger power grid or a finer thermal mesh. If, in the parallel run on $N$ processors, a fraction $p$ of the time is parallel work, the same work would take $(1 - p) + pN$ times as long on one processor, so the scaled speedup is
$$S = (1 - p) + pN.$$
With $p = 0.95$ and $N = 64$: $S = 0.05 + 60.8 = 60.85$.

The two laws do not conflict; they answer different questions, strong scaling (same problem, faster) and weak scaling (bigger problem, same time). EDA has both needs: faster turnaround on a fixed design, and full-chip analysis that used to require partitioning. Communication cost, not just serial code, usually limits weak scaling in distributed solvers.

## Parallelism in the elimination tree

For a sparse Cholesky factor, the elimination tree has parent$(j)$ = the row index of the first off-diagonal nonzero in column $j$ of $L$. Column $j$ can be factored once all its descendants are done, and columns in disjoint subtrees are independent. With unit work per column, the tree height is the critical path, and the speedup is at most $n/\text{height}$.

Example: 8 columns with parent = 3, 3, 4, 8, 6, 7, 8 for columns 1 to 7. The longest chain is 1 or 2, 3, 4, 8: height 4 columns, while 5, 6, 7, 8 is also 4. Speedup at most $8/4 = 2$. A tridiagonal matrix in natural order gives a path, height $n$, no parallelism; nested dissection gives a short, bushy tree whose independent subtrees run concurrently, which is why parallel sparse solvers use it.

## Going deeper

Amdahl, "Validity of the single processor approach to achieving large scale computing capabilities", AFIPS (1967); Gustafson, "Reevaluating Amdahl's law", Communications of the ACM (1988); Liu, "The role of elimination trees in sparse factorization", SIAM J. Matrix Anal. Appl. (1990).
