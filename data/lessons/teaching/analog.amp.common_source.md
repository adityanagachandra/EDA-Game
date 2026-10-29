# The common-source amplifier

The common-source stage is the basic voltage amplifier: input at the gate, output at the drain, source at a fixed potential. Every op-amp gain stage is a variant of it. Its small-signal analysis introduces the three quantities a designer and a sizing tool track for every stage: gain, output resistance and the pole formed at the output node.

## Gain is transconductance times output resistance

Replace the transistor by its small-signal model. The input voltage produces a drain current $g_m v_{in}$, which flows into the parallel combination of the load and the transistor's own $r_o$. Writing KCL at the drain gives $v_{out} = -g_m (R_D \parallel r_o) v_{in}$. With $g_m = 2$ mS, $R_D = 20$ kilohms and $r_o = 100$ kilohms, $R_D \parallel r_o = 16.67$ kilohms and the gain is -33.3 V/V. The minus sign means inversion. Raising $R_D$ raises gain but also the DC drop across it, eating headroom, which is why integrated stages use current-source loads: their high small-signal resistance gives gain without a large DC drop, pushing the gain toward the intrinsic limit $g_m r_o$.

## Source degeneration

A resistor $R_S$ between source and ground feeds back the output current: the gate-source voltage becomes $v_{in} - i_d R_S$. Solving gives an effective transconductance $G_m = g_m/(1 + g_m R_S)$, and, ignoring $r_o$, a gain of $-G_m R_D = -g_m R_D/(1 + g_m R_S)$. For $g_m = 4$ mS, $R_S = 1$ kilohm and $R_D = 20$ kilohms, $G_m = 0.8$ mS and the gain is -16, down from -80 without degeneration. When $g_m R_S \gg 1$ the gain approaches $-R_D/R_S$, a ratio of resistors that matches well and varies little with process. That is the general bargain of negative feedback: less gain, more linearity and predictability.

## The output pole

The output node sees resistance $R_{out} = R_D \parallel r_o$ and the load capacitance $C_L$ plus device capacitances. Together they form a pole at $f_p = 1/(2\pi R_{out} C_L)$. With $R_{out} = 16.67$ kilohms and $C_L = 1$ pF the time constant is 16.67 ns (kilohms times picofarads is nanoseconds) and $f_p = 9.55$ MHz. The gain-bandwidth product $g_m R_{out} \cdot 1/(2\pi R_{out} C_L) = g_m/(2\pi C_L)$ does not depend on $R_{out}$: raising the load resistance trades bandwidth for gain one for one. The input also sees $C_{gd}$ multiplied by the gain (the Miller effect), which can create a second pole at the input. For a sizing tool, every high-impedance node is a pole to account for in stability.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 3; Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits" (5th ed., 2009), chapter 3.
