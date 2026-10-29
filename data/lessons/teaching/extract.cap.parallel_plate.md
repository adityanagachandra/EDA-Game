# Area capacitance to the layer below

The simplest capacitance a wire has is the parallel-plate part: the field running straight down from its bottom face to a wide conductor below, such as the substrate or a dense layer of crossing wires. It is the first term in every extraction model and the one a back-of-envelope estimate uses. Understanding how it scales with width, dielectric thickness and permittivity tells a tool designer why area capacitance alone has become a poor estimate at advanced nodes.

## The parallel-plate formula

Two parallel conducting plates of area $A$, separated by $h$ of a dielectric with relative permittivity $\varepsilon_r$, have $C = \varepsilon_0 \varepsilon_r A / h$, ignoring the field that bulges out at the edges. For a wire, $A = W L$. In micrometre units, $\varepsilon_0 = 8.854 \times 10^{-3}$ fF/um. A wire 1 um wide and 180 um long, 0.1 um above a ground plane in a dielectric with $\varepsilon_r = 2.5$, has $C = 0.008854 \times 2.5 \times 1 \times 180 / 0.1 = 39.84$ fF. Silicon dioxide has $\varepsilon_r \approx 3.9$; low-k interlayer dielectrics are lower, which is why they were introduced: capacitance falls in proportion.

## Capacitance per unit length

Extraction tables and wire-load estimates work per unit length, because a wire's cross-section repeats along it. Per micron, $C/L = \varepsilon_0 \varepsilon_r W / h$. A 1 um wire 0.05 um above its plane in $\varepsilon_r = 3$ gives $0.008854 \times 3 / 0.05 = 0.531$ fF/um, or 531 aF/um. Attofarads per micron (1 aF = $10^{-3}$ fF) are the natural unit for minimum-width wires, whose total capacitance per micron is typically of order 0.1 to 0.2 fF, that is 100 to 200 aF. Per-length values let an extractor multiply by the parallel run length of each fractured segment instead of solving the field again.

## Scaling and its limits

Area capacitance goes as $\varepsilon_r W / h$. If a shrink scales width by 0.5 and dielectric height by 0.8 with the same material, area capacitance per length falls to $0.5/0.8 = 0.625$ of its old value. That looks like progress, but the wire's sides have not shrunk as much as its width: thickness is kept up to control resistance, so the sidewalls become a larger share of the wire's surface. The parallel-plate term then underestimates the true capacitance badly, because fringe fields from the sidewalls and coupling to neighbours dominate. This is why total capacitance per unit length has stayed roughly constant across many nodes, and why extractors must model fringe and coupling explicitly.

## Going deeper

Cheng, Lillis, Lin and Chang, "Interconnect Analysis and Synthesis" (2000), chapter on capacitance; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 4.
