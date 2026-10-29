# The MOS capacitor and oxide capacitance

The MOS capacitor (gate, dielectric, semiconductor) is the heart of every MOSFET. The oxide capacitance per unit area $C_{ox}$ sets how much channel charge a volt of gate bias induces, and so sets drive current, gate capacitance and threshold sensitivity. Equivalent oxide thickness (EOT), the number every process quotes, is just $C_{ox}$ expressed as a thickness of SiO$_2$. Compact models, cell characterization and gate-load estimates all start here.

## Oxide capacitance per unit area

The gate and the semiconductor surface form a parallel-plate capacitor separated by the dielectric. Per unit area,
$$C_{ox} = \frac{\varepsilon_{ox}}{t_{ox}}, \qquad \varepsilon_{ox} = 3.9\,\varepsilon_0 = 3.453 \times 10^{-13}\ \text{F/cm}.$$

Example: $t_{ox} = 2$ nm $= 2 \times 10^{-7}$ cm gives $C_{ox} = 1.73 \times 10^{-6}$ F/cm$^2$. Since 1 F/cm$^2$ equals $10^7$ fF/µm$^2$, that is 17.3 fF/µm$^2$. A handy rule: $C_{ox} \approx 34.5/t_{ox}$ fF/µm$^2$ with $t_{ox}$ in nm.

The full gate-to-channel capacitance is less than $C_{ox}$, because the inversion charge sits a little below the interface (quantum confinement) and a polysilicon gate depletes. Both act as extra series capacitances and add a few tenths of a nanometre to the electrical thickness in inversion. Metal gates removed poly depletion; confinement remains, and compact models account for it.

## Equivalent oxide thickness

Modern gate stacks use high-k dielectrics such as HfO$_2$ (relative permittivity around 20 to 25) instead of pure SiO$_2$. To compare stacks, processes quote the equivalent oxide thickness: the SiO$_2$ thickness that would give the same capacitance,
$$\text{EOT} = t_{hk}\frac{3.9}{k_{hk}}.$$

Example: 3 nm of a $k = 22$ dielectric has EOT $= 3 \times 3.9/22 = 0.53$ nm. Real stacks keep a thin SiO$_2$ interfacial layer for interface quality, and layers in series add: EOT $= t_{IL} + t_{hk} \cdot 3.9/k_{hk}$. A 0.5 nm interlayer under the same high-k film gives 1.03 nm.

EOT is what transistor electrostatics care about: $C_{ox}$, inversion charge and short-channel control. Physical thickness is what tunneling cares about. Keeping those two separate is the whole point of high-k dielectrics, and compact models carry both (BSIM uses separate electrical and physical oxide thickness parameters).

## Inversion charge from gate overdrive

Above threshold, extra gate voltage induces inversion charge almost entirely in the channel, because the depletion charge is pinned once the surface is inverted. The charge-sheet approximation gives
$$Q_{inv} = C_{ox}(V_{GS} - V_t), \qquad N_{inv} = Q_{inv}/q.$$

Example: EOT $= 1$ nm ($C_{ox} = 3.45 \times 10^{-6}$ F/cm$^2$) and overdrive 0.5 V. $Q_{inv} = 1.73 \times 10^{-6}$ C/cm$^2$, so $N_{inv} = 1.08 \times 10^{13}$ electrons per cm$^2$.

Drive current is charge times velocity: $I_D/W = Q_{inv} v$. That one product explains why process scaling pushed EOT down (more charge per volt) and why strain engineering pushed mobility and velocity up. It also bounds what any channel material can do: at around $10^{13}$ cm$^{-2}$, even saturated velocity gives on the order of a milliampere per micron, which is where modern devices sit.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), chapter 2; Robertson, "High dielectric constant gate oxides for metal oxide Si transistors", Reports on Progress in Physics (2006).
