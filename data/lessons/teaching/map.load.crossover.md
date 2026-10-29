# Choosing matches under unknown load

During bottom-up mapping, a node's fanout in the final cover is not known, so neither is its load. A small, fast-starting cell wins at light load and a strong cell wins at heavy load. The mapper must guess, bin, or defer the choice. Knowing where the winner flips, and what a wrong guess costs, is the core of load-aware mapping.

## The crossover load

Two matches with delays $d_A + R_A C$ and $d_B + R_B C$, where $d_A < d_B$ and $R_A > R_B$, are equal at
$$C^* = \frac{d_B - d_A}{R_A - R_B}.$$
Below $C^*$ the weaker cell A is faster; above it the stronger B is. With $d_A$ = 8 ps, $R_A$ = 1.6 kOhm, $d_B$ = 15 ps, $R_B$ = 0.5 kOhm, the crossover is 7 / 1.1 = 6.4 fF. A library with several drive strengths of one function forms a lower envelope of such lines, and each cell owns an interval of loads.

## Best arrival at a known load

When the load is known, evaluate every match: arrival = latest input + intrinsic + R C, and keep the minimum. At 31 fF with inputs at 61 ps, the matches above give 61 + 8 + 49.6 = 118.6 ps and 61 + 15 + 15.5 = 91.5 ps, so B wins. Mappers that handle load explicitly keep, at each node, the best match for each of several load bins, a small table instead of one number; the consumer later picks the entry for the load it creates. The DP stays bottom-up but each state now carries a load index.

## The cost of a wrong guess

If the mapper assumes a load on one side of the crossover and the real load falls on the other side, it picks the wrong cell. Assume 12 fF with $d_A$ = 14 ps, $R_A$ = 2.2 kOhm, $d_B$ = 21 ps, $R_B$ = 0.8 kOhm: crossover 5 fF, so B is chosen. If the real load is 3 fF, B takes 21 + 2.4 = 23.4 ps while A would take 14 + 6.6 = 20.6 ps, a 2.8 ps penalty. Remedies are load binning, gain-based mapping that removes load from the delay model, iterating mapping with loads from the previous pass, and leaving the drive-strength choice to post-mapping sizing.

## Going deeper

Touati, "Performance-Oriented Technology Mapping", PhD thesis, UC Berkeley, 1990; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 10, on delay-oriented library binding.
