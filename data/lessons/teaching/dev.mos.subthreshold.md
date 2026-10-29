# Subthreshold conduction

Below threshold a MOSFET does not turn off; its current falls exponentially, by a factor of ten for every subthreshold swing $S$ of gate voltage. That exponential tail is the off-current of every idle transistor and the reason threshold voltage cannot scale freely. Leakage power analysis, threshold-flavor choice and near-threshold design all start from the subthreshold swing.

## The subthreshold swing

Below threshold, the channel is weakly inverted and current flows by diffusion over the source barrier. The barrier height follows the surface potential, which the gate controls through a capacitive divider, so the current is exponential:
$$I_D \propto e^{(V_{GS} - V_t)/(n\,kT/q)}.$$
The slope factor $n \ge 1$ measures how imperfectly the gate controls the surface. The gate voltage needed for one decade of current is the subthreshold swing
$$S = n\frac{kT}{q}\ln 10.$$

Example: $n = 1.2$ at 300 K gives $S = 1.2 \times 25.85 \times 2.3026 = 71.4$ mV/decade.

The ideal $n = 1$ limit is 59.6 mV/decade at 300 K. No device using thermionic injection over a barrier can beat it, which is the so-called Boltzmann tyranny. It scales with absolute temperature, so a hot die has a weaker turn-off.

## Current below threshold

With the swing known, current at any gate voltage below threshold follows from current at threshold:
$$I_D(V_{GS}) = I_D(V_t) \cdot 10^{-(V_t - V_{GS})/S}.$$

Example: 200 nA at $V_{GS} = V_t$ and $S = 80$ mV/decade. At $V_{GS} = V_t - 240$ mV the current is $200 \times 10^{-3} = 0.2$ nA $= 200$ pA. At $V_{GS} = 0$ with $V_t = 0.32$ V, it is $200 \times 10^{-4} = 20$ pA.

This is how a leakage tool reasons about off-current: each off transistor leaks the current at threshold scaled down by $V_t/S$ decades. A 10 mV threshold shift changes leakage by about 30% at 80 mV/decade, so threshold variation turns into a wide, log-normal-like leakage distribution across a die. Because the sum over millions of devices is dominated by the low-threshold tail, the mean leakage of a die is well above the leakage of a nominal device, which statistical leakage analysis must account for.

## Swing from the capacitive divider

The slope factor comes from how gate voltage divides between the oxide and the depletion region. With depletion capacitance $C_{dep}$ and interface-trap capacitance $C_{it}$,
$$n = 1 + \frac{C_{dep} + C_{it}}{C_{ox}}.$$

Example: $C_{ox} = 20$ fF/µm$^2$, $C_{dep} = 4$ fF/µm$^2$, negligible traps. $n = 1.2$ and $S = 71.4$ mV/decade at 300 K.

Thinner oxide (larger $C_{ox}$) improves swing, while heavier body doping (larger $C_{dep}$) worsens it. In short planar channels the drain also couples to the barrier, which degrades swing further. Thin-body devices such as FinFETs, nanosheets and FD-SOI have almost no depletion capacitance under gate control, so they approach $n = 1$, around 60 to 70 mV/decade. That electrostatic advantage, more than any change in mobility, is why the industry moved to them.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 3.1.3; Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), chapter 6.
