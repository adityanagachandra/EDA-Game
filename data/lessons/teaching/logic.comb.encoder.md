# Encoders and priority encoders

An encoder is the inverse of a decoder: it turns one asserted line out of $2^n$ into its $n$-bit index. A plain encoder assumes exactly one input is active. Real request vectors have several bits set, so hardware uses priority encoders, which pick one by a fixed rule. They appear in interrupt controllers, arbiters, CAM match resolution and the leading-zero counters that normalize floating-point results.

## Highest index wins

A priority encoder with inputs $R_{n-1} \dots R_0$ outputs the index of the highest asserted input, plus a valid flag that is 0 when no input is asserted. For the 8-bit request vector $R_7 \dots R_0 = 00101100$, inputs 5, 3 and 2 are asserted; the output is 5.

Its logic is a chain of "nobody above me" terms: output $i$ wins when $R_i \cdot \bar{R}_{n-1} \cdots \bar{R}_{i+1}$. That chain is a prefix computation, so a wide priority encoder can be built in $O(\log n)$ levels with the same prefix trees used for fast adders.

## Lowest index wins

Priority can run the other way: the lowest asserted index wins, as in many fixed-priority arbiters where requester 0 is most important. For 00101100 the output is 2.

A neat identity finds the lowest set bit without a scan: $x \mathbin{\&} (-x)$ isolates it in two's complement, because negation inverts every bit above the lowest 1 and keeps that 1. For $x = 00101100$, $-x = 11010100$, and the AND is 00000100, which a plain encoder turns into 2. This is how a round-robin arbiter is often built: mask off requests below the last grant, then find the lowest remaining one.

## Leading-zero counting for normalization

After a floating-point subtraction, the result can have many leading zeros, and the hardware must shift it left until the top bit is 1. The shift amount is the leading-zero count, which is $w - 1 - p$ where $p$ is the index a highest-wins priority encoder returns for a $w$-bit word.

For the 16-bit word $0000\,0001\,0110\,1000$, the highest set bit is bit 8, so the leading-zero count is $15 - 8 = 7$. Fast adders often compute a leading-zero anticipation in parallel with the subtraction, predicting the count (sometimes off by one, corrected afterwards) rather than waiting for the sum. An equivalence checker comparing such a unit to a reference must cope with that internal speculation, which is one reason floating-point units are hard formal targets.

## Going deeper

Mano and Ciletti, "Digital Design" (5th ed., 2013), chapter 4; Muller et al., "Handbook of Floating-Point Arithmetic" (2nd ed., 2018), on leading-zero anticipation.
