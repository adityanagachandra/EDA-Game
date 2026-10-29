# Subtractors and comparators

An adder becomes a subtractor by inverting one operand and setting the carry-in, because $a - b = a + \bar{b} + 1$ modulo $2^n$. The carry-out and the overflow and sign flags of that subtraction then answer every comparison, signed or unsigned. ALUs implement comparison this way, and so do the condition codes a processor's branch logic reads.

## Subtraction with an adder

Since $\bar{b} + 1 \equiv -b \pmod{2^n}$, feeding $\bar{b}$ into an adder with carry-in 1 computes $a - b \bmod 2^n$. An adder-subtractor puts an XOR gate on each $b$ input, controlled by a subtract signal that also drives the carry-in.

For 8-bit $a = 75$, $b = 200$: $\bar{b} = 55$, and $75 + 55 + 1 = 131$. As an unsigned result that is $(75 - 200) \bmod 256 = 131$. A signed reading sees the same input pattern of $b$ as -56 and the same output pattern as $131 - 256 = -125$, the wrapped result of $75 - (-56) = 131$, a signed overflow. The bit pattern is the same either way; only the interpretation, and the flags, differ.

## Unsigned comparison from the carry-out

The carry-out $c_n$ of $a + \bar{b} + 1$ is 1 exactly when $a \ge b$ as unsigned numbers. The reason: the adder computes $a + (2^n - 1 - b) + 1 = 2^n + (a - b)$, which reaches $2^n$ (setting the carry-out) iff $a - b \ge 0$.

So the carry-out is an inverted borrow: $c_n = 0$ means the subtraction borrowed, $a < b$. For $a = 75$, $b = 200$ the true sum is $75 + 55 + 1 = 131 < 256$, so $c_n = 0$ and $a < b$. Processor flag conventions differ on whether they store $c_n$ or its complement as the carry flag after a subtraction, a detail an instruction-set simulator and the RTL must agree on.

## Signed comparison from N and V

For signed operands the sign of the result is not enough, because the subtraction can overflow. The rule is

$$a < b \iff N \oplus V = 1,$$

where $N$ is the result's sign bit and $V$ is the signed-overflow flag. Take 8-bit $a = 100$, $b = -50$: $a - b = 150$ overflows to -106, so $N = 1$ but $V = 1$, and $N \oplus V = 0$, correctly saying $a \ge b$. Without the overflow correction the sign bit alone would claim $100 < -50$.

An equality comparator needs no subtraction at all: $a = b$ iff every bit of $a \oplus b$ is 0, an XOR row followed by a NOR tree of depth $\lceil \log_2 n \rceil$.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 5.2.4 on ALUs and flags; Parhami, "Computer Arithmetic: Algorithms and Hardware Designs" (2nd ed., 2010), chapter 2.
