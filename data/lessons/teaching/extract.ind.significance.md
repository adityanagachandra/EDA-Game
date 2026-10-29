# When on-chip inductance matters

Most on-chip nets are resistive enough that an RC model is accurate. Inductance matters on wide, low-resistance wires carrying fast edges: global clocks, wide buses and power distribution. Because full inductance extraction is expensive, flows use screening rules to decide which nets need it. These rules compare inductive and resistive impedance at the signal's significant frequency, and the signal's rise time with the wire's time of flight.

## The knee frequency

A digital edge with 10 to 90 percent rise time $t_r$ has most of its spectral energy below the knee frequency $f \approx 0.35/t_r$. This is the frequency at which to judge whether inductance matters. A clock edge of 170 ps has $f \approx 0.35/170\text{ ps} = 2.06$ GHz; a 25 ps edge has 14 GHz. The rule follows from the bandwidth of a single-pole system, whose rise time and 3 dB bandwidth satisfy $t_r f_{3dB} \approx 0.35$. Faster edges, not higher clock frequencies, are what bring inductance into play: a slow-clocked chip with sharp drivers still needs inductance analysis on its widest wires.

## Inductive versus resistive impedance

Compare $\omega L$ with $R$ at the knee frequency. If $\omega L \ll R$, the wire is resistance-dominated and RC suffices. For a 1500 um global wire with 0.2 ohm/um and 0.6 pH/um, and a 50 ps edge, $f = 7$ GHz, $\omega L = 2\pi \times 7 \times 10^9 \times 900 \times 10^{-12} = 39.6$ ohm against $R = 300$ ohm. The ratio is 0.13, so this wire is RC. The same length of a very wide top-metal wire with 0.01 ohm/um gives $R = 15$ ohm and a ratio of 2.6: inductive. Thin local wires have resistance of ohms per micron and are always RC. The loop inductance here assumes a known return path, so the screen is an estimate, used to decide where to spend accurate extraction.

## Time of flight and damping

A wire behaves as a transmission line when its signal changes faster than a wave crosses it. The time of flight is $t_f = l\sqrt{L'C'}$ for per-length inductance $L'$ and capacitance $C'$. With $L' = 0.5$ pH/um and $C' = 0.3$ fF/um, $\sqrt{L'C'} = 0.0122$ ps/um, so a 2500 um wire has $t_f = 30.6$ ps. A common screen flags the wire for an RLC model when $t_r$ is shorter than about $2t_f$ and its total resistance is below about $2Z_0$, where $Z_0 = \sqrt{L'/C'}$; here $Z_0 = 40.8$ ohm. A wire with more resistance than that is overdamped and shows no ringing. Both conditions together identify the nets where inductance changes delay and overshoot enough to matter.

## Going deeper

Ismail and Friedman, "Effects of inductance on the propagation delay and repeater insertion in VLSI circuits", IEEE Transactions on VLSI Systems (2000); Johnson and Graham, "High-Speed Digital Design" (1993), on the knee frequency.
