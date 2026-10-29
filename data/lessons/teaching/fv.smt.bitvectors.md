# SMT over bit-vectors

Hardware is made of fixed-width words, and the SMT theory of bit-vectors models them exactly: arithmetic wraps modulo $2^n$, shifts drop bits, and signed and unsigned operations are distinct. Word-level model checkers, equivalence checkers and symbolic simulators send their queries to bit-vector solvers, which mostly answer by bit-blasting to SAT. Evaluating bit-vector terms precisely and estimating their bit-blasted size are both everyday skills for a formal tool designer.

## Bit-vector terms and wraparound

In SMT-LIB, a bit-vector of width $n$ is a string of $n$ bits, and most operations act on the unsigned value modulo $2^n$. With 4-bit $a = 1101$ (13) and $b = 0110$ (6): bvadd gives $19 \bmod 16 = 3$; bvsub gives 7; bvmul gives $78 \bmod 16 = 14$. Shifts take a bit-vector shift amount: bvshl $a$ by 2 gives $0100 = 4$, since the top bits fall off; bvlshr $a$ by 1 gives $0110 = 6$, filling with 0; bvashr $a$ by 1 copies the sign bit and gives $1110 = 14$, the signed value -2. The extract operation ((_ extract 3 2) a) keeps bits 3 down to 2, $11 = 3$, and concat joins bit-vectors. Because all of this is exact, a bit-vector model captures overflow bugs that integer arithmetic would hide.

## Bit-blasting cost

Bit-blasting translates each operation into a gate-level circuit and then into CNF with the Tseitin encoding. A ripple-carry adder for $n$-bit bvadd uses a half adder at bit 0 (one XOR and one AND, $4 + 3 = 7$ clauses and 2 gate variables) and full adders at bits 1 to $n - 1$ (two XORs, two ANDs and one OR, $8 + 6 + 3 = 17$ clauses and 5 gate variables each). Clauses: $7 + 17(n - 1)$; for $n = 8$, $7 + 119 = 126$. Variables, counting the $2n$ operand bits: $2n + 2 + 5(n - 1)$; for $n = 8$, $16 + 2 + 35 = 53$. Addition is linear in width, but an $n$-bit multiplier is quadratic, which is why solvers simplify at the word level (constant multiplications, shifts, extracts of known bits) before bit-blasting, and why wide multipliers stay hard for SAT.

## Signed interpretation

The same bits have an unsigned and a signed (two's complement) reading. In $n$ bits the top bit has weight $-2^{n-1}$, so an 8-bit pattern ranges over -128 to 127 when signed. Arithmetic is identical for both readings; only the interpretation differs. Example: 8-bit $100 + 50 = 150$, the pattern 10010110, which reads as $150 - 256 = -106$ when signed: signed overflow, though the unsigned sum is fine. Comparisons are where the readings really diverge: bvult and bvslt give different answers for 10010110 against 00000001 (150 is not below 1 unsigned, but -106 is below 1 signed). Division and remainder also come in unsigned and signed versions, and right shifts in logical and arithmetic versions. Many RTL bugs are a signed comparison written as unsigned, and a bit-vector solver will find them exactly.

## Going deeper

Kroening and Strichman, "Decision Procedures" (2nd ed., 2016), chapter 6; the SMT-LIB standard, the FixedSizeBitVectors theory.
