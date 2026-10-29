# Parallel-prefix adders

Every carry is a prefix of the generate-propagate sequence combined by one associative operator. That observation turns carry computation into the parallel-prefix problem, solvable in logarithmic depth, and it produces the whole family of fast adders used in processors: Sklansky, Kogge-Stone, Brent-Kung, Han-Carlson, Ladner-Fischer. Datapath synthesis tools choose among these structures per adder.

## The prefix operator

Define the span signals $(G_{i:j}, P_{i:j})$: the span of bits $j$ to $i$ generates a carry, or propagates one through. Two adjacent spans combine with the operator

$$(G_{hi}, P_{hi}) \circ (G_{lo}, P_{lo}) = (G_{hi} + P_{hi} G_{lo},\; P_{hi} P_{lo}).$$

The operator is associative (though not commutative), so the span $i{:}0$ can be built from any bracketing of the single-bit spans $(g_k, p_k)$. With carry-in 0, the carry into bit $i + 1$ is $c_{i+1} = G_{i:0}$, and the sum is $s_i = p_i \oplus G_{i-1:0}$. Associativity is the whole trick: it lets a tree evaluate all $n$ prefixes in $\lceil \log_2 n \rceil$ levels instead of a chain of $n$.

## Computing the prefix generates

For a concrete pair of operands the prefix generates are just the carries. Take $a = 10110110$, $b = 01100101$ (8 bits) with carry-in 0. The carries computed column by column are $c_1 = 0$, $c_2 = 0$, $c_3 = 1$, $c_4 = 0$, $c_5 = 0$, $c_6 = 1$, $c_7 = 1$, $c_8 = 1$, so $G_{i:0} = 1$ for $i = 2, 5, 6, 7$: four of the eight prefixes.

The tree computes exactly these values, only faster. For example $G_{7:4}$ and $G_{3:0}$ combine as $G_{7:0} = G_{7:4} + P_{7:4} G_{3:0}$. Here $G_{3:0} = 0$ (no carry leaves the low nibble) but $G_{7:4} = 1$ (bit 5 generates and bits 6 and 7 propagate), so $G_{7:0} = 1$.

## Depth, size and fanout

Any prefix network on $n$ inputs needs at least $\lceil \log_2 n \rceil$ levels, since the top output depends on all $n$ inputs through two-input nodes. Networks differ in node count, fanout and wiring at a given depth.

The Sklansky (divide-and-conquer) network reaches minimum depth with $\frac{n}{2}\log_2 n$ nodes, but at its last level one node drives $n/2$ others, a large fanout that needs buffering. For $n = 32$ that is $16 \cdot 5 = 80$ nodes in 5 levels, with a fanout of 16 at the final level. Kogge-Stone keeps fanout at 2 by using more nodes and long wires; Brent-Kung uses fewer nodes at nearly twice the depth. Choosing among them is a trade between logic depth, area, fanout and wiring congestion, which is why it is decided late, with physical information.

## Going deeper

Ladner and Fischer, "Parallel Prefix Computation" (JACM, 1980); Harris, "A Taxonomy of Parallel Prefix Networks" (Asilomar, 2003); Zimmermann, "Binary Adder Architectures for Cell-Based VLSI and their Synthesis" (PhD thesis, ETH Zurich, 1997).
