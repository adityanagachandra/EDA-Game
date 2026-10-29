# Two's complement arithmetic

Hardware stores integers as fixed-width bit vectors, and almost every datapath interprets signed vectors in two's complement. The reason is economic: one binary adder then adds signed and unsigned numbers alike, because both are arithmetic modulo $2^n$. A tool designer meets two's complement everywhere: in the semantics a synthesizer must preserve for a Verilog `signed` declaration, in how an equivalence checker models a subtractor, and in every width and sign-extension bug a lint tool tries to catch.

## Weights and range

An $n$-bit two's-complement word $b_{n-1} \dots b_1 b_0$ has value

$$x = -b_{n-1} 2^{n-1} + \sum_{i=0}^{n-2} b_i 2^i.$$

The top bit carries a *negative* weight. Everything else is ordinary binary. The range is therefore asymmetric: from $-2^{n-1}$ (only the top bit set) to $2^{n-1} - 1$ (every bit except the top one set). An 8-bit word covers -128 to 127; a 12-bit word covers -2048 to 2047. There is exactly one zero, unlike sign-magnitude or one's complement, and half of the $2^n$ patterns are negative. The top bit doubles as the sign bit, so a comparator for "is negative" is a single wire.

## Decoding a pattern

To read a pattern, either use the weights directly or use the modular view: interpret the bits as an unsigned integer $u$, and if the top bit is set, subtract $2^n$. Take the 8-bit pattern `11101010`. As unsigned it is $128+64+32+8+2 = 234$. The top bit is set, so the signed value is $234 - 256 = -22$. The same answer from the weights: $-128 + 64 + 32 + 8 + 2 = -22$.

The modular view is the one hardware actually implements. Signed and unsigned values with the same bit pattern are congruent modulo $2^n$, so an adder that computes $u_a + u_b \bmod 2^n$ is simultaneously a correct signed adder whenever the true signed sum fits in range.

## Negation and the asymmetric corner

Negation is "invert every bit, then add one", because $\bar{x} = 2^n - 1 - u$, so $\bar{x} + 1 = 2^n - u \equiv -u \pmod{2^n}$. For $x = 22$ in 8 bits, `00010110` inverts to `11101001`, plus one gives `11101010`, which we just decoded as -22. As an unsigned integer that pattern is $256 - 22 = 234$.

The asymmetric range creates one corner case every verification engineer remembers: negating $-2^{n-1}$ gives $-2^{n-1}$ again. In 8 bits, `10000000` inverts to `01111111`, plus one is `10000000`. The true answer $+128$ is out of range, and the hardware silently wraps. An absolute-value unit, a negator or a saturating datapath must handle this pattern explicitly, and a formal property checker will find it if the designer did not.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 1; Knuth, "The Art of Computer Programming", volume 2, section 4.1.
