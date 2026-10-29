# Effective capacitance of a resistive load

Liberty tables assume a cell drives a pure capacitor. A real routed net has resistance between the driver and the far capacitance, which shields part of it: during the fast output transition the driver does not see the full load. Effective capacitance is the single capacitance that makes the table give the right answer for the real resistive load. Every modern delay calculator computes it, iterating because the answer depends on the driver's own transition.

## Resistive shielding

Reduce the net to a pi model: $C_1$ at the driver pin, a resistance $R$, and $C_2$ at the far end. The far capacitor charges through $R$ with time constant $\tau = R C_2$. If the driver's output ramp is fast compared with $\tau$, the far capacitor barely charges during the ramp, and the driver sees mostly $C_1$. If the ramp is slow, the far capacitor keeps up and the driver sees nearly $C_1 + C_2$. With $R = 650$ ohm and $C_2 = 33$ fF, $\tau = 21450$ ohm-fF, or 21.45 ps. Using the total capacitance in the table would overestimate the delay at the driver pin, sometimes by a large margin on long resistive nets.

## Charge matching over a ramp

Model the driver as a ramp of duration $T$ at its pin. The near capacitor takes its full charge; the far capacitor, driven through $R$, reaches only the fraction $1 - (\tau/T)(1 - e^{-T/\tau})$ of the swing by the end of the ramp. Matching the total charge gives
$$C_{eff} = C_1 + C_2\left[1 - \frac{\tau}{T}\left(1 - e^{-T/\tau}\right)\right].$$
For $C_1 = 25$ fF, $C_2 = 54$ fF, $\tau = 21.6$ ps and $T = 171$ ps, $\tau/T = 0.1263$ and the far capacitor contributes 47.18 fF, so $C_{eff} = 72.18$ fF against a total of 79 fF. As $T/\tau$ grows, the bracket tends to 1 and $C_{eff}$ to the total; as it shrinks, $C_{eff}$ tends to $C_1$.

## Iterating with the cell table

The ramp time itself comes from the cell's table, evaluated at the load, which is $C_{eff}$: a fixed-point problem. Start from the total capacitance, look up the ramp time, compute a new $C_{eff}$, and repeat until the two agree. With a ramp time of $7 + 2C$ ps, starting at $C_{tot} = 37$ fF gives $T = 81$ ps, and with $\tau = 7.8$ ps the first update is $C_{eff} = 35.84$ fF. Since $C_{eff} \le C_{tot}$, the next ramp is shorter, which lowers $C_{eff}$ slightly again; a few iterations converge. Calculators in OpenSTA such as the dmp_ceff family fit a driver model to the table and refine this idea, and the waveform at the far end is then computed from the RC network rather than from the table.

## Going deeper

Qian, Pullela and Pillage, "Modeling the effective capacitance for the RC interconnect of CMOS gates" (IEEE TCAD, 1994); Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 5.
