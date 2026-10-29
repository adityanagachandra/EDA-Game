# State encoding

An FSM's states are abstract until they are assigned bit patterns. The encoding fixes the number of flip-flops, the complexity of the next-state and output logic, the switching activity, and how the machine behaves when it lands in an unused code. Synthesis tools such as Yosys and commercial compilers extract FSMs from RTL and re-encode them, so the trade-offs are part of the tool's own decision procedure.

## Flip-flop counts per encoding

For a machine with $S$ states:

- **Binary**: $\lceil \log_2 S \rceil$ flip-flops, states numbered $0, 1, \dots$.
- **Gray**: the same count, with codes taken from a Gray sequence so consecutive states differ in one bit.
- **One-hot**: $S$ flip-flops, exactly one set per state.
- **Johnson** (twisted ring): $\lceil S/2 \rceil$ flip-flops, since an $n$-bit Johnson counter cycles through $2n$ states.

For $S = 12$: binary and Gray need 4 flip-flops, one-hot 12, Johnson 6. One-hot spends flip-flops (cheap in FPGAs, where every logic cell has one) to make the next-state logic shallow: each state bit's next value is an OR of (predecessor bit AND condition) terms.

## Unused codes

An encoding with $b$ flip-flops has $2^b$ codes, of which $S$ are used. Binary with $S = 12$ leaves $16 - 12 = 4$ unused codes. One-hot with $S = 12$ leaves $2^{12} - 12 = 4084$: every pattern with zero or several bits set.

Unused codes matter because a radiation upset, a power-up without reset, or a bug can put the machine in one. Synthesis treats unreachable codes as don't-cares, which minimizes logic but leaves their behavior arbitrary; a "safe" FSM option instead adds logic so that every illegal code returns to a known state. One-hot machines are especially exposed, since most of their code space is illegal, and a single bit flip produces a two-hot pattern.

## Switching activity of an encoding

Each state transition toggles the flip-flops whose bits differ between the two codes: the Hamming distance. For a machine that steps cyclically through states $0, 1, \dots, S-1$ and back to 0, binary encoding toggles $\sum_{i} \text{popcount}(i \oplus ((i+1) \bmod S))$ bits per cycle of the loop. For $S = 8$ that is $1+2+1+3+1+2+1+3 = 14$ toggles per pass, where the last term is the wrap from 7 to 0. Gray encoding toggles exactly one bit per step when $S$ is a power of two, 8 per pass. One-hot toggles exactly two bits on every transition (one clears, one sets), 16 per pass.

Low-power state assignment chooses codes to minimize the probability-weighted Hamming distance over the transitions, which is a graph embedding problem, NP-hard in general and solved heuristically.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 9 on state encoding; the Yosys documentation for fsm_recode.
