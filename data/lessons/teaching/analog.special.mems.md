# MEMS design

Micro-electro-mechanical systems integrate moving mechanical structures with electronics: accelerometers, gyroscopes, microphones, resonators and switches. Their design couples mechanics, electrostatics and often fluids, and their electronics are analog front ends. For an EDA designer, MEMS shows a multi-physics flow in which finite-element models are reduced to lumped spring-mass-capacitor models that circuit simulators, often through Verilog-A, can use.

## The spring-mass model

Most MEMS structures behave near their first mode as a mass on a spring with damping: $m\ddot{x} + b\dot{x} + kx = F$. The resonant frequency is $f_0 = \frac{1}{2\pi}\sqrt{k/m}$. A proof mass of 10 ng ($10^{-11}$ kg) on a spring of 4 N/m resonates at $\frac{1}{2\pi}\sqrt{4 \times 10^{11}} = 100.7$ kHz. Damping, set by air or gas around the structure, gives the quality factor; vacuum packaging raises Q to thousands for resonators and gyroscopes. An accelerometer measures the displacement $x = ma/k$ under acceleration $a$, so sensitivity and bandwidth trade against each other through $k/m = (2\pi f_0)^2$. These lumped parameters are extracted from finite-element simulation of the real geometry.

## Capacitive transduction

Most MEMS sense and actuate electrostatically. A parallel-plate capacitor has $C = \varepsilon_0 A/g$ for plate area $A$ and gap $g$, ignoring fringing. Plates 100 micrometres square with a 2 micrometre gap give $C = 8.854 \times 10^{-12} \cdot 10^{-8}/(2 \times 10^{-6}) = 44.3$ fF. A displacement changes the gap and therefore the capacitance, often by attofarads, which a low-noise readout circuit converts into a voltage; differential structures with two gaps cancel common-mode effects. Because the signal is so small, parasitic capacitance of the interconnect and the readout's kT/C noise and charge injection dominate the design, making MEMS readout a precision analog problem of exactly the kind treated in the converter lessons.

## Pull-in and multi-physics modelling

Electrostatic force grows as $1/g^2$ while the spring force grows linearly with displacement, so beyond a critical voltage the plate snaps down onto the electrode. For a parallel-plate actuator the pull-in voltage is $V_{PI} = \sqrt{8kg^3/(27\varepsilon_0 A)}$, reached at a deflection of one third of the gap. With $k = 5$ N/m, $g = 2\ \mu$m and $A = 10^{-8}\ \mathrm{m}^2$, $V_{PI} = \sqrt{3.2 \times 10^{-16}/2.39 \times 10^{-18}} = 11.6$ V. Pull-in is a limit for sensors and the operating principle of switches. Capturing it requires a nonlinear coupled model of mechanics and electrostatics. MEMS flows build such models by finite-element analysis, reduce them to a few states (model order reduction), and express them in Verilog-A so they can be simulated with their readout electronics.

## Going deeper

Senturia, "Microsystem Design" (2001).
