# Logical effort

Logical effort measures how much worse a gate is than an inverter at driving a load, purely because of its topology. It turns gate sizing into arithmetic on a few numbers per cell, so it is how designers reason about paths and how many sizing and buffering heuristics in synthesis tools estimate delay before they have tables. Its definition rests on equal-drive sizing.

## Definition by input capacitance

Size a gate so that each input's worst-case pull-up and pull-down resistance equals a reference inverter's. The logical effort $g$ of an input is the ratio of its input capacitance to the inverter's input capacitance at that drive:
$$g = \frac{C_{in,\text{gate}}}{C_{in,\text{inv}}}.$$
Because both gates now drive a load equally well, $g$ says how many times more input capacitance the gate presents for the same output current. Equivalently, it is how much slower the gate is than an inverter at driving the same multiple of its own input capacitance. An inverter has $g = 1$ by definition. Logical effort depends on topology and on $\gamma$, not on the absolute size of the gate, which is why one number per input pin suffices.

## NAND and NOR formulas

With pMOS $\gamma$ times wider than nMOS for equal drive, the reference inverter's input capacitance is $1 + \gamma$. A NAND$n$ input drives an nMOS of width $n$ and a pMOS of width $\gamma$, so
$$g_{\text{NAND}n} = \frac{n + \gamma}{1 + \gamma}.$$
A NOR$n$ input drives an nMOS of width 1 and a pMOS of width $n\gamma$, so
$$g_{\text{NOR}n} = \frac{1 + n\gamma}{1 + \gamma}.$$
For $\gamma = 2$: NAND2 $4/3$, NAND3 $5/3$, NOR2 $5/3$, NOR3 $7/3$. For $\gamma = 3$ the NOR penalty is larger: NOR2 is $7/4 = 1.75$ while NAND2 is $5/4 = 1.25$. This is the quantitative reason libraries favour NAND-based logic and why synthesis tools choose NAND-heavy mappings in timing-critical cones.

## Compound gates and per-pin effort

For a compound gate, size by longest series path and divide each input's capacitance by $1 + \gamma$. With $\gamma = 2$, an AOI21 $\overline{AB + C}$ has nMOS widths 2, 2, 1 and pMOS widths 4, 4, 4, so inputs A and B have $g = 6/3 = 2$ and input C has $g = 5/3$. Different pins of one cell can thus have different logical efforts, and the faster pin should receive the critical signal. Logical effort grows with the number of inputs and with series stacking, which explains why very complex cells stop paying off: their per-input effort climbs faster than the stages they save. The model is linear and fitted at typical slews, so it is a sizing and reasoning tool, not a signoff timer.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapters 1 and 4; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.5.
