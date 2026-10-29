# Capacitors and inductors as circuit elements

Capacitors and inductors are the elements that give circuits memory. A capacitor relates current to the rate of change of voltage, an inductor relates voltage to the rate of change of current. Every interconnect delay, every supply droop and every ground-bounce glitch comes from these two relations, and every transient simulator spends its time integrating them.

## Capacitor current from voltage slew

A capacitor stores charge $q = Cv$, so its current is $i = C\,dv/dt$. A constant voltage draws no current; a fast voltage ramp draws a large one. Integrating gives $v(t) = v(0) + \frac{1}{C}\int_0^t i\,dt$, so the voltage across a capacitor cannot jump unless infinite current flows.

Example: a driver slews a 20 pF load by 1.8 V in 0.6 ns. The current during the ramp is $20 \cdot 1.8/0.6 = 60$ mA (picofarads times volts per nanosecond is milliamperes). Multiply by the number of outputs switching together and the supply network must deliver amperes in nanoseconds, which is where decoupling capacitance and ground bounce come from. On chip, femtofarads and picoseconds give the same scaling.

## Inductor voltage from current slew

An inductor links magnetic flux $\lambda = Li$, and by Faraday's law the voltage across it is $v = L\,di/dt$. A steady current develops no voltage; a fast change in current does. Integrating, $i(t) = i(0) + \frac{1}{L}\int_0^t v\,dt$, so inductor current cannot jump.

Example: a current through a 3 nH bond wire changes by 40 mA in 0.5 ns. The induced voltage is $3 \cdot 40/0.5 = 240$ mV (nanohenries times milliamperes per nanosecond is millivolts). That voltage appears between the die's ground and the board's ground, which is ground bounce. Faster edges make it worse, because $di/dt$ grows as rise times shrink, so the packaging inductance becomes a timing and noise problem rather than a footnote.

## Stored energy

Charging a capacitor to $V$ stores energy $E = \frac12 CV^2$ in its electric field; a current $I$ in an inductor stores $E = \frac12 LI^2$ in its magnetic field. Both follow from integrating $p = vi$ with the element law.

Example: 10 pF at 1.2 V stores $0.5 \cdot 10 \cdot 1.44 = 7.2$ pJ. A 5 nH inductor carrying 100 mA stores $0.5 \cdot 5 \times 10^{-9} \cdot 0.01 = 25$ pJ. These energies explain two facts a power or reliability tool relies on. Charging a capacitor from a fixed supply through any resistance draws $CV^2$ from the supply: half is stored, half is dissipated in the resistance, which is the origin of the $CV^2 f$ dynamic-power term. And the energy stored in an inductor must go somewhere when its current is interrupted, which is the source of inductive voltage spikes.

## Going deeper

Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 9; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 5 on dynamic power.
