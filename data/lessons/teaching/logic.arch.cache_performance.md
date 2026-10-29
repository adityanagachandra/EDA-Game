# Cache hits, misses and AMAT

A cache's value is measured by how often it hits and how much a miss costs. Average memory access time combines the two, and tracing an address stream through a small cache shows where misses come from: cold starts, limited capacity, and conflicts between addresses that share a set. These traces are what cache simulators compute, and the misses they count become CPI.

## Average memory access time

With hit time $t_h$, miss rate $m$ and miss penalty $t_p$,

$$AMAT = t_h + m \cdot t_p.$$

A 1-cycle L1 with a 5% miss rate and a 100-cycle penalty gives $AMAT = 6$ cycles. With an L2 the penalty itself becomes an access time: $AMAT = t_{h1} + m_1 (t_{h2} + m_2 t_{mem})$. With $t_{h2} = 10$, local L2 miss rate 20% and $t_{mem} = 100$: $1 + 0.05 (10 + 20) = 2.5$ cycles. Small changes in L1 miss rate matter a great deal, because they multiply the whole lower hierarchy.

## Tracing a direct-mapped cache

In a direct-mapped cache with $S$ sets of $B$-byte blocks, byte address $x$ belongs to block $\lfloor x / B \rfloor$, which maps to set $\lfloor x / B \rfloor \bmod S$. Each set holds one block; an access hits if that set already holds its block, otherwise it misses and replaces the set's block.

With 4 sets of 16 bytes, trace addresses 0, 4, 16, 64, 0, 20, 68: blocks 0, 0, 1, 4, 0, 1, 4 in sets 0, 0, 1, 0, 0, 1, 0. Results: miss (cold), hit, miss (cold), miss (block 4 evicts block 0 from set 0), miss (block 0 evicts block 4), hit, miss. Two hits in seven. Blocks 0 and 4 conflict in set 0 even though the cache has empty sets: conflict misses.

## Associativity and LRU

A 2-way set-associative cache of the same capacity has 2 sets of 2 blocks each, with least-recently-used replacement within a set. The same trace maps blocks 0, 0, 1, 4, 0, 1, 4 to sets 0, 0, 1, 0, 0, 1, 0. Set 0 can now hold blocks 0 and 4 together: miss, hit, miss, miss, hit, hit, hit. Four hits in seven.

Associativity removed the conflict misses, but it cannot remove cold misses (first touches) or capacity misses (a working set larger than the cache). Cache simulators classify misses into these three categories, and architects tune size, block size and associativity against the hit-time and energy costs from the organization topic. The same trace-driven method, at scale, decides the cache configuration a memory compiler is asked to build.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach" (6th ed., 2017), chapter 2 and appendix B; Hill and Smith, "Evaluating Associativity in CPU Caches" (IEEE Transactions on Computers, 1989).
