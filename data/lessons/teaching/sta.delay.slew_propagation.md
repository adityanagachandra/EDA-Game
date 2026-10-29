# Slew propagation: worst slew versus path slew

Delay depends on input slew, so every stage's output slew must be carried to the next stage. At a pin with several inputs, the timer has to decide which slew to pass on. Graph-based timers keep one slew per pin, the worst, and that choice is the main source of graph-based pessimism that path-based analysis later removes.

## Worst-slew merging at a pin

Consider a gate G with inputs A and B. The arc from each input has its own delay and its own output slew, both functions of that input's slew. A graph-based timer stores one late arrival and one late slew at G's output, each the worst over the input arcs: the arrival is $\max(AT_A + d_A, AT_B + d_B)$ and the slew is $\max(s_{out,A}, s_{out,B})$, computed separately. With arc slews of 124 ps from A and 103 ps from B, the timer stores 124 ps, even if the latest arrival came from B. The stored pair may not belong to any single real transition, but it bounds them all, because larger slews only increase downstream delays. OpenSTA follows this rule by default.

## Downstream arrival under graph-based analysis

The next gate H computes its delay from the merged slew. If H's arc delay is $a_H + k_H s$, the arrival at H's output is $AT_G + a_H + k_H \max(s_{out})$. For a G output arrival of 162 ps, a merged slew of 87 ps and $d_H = 36 + 0.3 s$, the arrival at H is $162 + 36 + 26.1 = 224.1$ ps. The latest arrival came through input A with a slew of 81 ps, but H was timed with B's 87 ps slew. Every downstream stage inherits the pessimism, and on deep paths through many merges it accumulates. A late slew that stems from a short, non-critical side input can inflate the delay of the critical path by several picoseconds.

## Path-based re-timing removes the pessimism

Path-based analysis re-times one specific path with the slews that actually travel along it. For the path through A, H's delay uses A's own output slew of 81 ps instead of 87 ps: $36 + 0.3 \times 81 = 60.3$ ps rather than 62.1 ps, so the arrival drops by 1.8 ps. In general the pessimism removed at one merge is $k_H (s_{max} - s_{path})$. Because PBA uses the real slew of the path being reported, its slack is always at least the GBA slack for that path. It is also much more expensive: it works one path at a time and cannot be updated in bulk, so tools use GBA for optimization and PBA to recover pessimism on the worst paths before signoff.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on path-based analysis; the OpenSTA documentation on slew propagation.
