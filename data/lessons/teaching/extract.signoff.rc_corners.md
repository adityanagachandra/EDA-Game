# RC corners and signoff scenarios

Interconnect varies with the process just as transistors do: metal width, thickness and dielectric height all have tolerances. Extraction captures this with RC corners, sets of extreme wire geometries, and signoff combines them with transistor corners and functional modes into scenarios. Knowing which RC corner is worst for which net, and how many scenarios result, explains both the cost of signoff and why no single corner can stand for the rest.

## Counting scenarios

A scenario is one mode (functional, test, low-power) at one PVT corner (process, voltage, temperature) with one RC corner. Setup and hold are signed off at different corners: setup at slow transistors with Cworst and RCworst wiring, hold at fast transistors with Cbest and RCbest, among others. With 2 modes, setup at 3 PVT corners each with 3 RC corners, and hold at 2 PVT corners each with 2 RC corners, there are $2 \times 3 \times 3 + 2 \times 2 \times 2 = 26$ scenarios. Extraction runs once per RC corner and is shared across the PVT corners that use it, with temperature applied to resistance. Multi-mode multi-corner timers analyse many scenarios in one run, sharing the timing graph and keeping per-scenario arrival times.

## Corner geometry

RC corners move geometry to extremes. In the Cworst corner, wires are wider and thicker and dielectrics thinner, so capacitance, especially coupling, is highest, and resistance is somewhat lower. In RCworst, wires are narrower and thinner, so resistance is highest and capacitance somewhat lower. Resistance goes as $1/(WT)$: with width 8% and thickness 10% below nominal, a 340 ohm wire becomes $340/(0.92 \times 0.9) = 410.6$ ohm, 21% higher. Corners are derived from the foundry's statistical process data at a stated sigma level, and the technology file for each corner describes the stack with those dimensions. Because width and thickness move together in a corner, these are correlated extremes, not independent worst cases.

## Which corner is worst

A net's delay is roughly $0.69[R_d(C_w + C_L) + R_w(C_w/2 + C_L)]$. The driver term depends on capacitance only, the wire term on the product. Short, driver-dominated nets are worst at Cworst; long, resistive nets can be worst at RCworst. For a 0.8 kohm driver, 2.7 kohm of wire, 36 fF of wire and 18 fF of load, Cworst (R scaled 0.95, C 1.1) gives 98.7 ps and RCworst (R 1.15, C 0.96) gives 104.6 ps: this long net is RC-limited. A design contains both kinds of net, often on the same path, so signoff runs both corners, and a path can even be worst at a mixed condition no standard corner represents, which is one motivation for statistical extraction.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure; the OpenROAD documentation on multi-corner extraction with OpenRCX.
