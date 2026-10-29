# Redistribution layers

A redistribution layer is a thin copper wiring level, insulated by polymer or inorganic dielectric, that moves connections from where the die's pads are to where the bumps or balls must be. RDLs appear in fan-in and fan-out WLP and on interposers. They are lithographically defined at a few micrometres line and space, far finer than a laminate substrate but far coarser than on-chip metal, and their resistance and density trade off directly.

## RDL trace resistance

An RDL trace of length $L$, width $w$ and thickness $t$ has $R = \rho L/(w t)$. With $\rho$ in micro-ohm cm, note that 1 micro-ohm cm is $10^{-2}$ ohm um, so with lengths in micrometres $R = 0.01 \rho L/(w t)$ ohms. Plated RDL copper is a little above bulk, around 2 micro-ohm cm. A 5 mm trace 5 um wide and 3 um thick has $R = 0.01 \times 2 \times 5000/15 = 6.67$ ohm. That is fine for a short signal into a high-impedance receiver but large for a long one, and the RC product with a few picofarads of load reaches tens of picoseconds. RDL thickness is the lever: thicker copper lowers resistance but makes fine line and space harder to pattern.

## IR drop on RDL power routes

Power often also travels through the RDL from a ball to the die's power pads, and the drop is $V = I R$. A power route of 2 ohm carrying 150 mA drops 300 mV, far too much for a 0.8 V supply whose budget might be a few percent. Designers therefore use wide RDL planes or meshes for power, many parallel paths, and place power balls close to the pads they feed. Because RDL copper is thin compared with substrate copper, an RDL-only package such as a substrate-less fan-out is more sensitive to IR drop and electromigration than a laminate package, and its power design is checked with the same grid analysis used on chip.

## Density against resistance

RDL density is $1000/(w+s)$ lines per millimetre per layer. Moving from a 10/10 um substrate-class pitch to a 2/2 um RDL gives $20/4 = 5$ times more wires per layer, so an edge that needed five substrate layers needs one RDL layer. The price is resistance: at fixed thickness-to-width ratio, resistance per unit length grows as the inverse square of the width. Fine RDL is therefore used for short, dense die-to-die wiring, while long or high-current routes stay on coarser layers. Multi-layer RDL stacks commonly combine a few fine layers near the die with coarser layers toward the balls, the same hierarchy logic as an on-chip metal stack.

## Going deeper

Lau, "Semiconductor Advanced Packaging" (2021), chapter on RDLs; Lau, "Fan-Out Wafer-Level Packaging" (2018).
