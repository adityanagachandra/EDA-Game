# Short-circuit power

During an input transition there is an interval when both the pull-up and pull-down conduct, and current flows straight from the supply to ground. This short-circuit or crowbar power is usually a small part of dynamic power, but it grows quickly with slow edges. It is why libraries measure internal power against input slew, and one reason transition limits exist.

## When both networks conduct

In an inverter the nMOS conducts once the input exceeds $V_{tn}$, and the pMOS until the input exceeds $V_{DD} - |V_{tp}|$. For inputs in between, both are on. As the input ramps through that window, a current spike flows from $V_{DD}$ to ground in addition to the current that charges or discharges the load. The energy lost per transition grows with the time spent in the window, so with the input transition time. If $V_{DD} < V_{tn} + |V_{tp}|$ the window does not exist and there is no short-circuit current at all, which is one minor benefit of very low supply voltages. With a static input, one device is always off and only leakage flows.

## Input slew versus output slew

The crowbar current depends on the voltage across the device that is turning off, and that is set by how fast the output moves. If the output is heavily loaded and moves slowly, then during the input transition the output is still near its starting rail, the turning-off device has little drain-source voltage, and little crowbar current flows. If the output is lightly loaded and fast, the output swings while both devices are on, and crowbar current is larger. The result is that short-circuit energy depends on the ratio of input to output transition times: matched slews keep it a small fraction of switching energy, while a slow input into a fast, lightly loaded gate is the worst case. Doubling the input slew roughly doubles the overlap time.

## How tools account for it

Because short-circuit energy depends on input slew and output load, it cannot be computed from capacitance alone. Characterization measures the total energy drawn from the supply during each transition, subtracts the energy that charges the external load, and records the remainder, which includes crowbar current and internal node charging, in Liberty internal_power tables indexed by input transition and output load. A power tool then looks up that energy per transition and adds the net's switching power separately. In a power report, short-circuit power therefore appears inside internal power, not switching power. Design rules that limit maximum transition keep input edges sharp and keep this component small, which is part of why buffering slow nets saves power as well as time.

## Going deeper

Veendrick, "Short-circuit dissipation of static CMOS circuitry and its impact on the design of buffer circuits" (IEEE JSSC, 1984); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 5.5.
