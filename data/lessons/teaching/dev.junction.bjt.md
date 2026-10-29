# The bipolar junction transistor

A bipolar junction transistor is two back-to-back junctions sharing a thin base. Logic moved to CMOS long ago, but bipolar devices remain in every CMOS chip: as vertical PNPs in bandgap references and temperature sensors, and as the parasitic thyristor behind latch-up. Bipolar physics is also the cleanest example of an exponential transconductor, the behavior a MOSFET shows below threshold.

## Collector current from base-emitter voltage

In forward-active operation the base-emitter junction is forward biased and the base-collector junction reverse biased. Electrons injected from the emitter diffuse across the thin p-type base and are swept into the collector by the reverse-biased junction field. The collector current is therefore set by the base-emitter voltage alone:
$$I_C = I_S e^{V_{BE}/(kT/q)},$$
ignoring the Early effect (a weak dependence on collector voltage through base-width modulation).

Example: $I_S = 10^{-15}$ A and $V_{BE} = 0.65$ V at 300 K. $V_{BE}/(kT/q) = 25.15$, so $I_C = 10^{-15} \times e^{25.15} = 8.3 \times 10^{-5}$ A $= 83$ µA.

Every 59.5 mV of $V_{BE}$ multiplies $I_C$ by 10 at 300 K. Bandgap references exploit this exactly: two transistors at different current densities produce a $\Delta V_{BE} = (kT/q)\ln(\text{ratio})$ proportional to absolute temperature.

## Current gain and terminal currents

The base also draws a small current, mostly holes injected back into the emitter and some recombination in the base. Its ratio to the collector current is the current gain:
$$\beta = \frac{I_C}{I_B}, \qquad I_E = I_C + I_B = I_C\left(1 + \frac{1}{\beta}\right).$$
A thin, lightly doped base and a heavily doped emitter give high $\beta$, typically 50 to several hundred.

Example: $I_C = 2$ mA and $\beta = 100$. The base draws 20 µA, and the emitter carries 2.02 mA.

In bulk CMOS the same structures form unintentionally. A p+ source in an n-well over a p-substrate is a vertical PNP; an n+ source in the substrate next to the well forms a lateral NPN. Coupled together they make a pnpn thyristor. If both parasitic transistors turn on with loop gain $\beta_{npn}\beta_{pnp} > 1$, the structure latches into a low-impedance state between supply and ground. Well and substrate taps keep the base resistances low enough to prevent this.

## Transconductance

Differentiating the exponential gives the transconductance:
$$g_m = \frac{dI_C}{dV_{BE}} = \frac{I_C}{kT/q}.$$

Example: $I_C = 1$ mA at 300 K gives $g_m = 1/0.02585 = 38.7$ mS. At 400 K the same current gives 29.0 mS.

$g_m$ depends only on the bias current and temperature, not on device size or process. This is the largest transconductance per unit current that any device whose current follows a Boltzmann barrier can reach. A MOSFET in strong inversion has $g_m/I_D = 2/V_{ov}$, far smaller; in weak inversion it approaches $1/(n\,kT/q)$, which is the bipolar value divided by the slope factor $n$. That comparison, $g_m/I_D$ as a figure of merit, is the basis of the $g_m/I_D$ design method for analog CMOS and explains why subthreshold MOSFETs are efficient but slow.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), chapter 5; Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits" (5th ed.), chapter 1.
