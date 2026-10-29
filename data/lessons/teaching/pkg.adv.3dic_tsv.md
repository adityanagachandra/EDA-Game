# 3D ICs and TSVs

A 3D IC stacks dies vertically and connects them through the silicon with through-silicon vias (TSVs): copper-filled holes insulated by an oxide liner. TSVs give vertical connections tens of micrometres long instead of millimetres of lateral wire, but they occupy area, stress the surrounding silicon and require thinned dies. Their geometry, aspect ratio and resistance are the first things a 3D floorplanner must know.

## Aspect ratio and die thickness

A TSV is etched by deep reactive-ion etching to a depth equal to the thinned die thickness $t$. Etch, liner deposition and copper fill are reliable only up to an aspect ratio $AR = t/d$ of roughly 10, sometimes more. So for a 5 um diameter TSV the die can be at most about 50 um thick, and a 100 um thick die needs TSVs of at least 10 um diameter. This couples three choices: thin dies allow small TSVs and dense vertical wiring but are fragile and need carrier wafers during processing; thick dies are robust but force large TSVs. Via-middle TSVs, formed after front-end transistors and before the back-end stack, are common for logic stacking.

## TSV resistance

A copper TSV is a cylinder of length $t$ and diameter $d$, so $R = \rho t/(\pi d^2/4)$. With $\rho = 1.7$ micro-ohm cm, which is 0.017 ohm um, a 50 um long, 5 um wide TSV has area 19.63 um^2 and $R = 0.017 \times 50/19.63 = 0.0433$ ohm, about 43 milliohm. That is small next to a millimetre of lateral wire, which is why vertical paths are attractive for power and signals. But a power network through hundreds of TSVs must still be modelled, because each TSV also carries a current limit for electromigration, and the current crowds where the TSV meets the thin back-end metal.

## TSV area overhead

Each TSV costs area on every die it passes through: the via itself plus a landing pad, and on the active die a keep-out zone. If $N$ TSVs each take a landing square of side $s$, the overhead fraction is $N s^2 / A_{die}$. Twenty thousand TSVs with 10 um landing pads on a 100 mm^2 die cost $20000 \times 100$ um^2 $= 2$ mm^2, or 2%. The overhead is why 3D partitioning tools minimise vertical connections and cluster them into arrays, why TSVs are placed in dedicated regions rather than scattered, and why a fine-grain partition that needs millions of vertical connections is only practical with hybrid bonding or monolithic 3D, not TSVs.

## Going deeper

Lim, "Design for High Performance, Low Power, and Reliable 3D Integrated Circuits" (2013); Garrou, Bower and Ramm (eds.), "Handbook of 3D Integration" (2008).
