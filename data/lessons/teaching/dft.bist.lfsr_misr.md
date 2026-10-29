# LFSRs and MISRs

Built-in self-test puts the pattern generator and the response checker on the chip. Both are linear feedback shift registers: an LFSR generates pseudorandom patterns, and a multiple-input signature register (MISR) compresses a long response stream into a short signature. They are tiny (one flip-flop per bit and a few XOR gates), run at full speed, and are analysable exactly with polynomial arithmetic over GF(2).

## The external-XOR LFSR

An $n$-stage **external-XOR** (Fibonacci) LFSR shifts $Q_1 \to Q_2 \to \cdots \to Q_n$ on each clock while $Q_1$ takes the XOR of a fixed set of tapped stages. With taps $Q_3 \oplus Q_4$ on 4 stages, the state 0010 has feedback $1 \oplus 0 = 1$, so the next state is 1001: the feedback enters $Q_1$ and the other bits move one place right. Repeating gives 1100, 0110, 1011, 0101. The **internal-XOR** (Galois) form puts XOR gates between stages instead; it generates the same sequences up to a change of basis and is faster because no gate chain feeds one flop. The all-zero state maps to itself, so LFSRs used as generators are seeded nonzero.

## The MISR

A **MISR** is the same register with one extra input per stage: on each clock, stage $i$ also XORs in response bit $d_i$. Starting from all zeros with taps $Q_1 \oplus Q_2 \oplus Q_4$ and input words 0001 then 0001: the first clock shifts 0000 to 0000 and adds 0001, giving 0001; the second has feedback $0 \oplus 0 \oplus 1 = 1$, shifts to 1000, and adds 0001, giving 1001. After the whole test the register holds the **signature**, which is compared with the fault-free signature computed by simulation. Because the MISR is linear, the signature of a faulty response equals the good signature XOR the signature of the error stream, which is what makes its analysis tractable.

## Polynomials over GF(2)

An LFSR's behaviour is captured by its **characteristic polynomial**. For the external-XOR form with taps $T$ (always including stage $n$), the output bit sequence satisfies $s_{k+1} = \sum_{t \in T} s_{k+1-t}$ over GF(2), with characteristic polynomial $x^n + \sum_{t \in T} x^{n-t}$; taps $\{3, 4\}$ give $x^4 + x + 1$. State sequences correspond to powers of $x$ modulo this polynomial, and signature analysis is polynomial division: for a single-input register in the internal-XOR form, the signature is the remainder of the response polynomial divided by the characteristic polynomial. The period, the aliasing behaviour and the randomness properties all follow from the polynomial.

## Going deeper

Bardell, McAnney and Savir, "Built-In Test for VLSI: Pseudorandom Techniques" (1987); Bushnell and Agrawal (2000), chapter 15; Golomb, "Shift Register Sequences" (1967).
