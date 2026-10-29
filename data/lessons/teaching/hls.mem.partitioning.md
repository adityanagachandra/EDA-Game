# Array partitioning: block, cyclic and complete

An array mapped to one RAM serves at most one or two accesses per cycle, which caps the II and the useful unroll factor of any loop that reads it. Array partitioning splits the array across several banks, each with its own ports. Choosing block, cyclic or complete partitioning, and the factor, is how an HLS designer matches memory bandwidth to the parallelism the schedule wants.

## Cyclic partitioning

Cyclic partitioning by a factor $B$ deals elements round-robin: element $i$ goes to bank $i \bmod B$ at offset $\lfloor i / B \rfloor$ inside it. With $B = 6$, element 50 lands in bank 2 at offset 8. Consecutive elements always fall in different banks, so a loop that reads $a[i], a[i+1], \dots, a[i+B-1]$ in one cycle touches every bank once. That makes cyclic partitioning the natural match for loops unrolled by a factor that divides or equals $B$ and that read neighbouring elements. Each bank is a smaller RAM, or a set of registers if it is small enough, and the address logic is cheap when $B$ is a power of two, because the modulo and division become bit selections. A non-power-of-two factor needs real modulo and division hardware on every access.

## Block partitioning

Block partitioning cuts the array into $B$ contiguous runs of $\lceil N / B \rceil$ elements: element $i$ goes to bank $\lfloor i / \lceil N / B \rceil \rfloor$ at offset $i \bmod \lceil N / B \rceil$. A 144-element array in 6 blocks has 24 elements per bank, so element 13 is in bank 0 and element 50 in bank 2. Neighbouring elements share a bank, so block partitioning does nothing for a loop that reads $a[i]$ and $a[i+1]$ together. It helps when parallel accesses are far apart: a loop that processes the two halves of an array simultaneously, or an outer loop unrolled so that each copy walks its own block. Complete partitioning, $B = N$, turns every element into its own register, which gives unlimited ports at the cost of $N$ registers and wide multiplexers on every access with a variable index.

## Choosing the factor

The factor should give exactly the bandwidth the schedule needs. If a loop is unrolled by U so that each cycle reads the U consecutive elements $a[Uk]$ through $a[Uk + U - 1]$, and each bank has $p$ ports, cyclic partitioning by $B$ spreads the U reads over $B$ banks, at most $\lceil U / B \rceil$ per bank. All reads fit in one cycle when $\lceil U / B \rceil \le p$, so the smallest factor is $B = \lceil U / p \rceil$: for U = 6 with dual-port banks, $B = 3$. A larger factor wastes RAM blocks, since each bank needs at least one physical block RAM however few elements it holds, while a smaller one raises the II. Multi-dimensional arrays can be partitioned along any dimension, and the right dimension is the one the unrolled loop index walks.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on arrays; Cong, Jiang, Liu and Zou, "Automatic memory partitioning and scheduling for throughput and power optimization" (ICCAD, 2009).
