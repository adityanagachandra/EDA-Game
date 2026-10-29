# Cache organization and address fields

A cache is an SRAM data array plus a tag array that remembers which memory blocks it holds. Its organization (capacity, block size, associativity) splits every address into tag, index and offset fields and fixes the size of both arrays. These numbers decide the memory macros a physical designer must place, the comparators on the critical path, and the hit rate the architecture sees.

## Tag, index and offset

A cache of capacity $C$ bytes with block (line) size $B$ bytes and associativity $a$ has $C / B$ lines in $S = C/(B a)$ sets. A byte address of $A$ bits splits as

- offset: $\log_2 B$ bits, the byte within the block;
- index: $\log_2 S$ bits, the set;
- tag: $A - \log_2 S - \log_2 B$ bits, stored to identify the block.

A 32 KB, 4-way cache with 64-byte blocks and 32-bit addresses has $32768 / (64 \cdot 4) = 128$ sets: 6 offset bits, 7 index bits and $32 - 7 - 6 = 19$ tag bits.

## Associativity changes the split

For fixed capacity and block size, doubling the associativity halves the number of sets, removes one index bit and adds one tag bit. The same 32 KB cache with 64-byte blocks has 512 sets when direct-mapped (9 index bits, 17 tag bits), 128 sets when 4-way, and 1 set when fully associative (0 index bits, 26 tag bits).

Higher associativity reduces conflict misses, but every way needs a tag comparator active on each access and a wider way-select mux, so hit time and energy rise. A fully associative cache is a CAM. This is why first-level caches stay at modest associativity (often 4 to 8 ways) and large last-level caches use more.

## Total storage

Each line stores its data ($8B$ bits), its tag, a valid bit and, for a write-back cache, a dirty bit. With $L = C/B$ lines,

$$\text{bits} = L \cdot (8B + \text{tag} + 1 + 1).$$

For the 32 KB, 4-way, 64-byte write-back cache: 512 lines of $512 + 19 + 2 = 533$ bits, 272896 bits in total, about 4.1% more than the data alone. Replacement state (LRU bits) adds a little more per set. In the physical design the data and tag arrays are separate SRAM macros, often with the tag array accessed first or in parallel, and their shapes come straight from this arithmetic.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach" (6th ed., 2017), appendix B; Smith, "Cache Memories" (ACM Computing Surveys, 1982).
