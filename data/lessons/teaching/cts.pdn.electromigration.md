# Electromigration in power grids

Electromigration (EM) is the slow transport of metal atoms by the electron flow. Over years it opens voids that raise resistance and eventually break wires, or pushes metal into extrusions that short neighbours. Power grid wires carry large, unidirectional currents, which makes them the primary EM risk. Foundries give current-density limits per layer and via; the PDN engine must size every strap and via array to meet them.

## Current density

The driver of EM is current density, $J = I/(wt)$, current divided by the wire's cross-section. A strap 6.5 µm wide and 0.85 µm thick carrying 45.5 mA has $J = 45.5/5.525 = 8.24$ mA/µm². Power-grid currents are unidirectional (DC-like), which is the worst case for EM, while signal nets carry bidirectional currents whose effect partly heals; signal EM is therefore checked with different, usually higher, limits on average, RMS and peak currents. Currents come from the IR solve: the drop difference across each branch divided by its resistance.

## Sizing wires to the limit

The foundry's EM rule gives a maximum $J_{max}$ for each layer at a reference temperature. The minimum width for current $I$ is $w = I/(J_{max}t)$: for 44.5 mA on a 0.65 µm thick layer with $J_{max} = 2.5$ mA/µm², $w \ge 27.4$ µm. Wider straps lower both EM risk and IR drop, so the two analyses usually agree; they disagree when a narrow wire carries current concentrated by a nearby strong supply path. EM rules are often given per unit width rather than per area; the check is the same once thickness is folded in.

## Vias, the weak point

Current crowds where a wire turns into a via, and via arrays are a common EM failure point. If each via is rated for 0.8 mA and a connection carries 71 mA, the array needs $\lceil 71/0.8 \rceil = 89$ vias. Real arrays share current unevenly, with the vias nearest the current source carrying more, so tools check each via individually from the solved grid. The limits fall with temperature, so a hot spot needs a larger margin. The lifetime model behind the limits is Black's equation.

## Going deeper

Black, "Electromigration: a brief survey and some recent results", IEEE Transactions on Electron Devices (1969); Lienig and Thiele, "Fundamentals of Electromigration-Aware Integrated Circuit Design" (2018).
