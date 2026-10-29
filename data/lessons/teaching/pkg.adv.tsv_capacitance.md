# TSV capacitance and delay

A TSV is a copper cylinder inside an oxide liner inside silicon, which makes it a coaxial capacitor to the substrate. That capacitance, typically tens of femtofarads, loads every signal that crosses a die boundary and enters timing, power and signal-integrity models. Extraction tools model it with the coaxial formula as a first estimate, refined for the depletion region in the silicon around the liner.

## The coaxial capacitance

For a conductor of radius $r$ inside a liner of outer radius $r + t_{ox}$ and length $L$, the coaxial formula gives $C = 2\pi \varepsilon_0 \varepsilon_{ox} L / \ln((r+t_{ox})/r)$. With $\varepsilon_0 = 8.854 \times 10^{-3}$ fF/um and $\varepsilon_{ox} = 3.9$, a TSV with $r = 2.5$ um, $t_{ox} = 0.2$ um and $L = 50$ um has $C = 2\pi \times 8.854 \times 10^{-3} \times 3.9 \times 50 / \ln(2.7/2.5) = 10.85/0.0770 = 141$ fF. This is the oxide capacitance only. In reality a depletion region forms in the silicon around the liner and adds a series capacitance, so the effective value is lower and depends on bias and doping, which is why 3D extraction uses voltage-dependent TSV models.

## TSV RC delay

A driver of output resistance $R_d$ charging the TSV plus a receiver load $C_L$ gives a 50% delay of about $0.69 R_d (C_{TSV} + C_L)$, because the TSV's own resistance of tens of milliohms is negligible. One kilo-ohm times one femtofarad is one picosecond. A 2 kilo-ohm driver into a 141 fF TSV plus 10 fF of load gives $0.69 \times 2 \times 151 = 208$ ps. That is slow compared with a gate delay, so TSV signals use stronger drivers, and dense TSV arrays use smaller diameters, which shrink both $r$ and $L$ and with them the capacitance. Hybrid-bond pads, with much smaller capacitance, make die crossings cheaper again.

## Liner thickness as a design lever

Capacitance depends on liner thickness through the logarithm. Doubling the liner from 0.2 to 0.4 um at $r = 2.5$ um changes $\ln(2.7/2.5) = 0.0770$ to $\ln(2.9/2.5) = 0.1484$, cutting capacitance by a factor of 0.519. A thicker liner lowers delay and switching energy, but it shrinks the copper core if the hole size is fixed, and thick liners are harder to deposit uniformly in deep holes. The logarithm also explains why shrinking the TSV's radius at fixed liner thickness increases the ratio $(r+t)/r$ and reduces capacitance per unit length. These sensitivities are exactly what a technology-exploration tool sweeps when it chooses a TSV process for a 3D stack.

## Going deeper

Katti, Stucchi, De Meyer and Dehaene, "Electrical modeling and characterization of through silicon via for three-dimensional ICs", IEEE Transactions on Electron Devices (2010).
