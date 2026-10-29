# The moment-matched pi model

A cell's delay depends on its load, but a resistive wire shields the driver from the far capacitance, so the total capacitance overstates the load. Delay calculators replace each net, as seen from its driver, with a three-element pi model, near capacitance, resistance and far capacitance, chosen to match the first three moments of the driving-point admittance. This O'Brien-Savarino model is the input to effective-capacitance computation in practically every static timer.

## The uniform line

The driving-point admittance of an RC load expands as $Y(s) = y_1 s + y_2 s^2 + y_3 s^3 + \dots$. A pi model, $C_{near}$ at the driver, then $R_\pi$, then $C_{far}$, matches the first three terms when $C_{far} = y_2^2/y_3$, $C_{near} = y_1 - C_{far}$ and $R_\pi = -y_3^2/y_2^3$. For a uniform distributed line of total $R$ and $C$, $y_1 = C$, $y_2 = -RC^2/3$ and $y_3 = 2R^2C^3/15$, which gives $C_{far} = 5C/6$, $C_{near} = C/6$ and $R_\pi = 12R/25$. A line of 1250 ohm and 18 fF becomes 3 fF, 600 ohm and 15 fF. Note that $y_1$ is always the total capacitance, so the model keeps the driver's DC load right while moving most of it behind the resistance.

## Admittance moments of a chain

For a lumped chain, compute $Y(s)$ from the far end. A capacitor alone has $Y = sC$. Behind a resistor $R$, a load $Y$ becomes $Y/(1 + RY)$, expanded as a power series; then add the capacitance at the new node. For a driver pin with 9 fF, 0.1 kohm to node 1 with 15 fF, and 0.1 kohm to node 2 with 12 fF, this gives $y_1 = 36$ fF, $y_2 = -87.3$ kohm fF$^2$ and $y_3 = 291.87$ kohm$^2$ fF$^3$. Then $C_{far} = 87.3^2/291.87 = 26.11$ fF and $C_{near} = 9.89$ fF. With small resistances most capacitance stays visible to the driver: $C_{far}$ here sits behind only a small $R_\pi$. The recursion is linear in the tree size, and branches simply add admittances at their junction.

## The resistance and its use

The pi resistance is $R_\pi = -y_3^2/y_2^3$, positive because $y_2$ is negative. For a chain with 8 fF at the pin, 0.8 kohm to 20 fF and 0.5 kohm to 16 fF, the moments give $R_\pi = 0.926$ kohm, $C_{near} = 8.54$ fF and $C_{far} = 35.46$ fF. The delay calculator then finds an effective capacitance: the single capacitance that draws the same charge from the driver as the pi load during the driver's transition. A slow driver sees nearly all 44 fF; a fast one sees little more than $C_{near}$. The effective capacitance indexes the Liberty table for the cell delay, and the pi model, or the full reduced network, is then used again to compute the wire delay from driver to each sink.

## Going deeper

O'Brien and Savarino, "Modeling the driving-point characteristic of resistive interconnect for accurate delay estimation", ICCAD (1989); Qian, Pullela and Pillage, "Modeling the effective capacitance for the RC interconnect of CMOS gates", IEEE Transactions on CAD (1994).
