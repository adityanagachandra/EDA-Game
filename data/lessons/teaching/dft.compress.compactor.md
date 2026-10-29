# Response compactors and X-masking

On the output side of a compressed design, many internal chains feed a few tester channels through a compactor. A space compactor is a network of XOR trees: each output is the parity of a group of chain outputs. It loses information in two ways, aliasing, where errors cancel, and X-corruption, where an unknown value poisons an output. Compactor design is about controlling both while keeping the output count small.

## XOR space compaction

Each shift cycle the chains present one bit each, and output $z_j$ is the XOR of the bits of the chains in its group. With chains c1 to c9 presenting 011001001 and groups $z_1 = $ (c4, c5, c7), $z_2 = $ (c2, c6, c9), $z_3 = $ (c1, c3, c8), the outputs are $z_1 = 0 \oplus 0 \oplus 0 = 0$, $z_2 = 1 \oplus 1 \oplus 1 = 1$, $z_3 = 0 \oplus 1 \oplus 0 = 1$, so 011. A single error in any chain flips its group's output, so single errors are always observed. That is the essential property: a defect usually corrupts one or a few cells.

## Aliasing

Two errors in the same cycle cancel if they fall in the same XOR group, because XOR of two flips is no flip. With group sizes $g_j$, the number of chain pairs that alias is $\sum_j \binom{g_j}{2}$. For 10 chains in groups of 4, 3 and 3, that is $6 + 3 + 3 = 12$ of the 45 pairs. Compactors reduce aliasing by connecting each chain to several outputs with distinct patterns, so the columns of the connection matrix differ and any two errors produce a nonzero output vector; codes with larger minimum distance tolerate more simultaneous errors. Errors in different cycles are seen separately, since a space compactor has no memory, unlike a MISR.

## Unknown values and masking

An X from an uninitialized flop, a bus with no driver or a timing-uncertain path enters an XOR and makes its output X, so every chain in that group becomes unobservable in that cycle. With X on c3 and c8, the groups containing them are lost: in a 10-chain, 4-output compactor that can leave only 5 chains observable. **X-masking** logic gates the chain carrying the X to a constant before the compactor, controlled per pattern or per cycle by mask data from the tester, so the other chains in the group remain observable. X-tolerant compactors spread each chain over several outputs so that an X hides fewer chains. Masking costs data and can hide real errors on the masked chain, so the first remedy is design: eliminate X sources.

## Going deeper

Mitra and Kim, "X-compact: an efficient response compaction technique," IEEE Transactions on CAD, 2004; Rajski et al., "Embedded deterministic test," 2004.
