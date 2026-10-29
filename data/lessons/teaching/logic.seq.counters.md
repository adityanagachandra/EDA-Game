# Ripple versus synchronous counters

A counter is the simplest sequential circuit with interesting timing. Built asynchronously, as a ripple counter, each stage is clocked by the previous stage's output; built synchronously, every bit shares the clock and an incrementer computes the next count. The comparison is a compact lesson in why synchronous design won: the ripple counter is smaller, but its outputs are skewed and its settle time grows with width.

## The ripple counter

In an $n$-bit ripple counter each bit is a toggle flip-flop. Bit 0 is clocked by the clock; bit $i$ is clocked by the falling of bit $i-1$ (for an up-counter built from positive-edge flops, by $\bar{Q}_{i-1}$ rising). A bit toggles only when the bit below it falls from 1 to 0, which is exactly binary carry propagation.

The worst case is a count like $0111\ldots1 \to 1000\ldots0$, where every bit toggles in sequence, each one clock-to-q after the previous. The last bit settles after $n \cdot t_{cq}$. For an 8-bit counter with $t_{cq} = 80$ ps that is 640 ps, and in between, the outputs pass through intermediate wrong values, a burst of glitches on any logic that decodes them.

## Reading a ripple counter safely

If logic clocked by the main clock samples a ripple counter, it must wait for the worst-case settle time plus its own setup time, so

$$T_{\min} = n \cdot t_{cq} + t_{su}.$$

With $n = 8$, $t_{cq} = 80$ ps and $t_{su} = 40$ ps, $T_{\min} = 680$ ps, $f_{\max} \approx 1470.6$ MHz, and it falls further as the counter widens. A synchronous counter instead computes $Q^+ = Q + 1$ with an incrementer and loads all bits on the same edge: its period is $t_{cq} + t_{inc} + t_{su}$, where the incrementer can use a fast prefix carry chain whose delay grows as $\log n$. All its outputs change together, one clock-to-q after the edge.

## Modulus and state count

A mod-$N$ counter cycles through $N$ states. In binary it needs $\lceil \log_2 N \rceil$ flip-flops and leaves $2^{\lceil \log_2 N \rceil} - N$ codes unused. A mod-10 counter needs 4 flip-flops and has 6 unused codes, which the design must handle (by reset, or by making each unused code lead back into the cycle) in case a soft error or a bad power-up lands there.

Ripple counters still appear as clock dividers and in low-power timers, where only the top bit is used as a slower clock. EDA tools must then treat each divided output as a generated clock with its own latency, which is exactly the kind of structure that complicates timing constraints and clock-domain-crossing analysis.

## Going deeper

Mano and Ciletti, "Digital Design" (5th ed., 2013), chapter 6; Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 5 on counters.
