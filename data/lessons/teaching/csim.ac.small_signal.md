# AC small-signal analysis

AC small-signal analysis combines the two previous ideas: linearize every device at the operating point, then solve the complex MNA system over frequency. Its output is the transfer function of the linearized circuit: gain, bandwidth, phase. A common-source stage with a load capacitor is the canonical example, and its gain and pole frequency are what a designer checks first in any AC plot.

## Low-frequency gain

For a common-source stage with transconductance $g_m$, drain resistor $R_D$ and transistor output resistance $r_o$, the output node sees $R_{out} = R_D \parallel r_o$ and the low-frequency gain is $A_0 = -g_m R_{out}$. With $g_m = 2$ mS, $R_D = 20$ k$\Omega$ and $r_o = 80$ k$\Omega$, $R_{out} = 16$ k$\Omega$ and $|A_0| = 32$, which is $20 \log_{10} 32 = 30.1$ dB. Millisiemens times kilohms is dimensionless. The minus sign is a phase of 180 degrees at low frequency. A simulator reports the same number at the low end of the sweep, provided the operating point placed the device in saturation.

## Gain at a frequency

A load capacitance $C_L$ at the output adds a pole: $A(j\omega) = -g_m R_{out}/(1 + j\omega R_{out} C_L)$. Its magnitude is
$$|A| = \frac{|A_0|}{\sqrt{1 + (f/f_p)^2}}, \qquad f_p = \frac{1}{2\pi R_{out} C_L}.$$
With $R_{out} = 16$ k$\Omega$ and $C_L = 1$ pF, $R_{out} C_L = 16$ ns and $f_p = 9.95$ MHz. At 20 MHz, $|A| = 32/\sqrt{1 + 4.04} = 14.3$. Far above the pole the gain falls by a factor of 10 per decade, and the product $|A| f$ approaches the unity-gain frequency $g_m/(2\pi C_L)$. The simulator computes each point by a complex solve, not by this formula.

## The pole frequency

The pole is where $|A|$ is 3 dB below $|A_0|$ and the added phase is -45 degrees: $f_p = 1/(2\pi R_{out} C_L)$. Kilohms times picofarads is nanoseconds, so with $R_{out} = 33.3$ k$\Omega$ and $C_L = 2$ pF the time constant is 66.7 ns and $f_p = 2.39$ MHz. AC analysis also reveals what this one-pole model omits: the gate-drain capacitance couples input to output, creating a right-half-plane zero and, with a source resistance, a Miller-multiplied input pole. A pole-zero analysis finds these roots directly from the generalized eigenproblem of $G$ and $C$.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., McGraw-Hill, 2017), chapter 6; Najm, "Circuit Simulation" (Wiley, 2010).
