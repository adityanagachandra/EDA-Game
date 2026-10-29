# Leakage currents

An off transistor still conducts. Subthreshold current, gate tunneling, junction leakage, band-to-band tunneling and gate-induced drain leakage add up to a static power that can rival switching power in idle blocks. Leakage drives threshold-flavor choice, power gating, standby modes and corner selection. Power analysis tools read it from Liberty tables that compact models produce, so a tool designer should know what each component is and what it depends on.

## The leakage components

Four mechanisms dominate an off MOSFET.

Subthreshold leakage flows through the channel, by diffusion over the source barrier, and scales as $10^{-V_t/S}$. It usually dominates in logic, especially for low-threshold devices at high temperature.

Gate leakage is direct tunneling through the dielectric. It fell sharply with high-k dielectrics, which allow a thicker physical layer at the same capacitance.

Junction leakage is reverse-bias current of the source and drain junctions: diffusion and generation currents, plus band-to-band tunneling where steep, heavily doped junctions (such as halo regions) create high fields.

Gate-induced drain leakage (GIDL) arises when the gate is low and the drain high: the strong field in the gate-drain overlap bends the bands enough for band-to-band tunneling, and current flows from drain to body. It grows at negative $V_{GS}$, which limits how far an off device can be pushed below threshold.

## What leakage depends on

Temperature: subthreshold leakage rises steeply as threshold falls and swing grows with temperature; a factor of several between 25 °C and 85 °C is common. Gate tunneling is nearly temperature independent.

Circuit state: two off transistors in series leak much less than one. The intermediate node rises slightly, giving the upper device negative $V_{GS}$, a body-effect threshold increase and reduced DIBL. This stack effect makes cell leakage depend on input state.

Body bias: reverse body bias raises threshold and cuts subthreshold leakage, but increases the junction field, so band-to-band and GIDL leakage rise. Beyond some bias, total leakage increases again.

Mitigation at the design level uses higher-threshold cells off the critical path, power gating with high-threshold header or footer switches that disconnect idle blocks, and retention schemes for state. Clock gating removes dynamic power but not leakage.

## How EDA models leakage

Liberty files give each cell a leakage power value per input state, because the inputs decide which transistors are off and whether they form stacks. Power tools weight these states by signal probabilities from simulation or propagation to estimate average leakage.

Leakage varies by an order of magnitude or more across process corners and temperature, so leakage numbers are only meaningful with their corner. A fast corner at high temperature is the worst-case leakage corner, while the same corner is the best case for delay.

Leakage also couples to temperature. Leakage heats the die and heat raises leakage. If leakage power rises faster with temperature than the package can remove heat, the loop runs away. Thermal-aware power analysis iterates temperature and leakage to a fixed point, and packaging decisions depend on the slope of that loop. For a tool designer this is a classic coupled nonlinear problem solved by fixed-point iteration.

## Going deeper

Roy, Mukhopadhyay and Mahmoodi-Meimand, "Leakage current mechanisms and leakage reduction techniques in deep-submicrometer CMOS circuits", Proceedings of the IEEE (2003); Rabaey, "Low Power Design Essentials" (2009).
