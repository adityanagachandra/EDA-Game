# The FO4 delay as a process yardstick

The fanout-of-four (FO4) delay, an inverter driving four copies of itself, is the standard yardstick of a process's speed. Expressing a cycle time or a path delay in FO4 units removes the process from the comparison: a design with 20 FO4 per cycle is equally aggressive in any node. Library teams quote it, architects plan pipelines with it, and timing tools use it to normalize reports.

## Defining FO4 from stage delay

An inverter has $g = 1$. Driving four identical inverters gives $h = 4$. Its stage delay is
$$d_{FO4} = g h + p_{inv} = 4 + p_{inv},$$
which is 5 in units of $\tau$ when $p_{inv} = 1$. In absolute time, $t_{FO4} = (4 + p_{inv})\tau$. Fanout four is chosen because it is close to the stage effort that minimizes delay in a buffer chain, so an FO4 inverter is a representative well-sized stage. With $\tau = 6$ ps and $p_{inv} = 1$, the FO4 delay is 30 ps.

## Measuring tau from FO4

In practice FO4 is measured, often with a ring oscillator whose stages each drive a fanout of four, and $\tau$ is backed out. If the measured FO4 delay is 25 ps and $p_{inv} = 1.2$, then $\tau = 25/5.2 = 4.81$ ps. That $\tau$, together with the per-gate $g$ and $p$ values, lets a designer estimate any path in the process without running a simulator. Because $p_{inv}$ depends on the ratio of diffusion to gate capacitance, it is also fitted: one sweeps the fanout, plots inverter delay against $h$, reads $\tau$ from the slope and $p_{inv}\tau$ from the intercept. That is a straight-line fit of measured delays, exactly what a characterization flow does on a larger scale.

## Thinking in FO4 units

A cycle time divided by the FO4 delay gives the logic depth available per pipeline stage. A 500 ps clock in a process with a 20 ps FO4 has 25 FO4 per cycle, of which a few go to flop overhead and clock skew. Processors have been designed at roughly 10 to 30 FO4 per cycle, and the number is a meaningful comparison between designs in different nodes because it measures how much logic each stage holds, not how fast the transistors are. For a tool designer, FO4-normalized numbers make QoR comparisons across technologies and corners readable, and an FO4 estimate is a quick sanity check on any timer result: a 3 FO4 path reported at 400 ps in a 20 ps process means something is wrong.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.4; Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapter 1.
