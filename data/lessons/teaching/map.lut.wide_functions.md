# LUT count and depth for wide functions

Simple counting arguments bound what any LUT mapper can achieve on wide functions. They tell a tool designer when a mapping result is already optimal, how deep a wide AND or parity must be, and how large a generic function can get. These bounds also explain why FPGA vendors add carry chains and wide multiplexers next to LUTs.

## LUT count for associative functions

A k-input LUT can combine at most k signals into one, so each LUT reduces the number of signals still to be combined by at most k - 1. Combining m inputs into one output therefore needs at least $\lceil (m - 1)/(k - 1) \rceil$ LUTs, and a tree of LUTs meets this bound for any associative operator: AND, OR, XOR. A 16-input AND in 4-LUTs needs $\lceil 15/3 \rceil = 5$ LUTs: four LUTs on the inputs and one combining their outputs. A 9-input OR in 3-LUTs needs $\lceil 8/2 \rceil = 4$.

## Minimum depth

A LUT tree of depth d can combine at most $k^d$ inputs, since each level multiplies the number of reachable inputs by at most k. The minimum depth of an m-input AND, or any function that depends on all m inputs, is therefore the smallest d with $k^d \ge m$, that is $\lceil \log_k m \rceil$. A 15-input AND in 5-LUTs needs depth 2, since 5 < 15 <= 25. Depth and count can conflict: a minimum-depth tree may need more LUTs than the count bound when m is not a convenient power of k.

## Generic functions by Shannon expansion

An arbitrary function of m > k inputs can be built by Shannon expansion, $f = x f_x + \bar{x} f_{\bar{x}}$: implement both cofactors on m - 1 inputs and select between them with a 2-to-1 multiplexer, which fits in one LUT when k >= 3. The count satisfies N(k) = 1 and N(j) = 2 N(j - 1) + 1, so $N(m) = 2^{m - k + 1} - 1$; a 7-input function in 5-LUTs uses 7 LUTs this way, at depth 3. This is an upper bound: a particular function often decomposes far better, and counting arguments show most functions of many inputs need a number of LUTs exponential in m. Vendors add dedicated multiplexers to absorb exactly these select stages.

## Going deeper

Murgai, Brayton and Sangiovanni-Vincentelli, "Logic Synthesis for Field-Programmable Gate Arrays" (1995); Shannon, "The synthesis of two-terminal switching circuits", Bell System Technical Journal, 1949, for the counting argument.
