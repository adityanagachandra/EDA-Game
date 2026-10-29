# Multiplier architectures

A multiplier forms partial products and adds them. Every architecture follows the same three steps: generate partial products, reduce them to two rows, and add the two rows with a fast carry-propagate adder. Architectures differ in how many partial products they create (Booth recoding) and how they reduce them (arrays, Wallace and Dadda trees). Multipliers dominate area in DSP and machine-learning accelerators.

## Partial products

Multiplying an $n$-bit unsigned multiplicand $A$ by an $m$-bit multiplier $B$ by long multiplication creates one row per multiplier bit: row $j$ is $A \cdot b_j$ shifted left by $j$, so its value is $b_j \cdot A \cdot 2^j$. Each partial-product bit is one AND gate, $a_i b_j$.

With $A = 1011$ (11) and $B = 1101$ (13), the rows are 1011 (row 0), 0 (row 1, since $b_1 = 0$), $1011 \ll 2 = 44$ (row 2) and $1011 \ll 3 = 88$ (row 3). Their sum is $11 + 44 + 88 = 143 = 11 \cdot 13$. The multiplication is now an addition problem: summing $m$ shifted rows.

## Sizes and widths

An $n \times m$ unsigned multiplier has $n m$ partial-product bits, so $n m$ AND gates, arranged in $m$ rows. The product needs $n + m$ bits: the largest product is $(2^n - 1)(2^m - 1) < 2^{n+m}$. A $16 \times 16$ multiplier has 256 partial-product bits and a 32-bit product.

For two's-complement operands the product of an $n$-bit and an $m$-bit signed number also fits in $n + m$ bits; the extreme case is $(-2^{n-1})(-2^{m-1}) = 2^{n+m-2}$, which needs the full $n + m$ bits as a positive signed value. Signed partial products require sign extension or the Baugh-Wooley rearrangement, which complements selected partial-product bits and adds constants so that every row stays positive.

## Array versus tree reduction

An array multiplier adds the rows one at a time with a grid of full adders, each row's adders taking the previous row's sums and carries. It is regular and easy to lay out, but its critical path grows linearly with the operand width, roughly two full-adder delays per bit for an $n \times n$ array, and its many unequal paths produce heavy glitching.

A tree multiplier instead reduces all rows in parallel with carry-save adders: a row of full adders takes three rows and produces two (a sum row and a carry row shifted left), without propagating any carry. The number of rows shrinks by a factor of about $3/2$ per level, so reduction takes $O(\log n)$ levels, followed by one fast carry-propagate adder. Wallace and Dadda trees are the two standard schedules for this reduction, and radix-4 Booth recoding halves the number of rows before it starts.

## Going deeper

Parhami, "Computer Arithmetic: Algorithms and Hardware Designs" (2nd ed., 2010), chapters 9 to 12; Baugh and Wooley, "A Two's Complement Parallel Array Multiplication Algorithm" (IEEE Transactions on Computers, 1973).
