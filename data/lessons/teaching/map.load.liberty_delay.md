# Load-dependent match delay from Liberty

A real cell's delay depends on the load it drives and the slew it receives, as its Liberty tables say. A delay-oriented mapper must evaluate each candidate match with that dependence, even though the final load is unknown while mapping. This concept is the arithmetic the mapper does per match; the next concepts deal with the unknown load.

## Arrival through a match

With a linear model, the delay from input pin i to the output is $d_i + R\, C_{load}$: an intrinsic pin delay plus the drive resistance times the load. The arrival at the match output is
$$a_{out} = \max_i (a_i + d_i) + R\, C_{load},$$
where $C_{load}$ is the sum of the fanout pin capacitances plus wire capacitance. Example: inputs at 51, 51 and 61 ps with pin delays 23, 18 and 23 ps, R = 0.5 kOhm, load 2 + 5 + 6 + 9 = 22 fF: the latest input term is 61 + 23 = 84 ps and the load term is 11 ps, so the output arrives at 95 ps. Kilohms times femtofarads give picoseconds.

## Reading an NLDM table

Liberty stores delay as a table indexed by input slew and output load. Between grid points the timer and the mapper interpolate bilinearly. With slews $s_1 < s_2$, loads $c_1 < c_2$, and fractions $t_s = (s - s_1)/(s_2 - s_1)$, $t_c = (c - c_1)/(c_2 - c_1)$, the delay is
$$(1 - t_s)(1 - t_c) v_{11} + (1 - t_s) t_c v_{12} + t_s (1 - t_c) v_{21} + t_s t_c v_{22}.$$
For entries 21, 42, 34, 61 ps at slews 20, 60 ps and loads 2, 8 fF, evaluated at 56 ps and 2.6 fF: $t_s = 0.9$, $t_c = 0.1$, giving 35.34 ps. The output slew comes from a second table the same way and feeds the next stage.

## Load couples neighbouring choices

A match's input pins are the previous stage's load. Choosing a larger cell lowers its own drive resistance but raises its input capacitance, slowing its driver. Example: upsizing from input capacitance 3 to 9 fF and resistance 2 to 0.67 kOhm, driving 34 fF behind a 1.5 kOhm driver, changes the driver's delay by 1.5 x 6 = 9 ps and the gate's by about -45 ps, a net gain. With a light load the same swap would lose. This coupling is why mapping cannot be fully bottom-up with true delays: a node's best choice depends on its consumers, which are chosen later.

## Going deeper

The Liberty reference manual (Synopsys, open distribution) for NLDM tables; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8.
