# S-parameters of RF blocks

At radio frequencies, voltages and currents are hard to measure directly, but incident and reflected waves on matched lines are not. Scattering parameters describe a block by how it reflects and transmits those waves, and they are what network analyzers measure, what RF simulators report and what RF block specifications are written in: input match as $S_{11}$, gain as $S_{21}$, isolation as $S_{12}$.

## Reflection and return loss

A port terminated in impedance $Z_L$ on a reference line of impedance $Z_0$, usually 50 ohms, reflects a fraction $\Gamma = (Z_L - Z_0)/(Z_L + Z_0)$ of the incident wave. For a one-port this is $S_{11}$. A 100 ohm resistive input gives $\Gamma = 50/150 = 0.333$. Return loss is $-20\log_{10}|\Gamma|$: here 9.54 dB. A perfect match reflects nothing and has infinite return loss; an open or short reflects everything, 0 dB. LNA specifications usually ask for $|S_{11}|$ below -10 dB across the band, meaning at most about 10% of the incident power is reflected. Matching networks of inductors and capacitors transform the device's input impedance toward $Z_0$ over the band.

## VSWR

A mismatch sets up a standing wave on the line whose maximum and minimum voltages are $1 + |\Gamma|$ and $1 - |\Gamma|$ times the incident amplitude. Their ratio is the voltage standing-wave ratio, $\mathrm{VSWR} = (1 + |\Gamma|)/(1 - |\Gamma|)$. From $|S_{11}| = -10$ dB, $|\Gamma| = 10^{-0.5} = 0.316$ and VSWR = 1.92. A VSWR of 1 is a perfect match; 2 corresponds to about -9.5 dB return loss. Antenna and power-amplifier specifications are often written in VSWR because it bounds the voltage stress on the line. All three descriptions, $\Gamma$, return loss and VSWR, carry the same information about $|S_{11}|$; a tool should store one and convert.

## Transmission and two-port blocks

A two-port has four S-parameters: $S_{11}$ and $S_{22}$ are input and output reflections with the other port matched, $S_{21}$ is forward transmission and $S_{12}$ reverse transmission. Each is complex, with magnitude and phase. An amplifier's gain in a matched system is $|S_{21}|^2$ as a power ratio, $20\log_{10}|S_{21}|$ in dB: $S_{21} = 3 + j4$ has magnitude 5 and gain 13.98 dB; its phase, 53.1 degrees, is the insertion phase. $|S_{12}|$ measures reverse isolation, which matters for stability and for LO leakage in receivers. Simulators compute S-parameters from small-signal AC analysis with port terminations, and harmonic balance extends the idea to large signals. S-parameter files (Touchstone format) are also how passive structures such as inductors and packages enter circuit simulation.

## Going deeper

Pozar, "Microwave Engineering" (4th ed., 2012), chapter 4; Lee, "The Design of CMOS Radio-Frequency Integrated Circuits" (2nd ed., 2004).
