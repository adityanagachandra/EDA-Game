# Carry-lookahead adders

A carry-lookahead adder computes carries directly from generate and propagate signals instead of waiting for them to ripple. Unrolling the carry recurrence gives two-level logic for each carry, and grouping bits into blocks with group generate and group propagate keeps the gates small. It is the conceptual bridge from the ripple adder to parallel-prefix adders.

## Generate and propagate per bit

Every column computes $g_i = a_i b_i$ and $p_i = a_i \oplus b_i$ in one gate delay, independent of every other column. For $a = 10110110$ and $b = 01100101$, the generate vector $a \mathbin{\&} b = 00100100$ has 2 set bits and the propagate vector $a \oplus b = 11010011$ has 5.

No column has both bits set, because $g_i = 1$ requires $a_i = b_i$ and $p_i = 1$ requires $a_i \ne b_i$. Columns with neither are kill columns ($a_i = b_i = 0$): they absorb any incoming carry. Every fast adder starts from these two vectors; only how it combines them differs.

## Unrolling the carry recurrence

Substituting $c_{i+1} = g_i + p_i c_i$ into itself gives each carry as a sum of products of $g$, $p$ and $c_0$:

$$c_1 = g_0 + p_0 c_0,\quad c_2 = g_1 + p_1 g_0 + p_1 p_0 c_0,$$
$$c_4 = g_3 + p_3 g_2 + p_3 p_2 g_1 + p_3 p_2 p_1 g_0 + p_3 p_2 p_1 p_0 c_0.$$

Each carry is a two-level AND-OR function, so all four carries of a 4-bit block are ready two gate delays after $g$, $p$ and $c_0$. With $g_3..g_0 = 0010$, $p_3..p_0 = 1101$ and $c_0 = 1$: $c_1 = 0 + 1 \cdot 1 = 1$, $c_2 = 1$ (from $g_1$), $c_3 = 0 + p_2 c_2 = 1$, $c_4 = 0 + p_3 c_3 = 1$. The carry vector $c_4 c_3 c_2 c_1 = 1111$.

## Blocks and their delay

The product terms grow quickly (the carry out of an $n$-bit block has $n + 1$ terms with up to $n + 1$ inputs), so lookahead logic is limited to blocks of about 4 bits. Each block also produces a group generate $G = g_3 + p_3 g_2 + p_3 p_2 g_1 + p_3 p_2 p_1 g_0$ and a group propagate $P = p_3 p_2 p_1 p_0$, so its carry-out is $G + P c_{in}$.

If blocks are rippled, a 16-bit adder of four 4-bit blocks has this delay under a simple model: $g, p$ take 1 gate delay; each block produces its carries 2 gate delays after its carry-in (its $G$ and $P$ are ready earlier); sums take 1 more XOR delay. Block 1 finishes at $1 + 2 = 3$, block $k$ at $2k + 1$, and the last sums at $2 \cdot 4 + 2 = 10$ gate delays, against $2 \cdot 15 + 2 = 32$ for a ripple adder with $d_c = d_s = 2$. Applying lookahead recursively to the block $G$ and $P$ signals gives logarithmic depth, which is the parallel-prefix idea.

## Going deeper

Weinberger and Smith, "A Logic for High-Speed Addition" (National Bureau of Standards Circular 591, 1958); Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 5.2.2.
