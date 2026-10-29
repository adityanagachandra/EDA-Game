# Transient thermal response

Temperature does not respond instantly to power. Every material stores heat, which acts as a thermal capacitance, so a power step produces an exponential rise with a thermal time constant. Short power bursts may therefore exceed the steady-state power budget without overheating, which is the basis of boost and turbo modes, and transient models predict how long they can last.

## Thermal capacitance

A body of volume $V$, density $\rho$ and specific heat $c_p$ stores heat with capacitance $C_{th} = \rho c_p V$, in J/K. Silicon has $\rho = 2330$ kg/m^3 and $c_p \approx 700$ J/(kg K), so about 1.63 J/(cm^3 K). A 20 mm by 20 mm die 0.75 mm thick has $V = 0.3$ cm^3 and $C_{th} = 0.489$ J/K, or 489 mJ/K. Copper heat spreaders store more, about 3.45 J/(cm^3 K), and heat sinks far more. The capacitances of the layers form a ladder with the resistances, the thermal analogue of an RC line, and each stage contributes its own time constant.

## The thermal time constant

A single lumped body with capacitance $C_{th}$ and resistance $R_{th}$ to ambient responds with time constant $\tau = R_{th} C_{th}$, exactly like an electrical RC. The die above, with 0.489 J/K and a 0.2 K/W local path, has $\tau = 97.8$ ms. Real packages have time constants from microseconds (transistor and wire) through milliseconds (die) and seconds (spreader) to minutes (heat sink and chassis). Each power step excites all of them. A transient thermal model is therefore a network of several RC stages, often fitted from a measured step response, and its structure is what compact transient models exchange between vendors.

## Response to a power step

For a lumped body, a power step $P$ applied at $t=0$ gives $\Delta T(t) = P R_{th}(1 - e^{-t/\tau})$. With $P = 50$ W, $R_{th} = 0.2$ K/W and $\tau = 100$ ms, the final rise is 10 K, and after 50 ms the rise is $10(1 - e^{-0.5}) = 3.93$ K. After one time constant the rise is 63% of its final value, after three about 95%. This is why a chip can briefly run above its sustained power: during a short burst the die has not yet reached steady state. Boost control integrates exactly this response, often with a multi-stage model, to decide how long the burst may continue before the junction limit is reached.

## Going deeper

Incropera and others, "Fundamentals of Heat and Mass Transfer", chapter 5 on transient conduction and the lumped capacitance method; JEDEC JESD51-14 on the transient dual interface method.
