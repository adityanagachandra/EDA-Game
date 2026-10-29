# Sheet resistance and wire resistance

Wire resistance comes from resistivity, geometry and the way layout is drawn. Because a metal layer has a fixed thickness, extractors work with sheet resistance and count squares. Modern copper lines are so narrow that surface and grain-boundary scattering raise their effective resistivity well above bulk, and vias add a resistance that often dominates short routes.

## Resistivity and sheet resistance

A uniform conductor of length $L$, width $W$ and thickness $t$ has $R = \rho L/(Wt)$. For a layer of fixed thickness, define the sheet resistance $R_s = \rho/t$, in ohm per square. Then $R = R_s L/W$, which depends only on the drawn shape.

Example: copper with an effective resistivity of 2.2 microohm cm ($2.2 \times 10^{-8}$ ohm m) and thickness 200 nm has $R_s = 2.2 \times 10^{-8}/2 \times 10^{-7} = 0.11$ ohm per square. Bulk copper is 1.72 microohm cm, but narrow lines have higher effective values because electrons scatter at surfaces and grain boundaries and because the barrier liner carries little current. Technology files therefore give width-dependent resistivity or sheet resistance tables rather than one number.

## Counting squares

The number of squares in a wire is $L/W$: a wire 100 um long and 0.1 um wide is 1000 squares regardless of its absolute size. Its resistance is $R_s$ times the square count.

Example: 0.08 ohm per square, 0.05 um wide and 200 um long gives 4000 squares and 320 ohm. Doubling the width halves the resistance and roughly doubles the area capacitance, which is the basic sizing trade for wires. Extractors fracture routed polygons into rectangles, count squares in each, and add corrections for bends and width changes, where current crowds at the inner corner. The resistance network they build feeds IR-drop analysis and RC delay calculation.

## Vias and parallel paths

A route that changes layers passes through vias, each with its own resistance, often a few ohms in advanced processes. Several via cuts in parallel divide that resistance by their number, and the layer segments add in series. Upper layers are usually thicker, with lower sheet resistance, which is why long routes are promoted upward.

Example: 100 um on a lower layer at 0.1 ohm per square and 0.05 um wide is 200 ohm; a via array of 2 cuts at 4 ohm each adds 2 ohm; 500 um on an upper layer at 0.02 ohm per square and 0.4 um wide adds 25 ohm, for 227 ohm in total. On a short route the via share dominates. Redundant vias also improve yield and electromigration margin. An extractor must model via arrays explicitly; lumping them hides the single-cut resistance a single failing via would expose.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), chapter 6 on interconnect; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8.
