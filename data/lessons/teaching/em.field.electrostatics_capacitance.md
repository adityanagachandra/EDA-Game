# Electrostatics and capacitance

Capacitance is a geometric property of conductors and the dielectric between them: charge per unit voltage. Every interconnect capacitance a timer uses comes from solving electrostatics, either exactly by a field solver or approximately by formulas fitted to one. The parallel-plate result is the anchor of all of them, and stacked dielectrics show how real process cross-sections combine.

## Field between parallel plates

Between two large parallel plates held $V$ apart at spacing $d$, Gauss's law with a curl-free field gives a uniform field $E = V/d$ normal to the plates, ignoring the fringing near their edges. The surface charge density on each plate is $\sigma_s = \varepsilon E$.

Example: 1.2 V across 2 nm of gate dielectric gives $E = 1.2/(2 \times 10^{-9}) = 600$ MV/m, which is 0.6 V/nm. Fields of this size are why dielectric breakdown and time-dependent dielectric breakdown limit supply voltages for thin oxides, and why reliability tools check field strength rather than voltage alone. Units are easy to slip: one volt per nanometre is 1000 MV/m, and one volt per micrometre is 1 MV/m.

## Parallel-plate capacitance

The charge on a plate of area $A$ is $Q = \sigma_s A = \varepsilon A V/d$, so
$$C = \frac{\varepsilon_0\varepsilon_r A}{d}.$$
In on-chip units, $\varepsilon_0 = 8.854 \times 10^{-3}$ fF/um, so a square micrometre of SiO2 ($\varepsilon_r = 3.9$) 1 um thick is 0.0345 fF.

Example: a 10 um by 20 um plate 0.2 um above ground in SiO2 gives $C = 8.854 \times 10^{-3} \cdot 3.9 \cdot 200/0.2 = 34.5$ fF. The formula ignores fringing, so it underestimates real wires, badly once the width is comparable to the height. Capacitance scales linearly with area and permittivity and inversely with spacing, which is why low-k dielectrics and thicker inter-layer dielectrics both reduce wire capacitance.

## Layered dielectrics

A process cross-section stacks dielectrics: an etch-stop or capping layer, a low-k bulk film, a barrier. For layers parallel to the plates, $D$ is continuous across each interface, so each layer drops voltage $D d_k/\varepsilon_k$, and the layers act as capacitors in series:
$$C = \frac{\varepsilon_0 A}{d_1/\varepsilon_{r1} + d_2/\varepsilon_{r2} + \cdots}.$$

Example: a 10 um by 10 um plate over 0.2 um of $\varepsilon_r = 2.7$ and 0.05 um of $\varepsilon_r = 7.0$: the denominator is $0.0741 + 0.0071 = 0.0812$ um, so $C = 8.854 \times 10^{-3} \cdot 100/0.0812 = 10.9$ fF. Without the thin layer the plate would see 11.95 fF through the low-k film alone; adding it in series lowers the total by about 9%, even though it is only a fifth of the thickness. Extractors carry the full stack in their technology file for exactly this reason: an effective single permittivity is only correct for one geometry.

## Going deeper

Griffiths, "Introduction to Electrodynamics" (4th ed.), chapters 2 and 4; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 4 on wire capacitance.
