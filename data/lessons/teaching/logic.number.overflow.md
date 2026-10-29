# Overflow, sign extension and shifts

Fixed-width arithmetic wraps. Whether a wrap is a bug depends on the interpretation: an unsigned counter is meant to wrap, a signed accumulator usually is not. Datapath designers detect overflow, widen operands by sign extension before mixing widths, and divide by powers of two with arithmetic shifts. Each of these is a frequent source of silent RTL bugs, and each is something a synthesis or equivalence tool must model exactly as the language defines it.

## Wrapped sums and signed overflow

An $n$-bit adder returns $(a + b) \bmod 2^n$ read back in two's complement. When the true sum lies outside $[-2^{n-1}, 2^{n-1}-1]$, the result differs from it by exactly $2^n$: that is signed overflow. In 8 bits, $100 + 50 = 150$ wraps to $150 - 256 = -106$.

Overflow can only happen when both operands have the same sign, and it shows up as a result of the opposite sign. In hardware the cheap test is $V = c_{n} \oplus c_{n-1}$: the carry into the top bit differs from the carry out of it. Adding two negatives, $-100 + (-50) = -150$, wraps to $-150 + 256 = 106$, positive, so $V = 1$. Adding operands of opposite signs never overflows, no matter how large they are.

## Sign extension

To widen an $n$-bit signed value to $m > n$ bits without changing it, copy the sign bit into all the new upper positions. The pattern `1010` (value -6 in 4 bits) becomes `11111010` in 8 bits, still -6: the new bits contribute $-2^{7} + 2^{6} + 2^{5} + 2^{4} = -2^{4}$, and the old top bit, now an ordinary $+2^{3}$, brings the total for the upper part to $-2^{3}$, exactly the weight that bit carried before. Read as unsigned, the 8-bit pattern is 250.

Zero extension (filling with zeros) is right for unsigned values and wrong for negative signed ones: `00001010` is $+10$. In Verilog, mixing a signed and an unsigned operand makes the whole expression unsigned, so the narrower operand is zero-extended even if it was declared signed. Lint tools and equivalence checkers flag exactly this.

## Arithmetic shifts round toward minus infinity

An arithmetic right shift by $s$ shifts the bits right and refills the top with copies of the sign bit. Its value is $\lfloor x / 2^s \rfloor$, the floor, not the truncation toward zero that integer division in C performs. For positive $x$ the two agree. For negative $x$ they differ whenever bits are shifted out: $-13 \gg 2$ gives $\lfloor -3.25 \rfloor = -4$, while C's $-13/4$ is -3.

Check in 8 bits: -13 is `11110011`; shifting right twice with sign fill gives `11111100`, which is -4. A hardware divider by a constant power of two that must match software division adds a bias of $2^s - 1$ to negative operands before shifting. Forgetting that bias is a classic mismatch between a C reference model and the RTL.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 5; IEEE Std 1800-2017 (SystemVerilog), the rules for expression signedness and width.
