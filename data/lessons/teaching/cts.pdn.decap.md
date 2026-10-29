# Decoupling capacitance

Decoupling capacitors are local charge reservoirs between VDD and VSS. When a burst of switching draws more current than the grid and package can deliver in time, the decap supplies the charge and its voltage falls a little. Sizing decap is the charge balance $Q = C\,\Delta V$. A PDN engine places decap cells in free sites, and an IR ECO loop adds them near hot spots, so the sizing arithmetic is used constantly.

## Charge balance

If a block's current rises by $\Delta I$ for a time $\Delta t$ before the grid responds, the decap must supply $Q = \Delta I\,\Delta t$ while its voltage drops by at most $\Delta V$, so
$$C \ge \frac{\Delta I\,\Delta t}{\Delta V}.$$
Units: milliamps times nanoseconds give picocoulombs, and picocoulombs over millivolts give nanofarads. For 1750 mA over 0.7 ns with a 65 mV budget, $Q = 1225$ pC and $C = 18.85$ nF. Real response times come from the grid and package impedance, so $\Delta t$ is itself an estimate.

## Charge from a current pulse

Switching current is a pulse, not a step. The charge is the integral of the current, so a triangular pulse peaking at $I_p$ over width $t_w$ carries $\tfrac{1}{2}I_p t_w$. A 200 mA peak over 1.7 ns carries 170 pC, which with a 35 mV budget needs $170/35 = 4.86$ nF. Libraries characterise each cell's switching current waveform, and power tools sum them per region and per clock edge to find the worst charge demand.

## Intrinsic and explicit decap

Not all decap is added on purpose. The gate and junction capacitance of cells that are not switching, and the grid's own capacitance, act as intrinsic decap. Explicit decap cells, usually large MOS capacitors, cover the shortfall. A block needing 3.5 nF that has 1.75 nF intrinsic needs 1.75 nF explicit; at 20 fF/µm² that is 87 500 µm². Decap is not free: thin-oxide decap leaks, it takes area and placement sites that later timing ECOs might need, and too much on-die capacitance with package inductance creates a resonance that must be damped.

## Going deeper

Popovich, Mezhiba and Friedman, "Power Distribution Networks with On-Chip Decoupling Capacitors" (2008); Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems" (2007).
