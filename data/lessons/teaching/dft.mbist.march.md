# March tests for memories

A March test applies a fixed sequence of reads and writes to every address in a fixed order, then the same for the next element. Its length is linear in memory size, so it scales to millions of words, and each algorithm is designed to detect a known set of functional faults. Memory BIST controllers are hardware implementations of March algorithms, so counting operations is counting test cycles.

## March notation

A March test is a sequence of **March elements**. Each element has an address order, $\Uparrow$ for ascending, $\Downarrow$ for descending and $\Updownarrow$ for either, and a list of operations applied in turn to each address before moving to the next: w0 and w1 write, r0 and r1 read and expect the stated value. **March C-** is
$$\{\Updownarrow(w0); \Uparrow(r0,w1); \Uparrow(r1,w0); \Downarrow(r0,w1); \Downarrow(r1,w0); \Updownarrow(r0)\}.$$
Its complexity is the total number of operations per cell times $n$: $1 + 2 + 2 + 2 + 2 + 1 = 10$, so **10n**. MATS+ is 5n, MATS++ and March X are 6n, March Y 8n, March A 15n and March B 17n.

## Test time

Total operations are the complexity times the number of words, and with one operation per cycle the test time follows directly. March B (17n) on a 16K-word memory with a 5 ns cycle needs $17 \times 16{,}384 \times 5$ ns $= 1392.64$ microseconds. A 10n test on a 1M-word memory at 1 ns takes about 10.5 ms, which is why March algorithms, not quadratic tests such as GALPAT, are used for large arrays. Many memories on a chip are tested in parallel under power limits, and the scheduler balances test time against peak current.

## What each element buys

Each element is there for a reason. The initial $\Updownarrow(w0)$ sets a known background. Reading the expected value then writing its complement detects stuck-at and transition faults. Doing that in both address orders catches coupling faults whose victim is above or below the aggressor, and address decoder faults. March C- detects stuck-at, transition, address decoder and unlinked inversion, idempotent and state coupling faults. For **word-oriented** memories, intra-word coupling needs several **data backgrounds**, typically $\log_2 B + 1$ for $B$-bit words, so the test repeats: March X on 2048 words of 4 bits costs $6 \times 2048 \times 3 = 36{,}864$ operations.

## Going deeper

van de Goor, "Testing Semiconductor Memories" (1991); Bushnell and Agrawal (2000), chapter 9.
