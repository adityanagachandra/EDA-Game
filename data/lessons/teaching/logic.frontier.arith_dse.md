# Agent design-space exploration of arithmetic

Adders and multipliers have a rich, well-understood architecture space: ripple, carry-select, Sklansky, Brent-Kung, Kogge-Stone and hybrids for addition; array, Booth, Wallace, Dadda and compressor trees for multiplication. An agent exploring that space for a given width and timing target is a natural fast-and-accurate loop: analytic models prune, synthesis and placement decide. This concept is about running that loop correctly.

## Fast filters from closed forms

The closed forms from this domain are a fast model. For an $n$-bit adder: ripple has depth $n$ carry steps and $n$ cells; Brent-Kung $2\log_2 n - 1$ levels and $2n - 2 - \log_2 n$ nodes; Sklansky $\log_2 n$ levels and $\frac{n}{2}\log_2 n$ nodes with fanout $n/2$; Kogge-Stone $\log_2 n$ levels and $n\log_2 n - n + 1$ nodes. Given a delay per level and a budget, an agent can discard every architecture that cannot possibly meet timing and rank the survivors by node count, in microseconds.

The cheapest candidate that meets the budget on this model is the starting point, not the answer, because the model ignores wires, fanout loading and placement.

## Recursing to the accurate engines

The accurate evaluation is synthesis, placement and timing on the real library. It is expensive, so the agent spends it on the frontier: the few candidates within the fast model's error margin of each other, or of the budget. Kogge-Stone looks best on depth, but its long wires and high node count often lose to Sklansky or Han-Carlson once placed; Brent-Kung looks worst on depth but may meet a loose target at half the area.

The loop should record each fast-model estimate next to its accurate result. The gap between them, per architecture and width, calibrates the fast model, and a systematic gap (Kogge-Stone always slower than predicted) is exactly the knowledge worth keeping.

## Objectives and traps

The objective must match what the design needs: minimum area subject to timing, or minimum delay subject to area, or energy per operation at a throughput. An agent optimizing a single scalar (say area plus a weighted delay) may land on a point that violates a hard constraint the weight did not capture. Pareto fronts are safer: keep every candidate not dominated on (delay, area, power).

Two traps recur. Synthesis tools already choose adder architectures internally, so an explicit architecture in RTL may be re-optimized away; the agent must check what was actually built. And a comparison at one corner or one width does not transfer: the winner at 16 bits is often not the winner at 64.

## Going deeper

Harris, "A Taxonomy of Parallel Prefix Networks" (Asilomar, 2003); Zimmermann, "Binary Adder Architectures for Cell-Based VLSI and their Synthesis" (1997); Roy et al., "PrefixRL: Optimization of Parallel Prefix Circuits using Deep Reinforcement Learning" (DAC, 2021).
