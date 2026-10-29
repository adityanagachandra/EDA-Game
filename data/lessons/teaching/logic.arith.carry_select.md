# Carry-select and carry-skip adders

Two classic tricks speed up an adder without full lookahead. A carry-select adder computes each block's sum twice, once for each possible carry-in, and picks the right one with a mux when the real carry arrives. A carry-skip adder lets a carry jump over a block whose bits all propagate. Both trade area for delay, and both have an optimal block size that is a nice small optimization problem.

## Speculate both carries, then select

Split an $n$-bit adder into blocks of $k$ bits. The first block is an ordinary ripple adder. Every other block contains two ripple adders, one assuming carry-in 0 and one assuming carry-in 1, plus muxes. All blocks compute in parallel and finish at $k\, t_c$. The real carry then passes from block to block through one mux each:

$$t = k\, t_c + \left(\frac{n}{k} - 1\right) t_{mux}.$$

For $n = 32$, $k = 8$, $t_c = 2$ and $t_{mux} = 2$ gate delays: $16 + 3 \cdot 2 = 22$ gate delays to the carry-out, against 64 for the ripple adder. The cost is nearly double the adder area. This is Shannon expansion on the carry-in: each block computes both cofactors and the mux chooses.

## Choosing the block size

The delay $k\, t_c + (n/k - 1) t_{mux}$ has two competing terms: bigger blocks ripple longer, smaller blocks need more muxes. Treating $k$ as continuous, the minimum is at $k^* = \sqrt{n\, t_{mux} / t_c}$, giving delay $O(\sqrt{n})$.

For $n = 64$ with $t_c = t_{mux} = 2$, the divisors of 64 give: $k = 4$: $8 + 15 \cdot 2 = 38$; $k = 8$: $16 + 7 \cdot 2 = 30$; $k = 16$: $32 + 3 \cdot 2 = 38$. The best uniform block is $k = 8 = \sqrt{64}$, at 30 gate delays. Non-uniform blocks do better: since later blocks have longer to compute before their select arrives, they can be larger (block sizes growing 1, 2, 3, ...), which is the square-root carry-select adder.

## Carry-skip

A carry-skip (carry-bypass) adder adds a mux per block that forwards the block's carry-in directly to its carry-out when every bit in the block propagates ($P = p_{k-1} \cdots p_0 = 1$). In the worst case, a carry generated in the first block ripples through it, skips every middle block, and ripples through the last:

$$t \approx k\, t_c + \left(\frac{n}{k} - 2\right) t_{skip} + k\, t_c.$$

With $n = 32$, $k = 4$, $t_c = 2$, $t_{skip} = 2$: $8 + 6 \cdot 2 + 8 = 28$ gate delays. Carry-skip costs much less area than carry-select (one mux and an AND per block), but its delay is harder for a static timer to see: the long ripple path through every block exists structurally but cannot be sensitized, a false path. Without false-path handling the timer reports the pessimistic ripple delay.

## Going deeper

Koren, "Computer Arithmetic Algorithms" (2nd ed., 2002), chapter 5; Parhami, "Computer Arithmetic: Algorithms and Hardware Designs" (2nd ed., 2010), chapter 7.
