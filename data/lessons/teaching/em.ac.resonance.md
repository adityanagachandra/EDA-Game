# Resonance and quality factor

At resonance the inductive and capacitive reactances of a circuit cancel. A series resonance makes impedance minimal; a parallel resonance makes it maximal. Both appear constantly in EDA: package and decap anti-resonances set power-delivery noise, clock and oscillator tanks set frequencies, and parasitic resonances in long nets cause ringing. The quality factor measures how sharp the resonance is.

## Resonant frequency

The reactances $\omega L$ and $1/(\omega C)$ are equal when $\omega_0 = 1/\sqrt{LC}$, so the resonant frequency is
$$f_0 = \frac{1}{2\pi\sqrt{LC}}.$$
At $f_0$ a series LC has zero net reactance, and a series RLC presents only $R$; a parallel LC has zero net susceptance and, ideally, infinite impedance.

Example: $L = 10$ nH and $C = 10$ pF give $LC = 10^{-19}$ s$^2$, $\sqrt{LC} = 3.16 \times 10^{-10}$ s and $f_0 = 503$ MHz. Quadrupling $C$ halves $f_0$. A useful shortcut: one nanohenry times one picofarad is $10^{-21}$ s$^2$, so $f_0 \approx 5033/\sqrt{L_{nH} C_{pF}}$ MHz; for 10 nH and 10 pF that is $5033/10 = 503$ MHz. Resonance is a property of the pair, so the same frequency appears whether the elements are wired in series or in parallel.

## Quality factor and bandwidth

For a series RLC, the quality factor is $Q = \omega_0 L/R = \frac{1}{R}\sqrt{L/C}$, the ratio of reactance at resonance to resistance. It equals $1/(2\zeta)$. The 3 dB bandwidth of the resonance is $\Delta f = f_0/Q = R/(2\pi L)$.

Example: $R = 2$ ohm, $L = 10$ nH and $C = 10$ pF. $\sqrt{L/C} = 31.6$ ohm, so $Q = 15.8$; with $f_0 = 503$ MHz the bandwidth is $R/(2\pi L) = 31.8$ MHz. Doubling $R$ halves $Q$ and doubles the bandwidth. At resonance, the voltage across the capacitor is $Q$ times the source voltage, so high-Q parasitic resonances amplify noise. A high Q is desirable in an oscillator tank and dangerous in a power-delivery network, where designers deliberately add damping.

## Anti-resonance in a power delivery network

Seen from the die, package inductance $L$ (with a small series resistance $R$) is in parallel with on-die decoupling capacitance $C$. This parallel tank has an anti-resonance near $f = 1/(2\pi\sqrt{LC})$, where its impedance peaks. For high $Q$ the peak is
$$Z_{peak} \approx \frac{L}{RC} = Q^2 R.$$

Example: 0.5 nH of package inductance with 5 milliohm of resistance and 100 nF of decap. $f = 1/(2\pi\sqrt{5 \times 10^{-17}}) = 22.5$ MHz, $Q = \sqrt{L/C}/R = 0.0707/0.005 = 14.1$, and $Z_{peak} = 0.5 \times 10^{-9}/(0.005 \cdot 10^{-7}) = 1$ ohm. A 1 A current step with energy near that frequency would ring the supply by up to 1 V. Power-integrity tools sweep this impedance against a target impedance and add damping or more decap where it peaks.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems", chapters 1 and 2; Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 14.
