# Deterministic parallelism

A deterministic engine gives bit-identical results for identical inputs, every run. Parallelism threatens this even without bugs: threads finish in different orders, floating-point sums depend on order, and shared random generators hand out numbers by timing. EDA teams insist on determinism because regression testing, bug reproduction, result caching and QoR comparison all assume it. A wirelength that drifts by 0.01% between identical runs hides real 0.01% regressions and makes field bugs impossible to reproduce.

## Floating-point sums depend on order

Floating-point addition rounds after every operation, so it is not associative: $(a + b) + c$ can differ from $a + (b + c)$. In decimal arithmetic rounded to 4 significant digits, $(1000 + 0.4) + 0.4$ rounds to 1000 at each step, while $1000 + (0.4 + 0.4) = 1000 + 0.8 = 1000.8$, which rounds to 1001. Double precision behaves the same way in its last bits.

A parallel reduction (the total wirelength, the density penalty, a gradient sum) combines partial sums in whatever order threads finish. Run twice, it can produce totals that differ in the last digits, with no race and no bug. In an iterative optimizer that difference does not stay small: it flips a comparison, a different move is accepted, and the final placements diverge. So bitwise reproducibility, not just "close enough," is the usual requirement.

## Fixed partitions and ordered merges

The cure is to make every rounding step and every merge independent of timing. Three techniques cover most cases. First, partition statically: split the index range into fixed chunks and compute a partial result per chunk. Second, merge in a fixed order: combine chunk results by chunk index, never as they arrive, and write per-item results into slots keyed by item index rather than appending to a shared list. Third, key randomness to work, not to threads: derive each item's random stream from the run seed and the item index, so the same item sees the same numbers whatever thread runs it.

If the chunking is chosen independently of the thread count, the result is also identical on 4 threads and on 32. A mutex is not a fix: it removes a race, but the order in which threads take the lock still varies, so an order-dependent result still varies.

## Deterministic tie-breaking and thread-count independence

Even single-threaded code hides nondeterminism. Iterating a hash map whose order depends on insertion history or a randomized hash seed, sorting by pointer address, reading uninitialized memory, or picking "the first" of equal-cost candidates from an unordered container all make results vary between runs. The rule is to break every tie with a stable key, such as an instance index or a name, and to sort by persistent ids, never by addresses.

The payoff is practical. A bug found on a 64-core server reproduces on a laptop. A regression that changes QoR on a commit touching only comments points to a determinism bug, not noise. Differential tests can demand identical outputs. Engines such as OpenROAD's tools aim to give the same result for the same inputs and seed, and treat any violation as a defect. The cost is usually a few percent of speed, from static partitioning and ordered merges.

## Going deeper

Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic" (ACM Computing Surveys, 1991); Demmel and Nguyen, "Fast Reproducible Floating-Point Summation" (ARITH 2013); Bocchino et al., "Parallel Programming Must Be Deterministic by Default" (HotPar 2009).
