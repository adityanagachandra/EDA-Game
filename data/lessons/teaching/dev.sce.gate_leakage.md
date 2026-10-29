# Gate tunneling leakage and high-k dielectrics

As SiO$_2$ gate dielectrics thinned toward 1.2 nm, direct tunneling current through them became a significant part of standby power. High-k dielectrics solved this by making the layer physically thicker at the same capacitance. Computing physical thickness, tunneling reduction and chip-level gate current shows why the change was worth new materials and new compact-model terms, and why gate leakage still appears in leakage reports.

## Physical thickness at a given EOT

Capacitance depends on $\varepsilon/t$, so a dielectric with relative permittivity $k$ reaches the capacitance of SiO$_2$ of thickness EOT at a physical thickness
$$t_{hk} = \text{EOT}\cdot\frac{k}{3.9}.$$

Example: EOT $= 0.9$ nm with HfO$_2$-like $k = 20$ needs $t_{hk} = 0.9 \times 20/3.9 = 4.6$ nm, about five times thicker than the equivalent SiO$_2$.

In practice the stack includes a thin SiO$_2$-like interfacial layer, often a few tenths of a nanometre, which uses up part of the EOT budget at low permittivity. The high-k layer must then supply the rest. High-k materials also bring problems a device engineer must manage: more interface and bulk traps (which feed PBTI), remote phonon scattering that lowers mobility, and threshold shifts from charge and dipoles at the interfaces, which is part of why metal gates and work-function engineering arrived together with high-k.

## How much leakage high-k saves

Direct tunneling scales as $e^{-2\kappa t}$, and $\kappa$ depends on the barrier height and tunneling mass. High-k dielectrics have lower barriers than SiO$_2$ (HfO$_2$ around 1.5 eV versus about 3.1 eV for electrons), so $\kappa$ is smaller, but the thickness gain dominates.

Example with equal prefactors: SiO$_2$ at 1.0 nm with $\kappa = 6.4$ nm$^{-1}$ gives an exponent of 12.8. A $k = 20$ layer at the same EOT is 5.13 nm thick; with $\kappa = 3.0$ nm$^{-1}$ its exponent is 30.8. The difference, 18.0, is $18.0/2.303 = 7.8$ decades.

This estimate is crude: prefactors, interfacial layers and image-force barrier lowering all matter. But it captures why the industry accepted a materials change: orders of magnitude of gate leakage at the same capacitance. Compact gate-current models keep the exponential thickness dependence, with parameters fitted to measured gate current per area.

## Gate leakage of a whole chip

Gate current scales with gate area, so chip-level gate leakage is current density times total biased gate area:
$$I_g = J_g \cdot N \cdot \overline{W L}.$$

Example: $10^9$ transistors with average gate area $0.1 \times 0.02$ µm$^2$ give $2 \times 10^6$ µm$^2$ $= 0.02$ cm$^2$. At $J_g = 1$ A/cm$^2$, the chip draws 20 mA of gate current, a large standby budget for a mobile chip.

Two details matter for a power tool. Gate leakage flows only where the gate is biased relative to the channel or source and drain, so it depends on logic state, as subthreshold leakage does. And it is nearly independent of temperature, unlike subthreshold leakage, so its share of total leakage is largest at cold corners. Liberty leakage values include it through characterization, but separating the components is what lets a designer reason about which mitigation (power gating, flavor swaps, voltage reduction) will help.

## Going deeper

Wilk, Wallace and Anthony, "High-k gate dielectrics: current status and materials properties considerations", Journal of Applied Physics (2001); Robertson, Reports on Progress in Physics (2006).
