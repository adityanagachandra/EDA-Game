# Barrel and logarithmic shifters

Shifts and rotates by a variable amount appear in every ALU, in floating-point alignment and normalization, and in cryptographic datapaths. A shifter is a mux network, and its structure is a clean example of trading area against depth: a logarithmic shifter decomposes the shift amount into powers of two and applies one stage per bit.

## Shift and rotate semantics

For an $n$-bit word, a logical left shift by $s$ moves every bit up $s$ places and fills the bottom with zeros, which is multiplication by $2^s$ modulo $2^n$. A logical right shift fills the top with zeros. A rotate (circular shift) reinserts the bits that fall off one end at the other: rotate-left of $x$ by $s$ is $((x \ll s) \mid (x \gg (n - s))) \bmod 2^n$.

Rotate the 8-bit word 10110001 (177) left by 3: the three top bits 101 wrap to the bottom, giving $10001101 = 141$. A logical left shift by 3 instead gives $10001000 = 136$.

## A logarithmic shifter

Write the shift amount in binary, $s = \sum_j s_j 2^j$. A logarithmic shifter has $\lceil \log_2 n \rceil$ stages; stage $j$ either passes its input through or shifts it by $2^j$, depending on $s_j$. Each stage is a row of $n$ two-input muxes, so the whole shifter uses

$$n \lceil \log_2 n \rceil$$

muxes and has depth $\lceil \log_2 n \rceil$ mux delays. A 32-bit rotator takes $32 \cdot 5 = 160$ muxes in 5 levels. The alternative, a crossbar-style barrel shifter where each output picks directly from all $n$ inputs, has one level of $n$-input muxes ($n^2$ switches) and is preferred only when $n$ is small or the layout favors it.

## Which stages act

For a given shift amount, only the stages whose bit $s_j = 1$ actually shift; the others pass through. A shift of 13 in a 16-bit shifter ($13 = 1101_2$) uses the 8, 4 and 1 stages, three of the four.

Because rotations by different amounts compose by adding the amounts, the order of the stages does not affect the result, and designers place the stage whose select signal arrives latest closest to the output so that late select bits see the shortest path. For arithmetic right shifts the vacated positions fill with the sign bit, which adds a little logic at each stage; a shared shifter supporting logical, arithmetic and rotate modes computes a fill bit per stage and muxes it in. Synthesis of a Verilog `>>>` with a variable amount produces exactly this structure.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), chapter 11 on shifters; Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 5.
