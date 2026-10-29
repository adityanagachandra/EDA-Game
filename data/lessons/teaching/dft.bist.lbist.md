# Logic BIST and test points

Logic BIST applies pseudorandom patterns from an on-chip generator to the scan chains and compacts the responses in a MISR, with no tester data beyond a seed and a signature. It enables in-field and power-on self-test, now required in safety-critical parts. Its weakness is random-pattern-resistant faults, which are attacked by inserting test points or weighting the pseudorandom patterns.

## STUMPS and test time

The standard architecture is **STUMPS** (self-testing using MISR and parallel shift register sequence generator). An LFSR feeds the internal scan chains through a **phase shifter**, an XOR network that decorrelates adjacent chains so that they do not receive shifted copies of the same sequence. All chains shift in parallel; their outputs feed a MISR. Each pattern shifts for the longest chain's length while the previous response shifts into the MISR, then captures. A session of $P$ patterns on chains of length $L$ takes $P(L + 1)$ cycles: 1000 patterns at $L = 106$ and 100 MHz take 1.07 ms. BIST sessions use far more patterns than ATPG, often tens of thousands, because pseudorandom patterns are inefficient.

## Test points

A fault excited with probability $p_e$ and propagated with probability $p_p$ has detection probability $p_e p_p$, and needs about $\ln(1/0.05)/(p_e p_p)$ patterns for 95 percent confidence. An **observation point** taps the hard-to-observe line into a scan cell, making $p_p = 1$. With $p_e = 2^{-8}$ and $p_p = 2^{-9}$, the pattern count falls from about 392,656 to 766. A **control point** inserts an AND or OR gate driven by a test signal to force a hard-to-control line, raising $p_e$. SCOAP or probabilistic measures pick the locations; each point costs area and possibly timing, so tools insert as few as reach the coverage target.

## Weighted random patterns

Instead of changing the circuit, bias the patterns. If a cone needs many inputs at 1, drive each to 1 with probability $w > 0.5$: the probability that all $k$ inputs are 1 rises from $0.5^k$ to $w^k$. For $k = 5$ and $w = 0.8$, by a factor $(0.8/0.5)^5 = 10.49$. The same weight makes faults that need a 0 on those inputs harder, so weight sets are computed per region or per session. BIST must also guarantee a known signature, which means every X source (uninitialized memories, multicycle paths, bus contention) must be blocked; X-bounding is often the dominant LBIST design effort.

## Going deeper

Bardell and McAnney, "Self-testing of multichip logic modules," International Test Conference, 1982; Wang, Wu and Wen (2006), chapter 5; Bushnell and Agrawal (2000), chapter 15.
