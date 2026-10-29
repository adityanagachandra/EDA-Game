# Odd and even modes of coupled lines

Two traces routed side by side over a plane are coupled: each has mutual capacitance and mutual inductance to the other. Any pair of signals on them can be split into two modes, the odd mode (opposite voltages) and the even mode (equal voltages), each with its own impedance and velocity. Differential signaling lives in the odd mode, common-mode noise in the even mode. Field solvers in board tools report both, and every differential-pair rule is a statement about them.

## Two modes of a coupled pair

Drive line 1 with $V_1$ and line 2 with $V_2$. Write them as a differential part $V_1 - V_2$ and a common part $(V_1 + V_2)/2$; the **odd mode** carries the first (the lines at $+V$ and $-V$) and the **even mode** the second (both at $+V$). For a symmetric pair the modes are independent waves. In the odd mode the neighbor is at the opposite voltage, so the mutual capacitance $C_m$ charges through twice the voltage and adds to each line's capacitance, while the opposite currents make the mutual inductance $L_m$ subtract: $L_{odd} = L - L_m$, $C_{odd} = C + C_m$. In the even mode the reverse holds. Hence $Z_{odd} = \sqrt{(L - L_m)/(C + C_m)}$ falls below the uncoupled $Z_0$ and $Z_{even}$ rises above it as coupling tightens.

## Differential and common-mode impedance

A differential receiver sees the two lines in series, line to line, so the **differential impedance** is $Z_{diff} = 2 Z_{odd}$. A common-mode signal drives both lines in parallel against the reference, so the **common-mode impedance** is $Z_{cm} = Z_{even}/2$. Example: $Z_{odd} = 48$ ohms and $Z_{even} = 58$ ohms give $Z_{diff} = 96$ ohms and $Z_{cm} = 29$ ohms. A single 100 ohm resistor across the pair terminates only the odd mode; common-mode current flows through both lines in the same direction and never crosses that resistor. To terminate both modes, a T network puts $Z_{odd}$ from each line to a center node and $(Z_{even} - Z_{odd})/2$ from the node to ground: in the odd mode the node is a virtual ground, and in the even mode each line sees $Z_{even}$.

## Velocities, crosstalk and reference planes

In **microstrip** the odd mode puts more of its field in the air between and above the traces, so it sees a lower effective dielectric constant and travels faster than the even mode. That velocity difference is what produces far-end crosstalk: the capacitive and inductive coupling terms arrive at different times instead of cancelling. In **stripline** in a uniform dielectric both modes travel at one speed, and far-end crosstalk nearly vanishes. Board pairs are edge-coupled over a plane, and most of each trace's field still ends on the plane, which sets the impedances and returns common-mode current. A pair that crosses a split in its reference plane loses that symmetry, and part of its differential signal converts into common mode, which radiates and closes the eye.

## Going deeper

Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs" (2009), chapter 3, on coupled transmission lines. Pozar, "Microwave Engineering" (4th ed., 2011), section on coupled lines.
