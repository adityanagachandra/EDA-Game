# Trace resistance from copper weight

A trace is a resistor of rectangular cross-section. Its resistance sets the DC voltage drop on power traces, the loss at low frequency, and the self-heating that limits current. Board tools compute it from the stackup's copper weight and the trace geometry, and a power-integrity tool sums it over thousands of shapes. The model is $R = \rho L / (w t)$ with the resistivity of copper, corrected for temperature.

## Resistance from geometry

For a uniform conductor, $R = \rho L / A$ with $A = w t$. Copper's resistivity is about $1.72 \times 10^{-8}\ \Omega\cdot$m at 20 C. The thickness $t$ comes from the copper weight: 1 oz is about 35 micrometres, 2 oz about 70. Example: a trace 100 mm long and 0.25 mm wide in 1 oz copper has $A = 0.25 \times 10^{-3} \cdot 35 \times 10^{-6} = 8.75 \times 10^{-9}$ square metres, so $R = 1.72 \times 10^{-8} \cdot 0.1 / 8.75 \times 10^{-9} \approx 0.197\ \Omega$, about 197 milliohms. A handy form is sheet resistance: 1 oz copper is about 0.5 milliohm per square, so the trace above, $100/0.25 = 400$ squares long, has about 200 milliohms. Unit discipline is the main source of error: convert every length to metres first.

## Voltage drop on power traces

Ohm's law turns resistance into the drop a load sees. A 0.2 ohm trace carrying 2 A drops 0.4 V, which would be fatal for a 1.0 V core rail with a 3% tolerance (30 mV). That is why power is distributed on planes or wide pours, and why regulators use **remote sense**, measuring the voltage at the load and adjusting their output to cancel the drop. On a plane, current spreads, and the resistance between a regulator and a BGA depends on the plane's geometry and the vias that feed it; a DC power-integrity tool solves a resistive mesh, a sparse linear system of nodal equations, to find the drop at every pin. A single trace is the one-dimensional special case.

## Temperature and frequency effects

Copper's resistance rises with temperature by about 0.393% per C: $R_T = R_{20}(1 + 0.00393 (T - 20))$. A trace at 85 C has about 26% more resistance than at 20 C, so a drop budget must be checked hot, and a trace that heats itself raises its own resistance, a feedback that current-capacity rules account for. At high frequency the **skin effect** confines current to a layer near the surface, about 2 micrometres thick at 1 GHz in copper, so a 35 micrometre trace uses only a fraction of its cross-section and its AC resistance grows roughly with the square root of frequency. Surface roughness adds more loss. DC resistance is therefore a lower bound for signal loss, accurate only for power and slow signals.

## Going deeper

Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic" (1993), chapter 4, on resistance and skin effect. IPC-2152, on conductor heating.
