# Phasors and sinusoidal steady state

In sinusoidal steady state every voltage and current in a linear circuit oscillates at the source frequency, differing only in amplitude and phase. Phasors capture those two numbers as one complex number, which turns differential equations into algebra. AC analysis in SPICE, S-parameter extraction and power-delivery impedance plots are all phasor computations, so phasor arithmetic is the working language of frequency-domain engines.

## Sinusoids as rotating complex numbers

Write $v(t) = V_m\cos(\omega t + \phi) = \mathrm{Re}\{V_m e^{j\phi} e^{j\omega t}\}$. The complex number $\mathbf{V} = V_m e^{j\phi}$, written $V_m\angle\phi$, is the phasor; the factor $e^{j\omega t}$ is common to every signal and is dropped. Differentiation becomes multiplication by $j\omega$, and adding same-frequency sinusoids becomes adding complex numbers.

Example: $3\cos(\omega t)$ plus $4\cos(\omega t + 90^\circ)$ mA. The phasors are 3 and $4j$, whose sum is $3 + 4j$ with magnitude 5 mA and angle $53.1^\circ$. Adding the amplitudes directly would give 7 mA, which is wrong unless the phases match. KCL and KVL hold for phasors exactly as they hold for instantaneous values, because they are linear.

## Element impedances

Substituting phasors into the element laws gives impedances, the ratio $\mathbf{V}/\mathbf{I}$: a resistor has $Z_R = R$, an inductor $Z_L = j\omega L$ and a capacitor $Z_C = 1/(j\omega C) = -j/(\omega C)$. Their magnitudes $\omega L$ and $1/(\omega C)$ are the reactances, with $\omega = 2\pi f$.

Example: at 100 MHz a 10 pF capacitor has $|Z_C| = 1/(2\pi \cdot 10^8 \cdot 10^{-11}) = 159.2$ ohm, and a 10 nH inductor has $|Z_L| = 2\pi \cdot 10^8 \cdot 10^{-8} = 6.28$ ohm. At 1 GHz those become 15.9 ohm and 62.8 ohm. A decoupling capacitor therefore looks like a short only at high enough frequency, and its own lead inductance eventually dominates, which is the whole story of power-delivery design: every capacitor is useful only over a band of frequencies.

## Phase between voltage and current

Because $\mathbf{I} = \mathbf{V}/Z$, the current's phase is the voltage's phase minus the impedance angle. Inductive impedances (positive angle) make current lag voltage; capacitive impedances (negative angle) make it lead.

Example: 100 ohm in series with 1 pF at 1 GHz. $1/(\omega C) = 159.2$ ohm, so $Z = 100 - 159.2j$ with angle $-57.9^\circ$. The current leads the source voltage by $\arctan(1/(\omega RC)) = \arctan(1.592) = 57.9^\circ$. At low frequency the capacitor dominates and the lead approaches $90^\circ$; at high frequency the resistor dominates and it approaches zero. Phase is what turns into delay and into the phase margin of feedback loops, so frequency-domain engines always report it with the magnitude.

## Going deeper

Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 13; Nilsson and Riedel, "Electric Circuits", chapter 9.
