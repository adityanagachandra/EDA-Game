# Miller compensation and pole splitting

An uncompensated two-stage op-amp has two low-frequency poles close together and would oscillate in feedback. Miller compensation, a capacitor across the second stage, separates them: it pushes one pole down and the other up. It also creates a right-half-plane zero. These three locations, dominant pole, output pole and zero, are what sizing tools must place correctly.

## The Miller effect

A capacitor $C_c$ connected across an inverting stage of gain $-A_2$ draws current $C_c(1 + A_2)\,dv/dt$ from its input side, because the voltage across it changes $1 + A_2$ times as fast as the input. Seen from the input, it looks like a capacitance $C_c(1 + A_2)$ to ground. A 2 pF capacitor across a stage with gain 50 appears as 102 pF. Placed at the high-resistance output of the first stage, this large effective capacitance creates a very low pole, the dominant pole of the op-amp. The same effect is what makes the gate-drain capacitance of a common-source stage limit its input bandwidth; in compensation it is exploited deliberately.

## Pole splitting

Without $C_c$ the poles are $1/(R_1 C_1)$ and $1/(R_2 C_2)$ at the two stage outputs. With $C_c$, analysis of the two-node circuit gives approximately $p_1 \approx 1/(g_{m6}R_1 R_2 C_c)$ and $p_2 \approx g_{m6}C_c/(C_1 C_2 + C_c(C_1 + C_2))$. When $C_c$ and the load $C_2 = C_L$ dominate $C_1$, the second reduces to $p_2 \approx g_{m6}/C_L$. The dominant pole moves down and the output pole moves up: the poles split. For $g_{m6} = 2$ mS and $C_L = 2$ pF, $f_{p2} \approx 2000/(2\pi \cdot 2) = 159$ MHz. The physical picture: at high frequency $C_c$ shorts the second stage's gate to its drain, making it a diode-connected device whose low resistance $1/g_{m6}$ drives the load.

## The right-half-plane zero and its fix

$C_c$ also provides a feedforward path from the second stage's input to its output. At the frequency where the current through $C_c$ equals $g_{m6}v_{gs}$, the two cancel, giving a zero at $z = g_{m6}/C_c$ in the right half plane. It raises magnitude like a normal zero but lags phase like a pole, which erodes phase margin. With $g_{m6} = 2$ mS and $C_c = 1$ pF, $f_z = 318$ MHz. The standard fix is a nulling resistor $R_z$ in series with $C_c$, which moves the zero to $1/(C_c(1/g_{m6} - R_z))$. Choosing $R_z = 1/g_{m6}$, here 500 ohms, pushes the zero to infinity; a slightly larger $R_z$ moves it into the left half plane, where it can be placed on $p_2$ to cancel it.

## Going deeper

Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits" (5th ed., 2009), chapter 9; Allen and Holberg, "CMOS Analog Circuit Design" (3rd ed., 2011), chapter 6.
