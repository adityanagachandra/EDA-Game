# Amortized analysis

Amortized analysis bounds the total cost of a sequence of operations, even when single operations are occasionally expensive. Boolean engines are full of such structures: unique tables and computed tables in BDD packages that double when full, clause databases that grow and are periodically reduced, the trail of a SAT solver, incremental timing graphs. A designer who knows the amortized cost can promise constant time per operation where a worst-case bound would say linear.

## The binary counter

Incrementing a binary counter flips all trailing 1s to 0 and one 0 to 1. A single increment can flip many bits (0111 to 1000 flips four), but bit $j$ flips only once every $2^j$ increments. Over $N$ increments from 0 the total is $\sum_j \lfloor N / 2^j \rfloor$, which equals $2N - s(N)$, where $s(N)$ is the number of 1s in the binary form of $N$.

For $N = 8$ the flips per increment are 1, 2, 1, 3, 1, 2, 1, 4, totaling 15 $= 16 - 1$. So the amortized cost is under 2 flips per increment, although the worst single increment costs $\log_2 N + 1$. The same argument bounds Gray-code enumeration and the cost of incrementally stepping through cofactor assignments.

## Doubling hash tables

A hash table that doubles its capacity when full copies every stored entry into the new array. Starting from capacity 1, the copies happen at sizes 1, 2, 4, ..., up to the largest power of two below $N$, so the total number of copies after $N$ inserts is less than $2N$. With 9 inserts the table grows at sizes 1, 2, 4 and 8, copying $1 + 2 + 4 + 8 = 15$ entries.

BDD unique tables, AIG structural-hashing tables and SAT watch lists all grow this way. Growing by a constant factor is essential: growing by a fixed amount (say 1,000 slots each time) makes the total copying quadratic in $N$.

## Cost per operation

Add the direct cost and the copying cost, then divide by the number of operations. For the 9 inserts above: 9 units of insertion plus 15 units of copying is 24, or about 2.67 per insert. The potential method gives a clean bound: with $\Phi = 2 \cdot \text{size} - \text{capacity}$, each insert has amortized cost at most 3, because the potential saved by cheap inserts pays for the next doubling.

Amortized bounds are about totals, not latency. A single resize of a unique table holding $10^8$ nodes still stalls the tool, which is why production BDD packages resize incrementally or during garbage collection, and why an interactive or agent-driven loop that needs predictable response times may prefer worst-case structures.

## Going deeper

Tarjan, "Amortized computational complexity", SIAM Journal on Algebraic and Discrete Methods (1985). Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter 16.
