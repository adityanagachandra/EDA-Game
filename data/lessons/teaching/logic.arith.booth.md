# Radix-4 Booth recoding

Radix-4 (modified) Booth recoding rewrites a two's-complement multiplier in base 4 with digits from $\{-2, -1, 0, 1, 2\}$. Each digit selects a partial product that is a trivial multiple of the multiplicand (zero, plus or minus itself, plus or minus twice itself), and the number of partial products halves. Nearly every high-performance multiplier uses it.

## Half as many partial products

A radix-2 multiplier creates one partial product per multiplier bit: $n$ rows for an $n$-bit multiplier. Radix-4 Booth examines the multiplier two bits at a time with one bit of overlap and creates one partial product per pair. For an $n$-bit two's-complement multiplier with $n$ even, that is $n/2$ partial products.

A $16 \times 16$ signed multiplier drops from 16 rows to 8; a $54 \times 54$ (double-precision significand) from 54 to 27. Each Booth row is a little wider, since it may hold $2A$, and needs a sign-handling bit, but halving the rows shrinks the reduction tree by about a level and saves substantial area.

## The recoding rule

Append $y_{-1} = 0$ below the multiplier's least significant bit. Digit $i$ ($i = 0, \dots, n/2 - 1$) looks at bits $y_{2i+1}, y_{2i}, y_{2i-1}$:

$$d_i = -2 y_{2i+1} + y_{2i} + y_{2i-1}.$$

So the triplet 000 or 111 gives 0, 001 and 010 give $+1$, 011 gives $+2$, 100 gives -2, and 101 and 110 give -1. Recode the 8-bit $y = 01101110$ (110): triplets from the bottom are $(y_1 y_0 y_{-1}) = 100 \to -2$, $(y_3 y_2 y_1) = 111 \to 0$, $(y_5 y_4 y_3) = 101 \to -1$, $(y_7 y_6 y_5) = 011 \to 2$. Check: $-2 + 0 \cdot 4 - 1 \cdot 16 + 2 \cdot 64 = 110$.

## Selecting and negating multiples

Digit $d_i$ selects the partial product $d_i \cdot A$, shifted left by $2i$ bits. With $A = 13$ and the digits above, the partial products are -26, 0, $-13 \cdot 16$ and $26 \cdot 64$; their sum is $13 \cdot 110 = 1430$.

Every multiple is cheap: $A$ is a wire, $2A$ a one-bit shift, and a negative multiple is the bitwise inverse plus one. The "plus one" is not added immediately; it is placed as an extra bit in the reduction tree so no carry chain is needed per row. The sign of each row is handled by a few extra bits rather than full sign extension. Booth's subtlety for verification is that the partial products no longer correspond to multiplier bits, so a word-level equivalence checker or an algebraic method must reason about the recoded sum rather than the original bit structure.

## Going deeper

Booth, "A Signed Binary Multiplication Technique" (Quarterly Journal of Mechanics and Applied Mathematics, 1951); MacSorley, "High-Speed Arithmetic in Binary Computers" (Proceedings of the IRE, 1961); Parhami, "Computer Arithmetic" (2nd ed., 2010), chapter 10.
