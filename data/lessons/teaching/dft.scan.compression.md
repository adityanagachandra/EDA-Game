# Scan compression

Tester channels are few and expensive, while designs have millions of scan flops. Scan compression puts a decompressor between the tester's scan-in channels and many short internal chains, and a compactor between the internal chains and the scan-out channels. Test data volume and shift time fall by roughly the ratio of internal chains to channels, typically ten to a hundred times, because ATPG patterns specify only a few percent of their bits.

## Architecture

With $F$ flops and $c$ tester channels, uncompressed scan drives $c$ chains of length $\lceil F/c \rceil$. A compressed design builds $N \gg c$ internal chains of length $\lceil F/N \rceil$. Each shift cycle the tester supplies $c$ bits, and a linear **decompressor** (a ring generator or LFSR plus a phase shifter of XOR gates, in EDT-style schemes) expands them into $N$ bits. On the output side, a **compactor** of XOR trees folds $N$ chain outputs into $c$ tester outputs. With 1,036,342 flops, 16 channels and 1600 internal chains, the shift length per pattern drops from 64,772 to 648, about a factor of 100.

## Data volume and ratio

Stimulus volume is patterns times channels times shift cycles. With the same pattern count, the compression ratio is the ratio of shift lengths, approximately $N/c$. For 11,146 patterns, 16 channels and 160 internal chains on 456,058 flops, it is $28{,}504/2{,}851 = 10.0$. Compression works because ATPG cubes are sparse: a typical pattern specifies 1 to 5 percent of its bits as care bits, and the decompressor only has to match those. Response compaction reduces output volume by the same ratio.

## Pattern inflation

The decompressor's outputs are linear functions of fewer free variables, so some cubes cannot be encoded: too many care bits, or a linearly dependent combination. ATPG with compression then splits work across more patterns, so the compressed pattern count is higher. With 1076 uncompressed patterns and 12 percent inflation to 1206, a raw chain ratio of 20 becomes an effective ratio of about 17.8. Unknown values (X) in responses are the other limit, because an X entering a compactor corrupts its outputs; X-masking logic and X-tolerant compactors address them.

## Going deeper

Rajski, Tyszer, Kassab and Mukherjee, "Embedded deterministic test," IEEE Transactions on CAD, 2004; Touba, "Survey of test vector compression techniques," IEEE Design and Test, 2006.
