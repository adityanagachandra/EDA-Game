# Core utilization, aspect ratio and die size

Before any cell is placed, the floorplan must decide how big the core is. Too small and the cells do not fit or cannot be routed; too large and the wires grow long and the die costs more. Core utilization, aspect ratio and the I/O margin turn the netlist's total cell area into a concrete core and die rectangle, and every placer and router inherits that choice.

## Utilization from cell area

Core utilization is the total standard-cell area divided by the core area. Designers choose a target, commonly somewhere around 50% to 80%, leaving whitespace for routing, buffers added during optimisation, clock tree cells, tap and filler cells. The required core area is the cell area divided by the target: 200000 um$^2$ of cells at 70% utilization needs $200000/0.7 \approx 285714.3$ um$^2$ of core. Macros are usually counted separately, by adding their area plus halos to the core. The target is a trade-off: higher utilization shortens wires and shrinks the die, while lower utilization eases routing and timing closure. A placer's own target density is a local version of the same idea, applied bin by bin.

## Aspect ratio and core dimensions

The aspect ratio fixes the shape: with ratio $R = H/W$ and area $A$, the core width is $W = \sqrt{A/R}$ and the height is $H = RW$. For $A = 285714.3$ um$^2$ and $R = 1$, both are about 534.5 um. For $R = 2$ the core is $\sqrt{285714.3/2} \approx 378.0$ um wide and 756.0 um tall. In practice the height is then rounded to a whole number of rows and the width to a whole number of sites. Square cores minimise the average Manhattan distance between random points and are the default; other ratios follow from the package, from neighbouring blocks in a hierarchical design or from the shape of large macros.

## From core to die

The die adds a margin around the core for the I/O ring, pad cells or bumps, power rings and seal ring. With a uniform margin $m$, the die is $(W + 2m)(H + 2m)$. A 534.5 um square core with a 20 um margin gives a die about 574.5 um on a side, 330050 um$^2$. Small blocks pay heavily for the margin: a 100 um core with the same margin loses about 31% of its die area to it. Because die area drives cost and yield, while utilization drives routability, an agent tuning a floorplan treats utilization, aspect ratio and margin as Tier 1 parameters and measures their effect with a fast placement and congestion estimate before committing.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3; the OpenROAD documentation for initialize_floorplan, which takes utilization, aspect ratio and core margins.
