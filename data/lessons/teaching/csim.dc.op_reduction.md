# The DC operating point problem

The DC operating point is the first thing almost every analysis computes: transient starts from it, AC and noise linearize around it. It is the steady state with all sources at their DC values and all time derivatives zero. That makes capacitors open circuits and inductors short circuits, and turns the circuit into a nonlinear resistive network. Knowing what the simulator actually builds for this problem, including the elements it keeps but ignores, avoids many surprises in .op output.

## Capacitors open, inductors short

At DC, $i_C = C\, dv/dt = 0$, so a capacitor carries no current and acts as an open circuit, and $v_L = L\, di/dt = 0$, so an inductor has no voltage across it and acts as a short. Reduce the circuit accordingly and solve. A source of 3 V through 1 k$\Omega$ into node a, an inductor from a to b, 2 k$\Omega$ from b to ground, a capacitor from a to ground and a capacitor from b to node c with 1 k$\Omega$ from c to ground gives $v_a = v_b = 3 \cdot 2/3 = 2$ V (the inductor joins a and b) and $v_c = 0$, because no DC current flows through the capacitor into c.

## Currents at the operating point

Branch currents follow from the reduced circuit. An inductor's DC current is whatever the resistive network forces through the short; a capacitor's DC current is zero. In the example the inductor current is the current that leaves node b through the 2 k$\Omega$ resistor, $2/2 = 1$ mA, because the capacitor at b carries none. If a resistor also hangs from node a to ground, the source current splits, but only the part continuing to b passes through the inductor. These currents are the initial conditions that transient analysis inherits: at $t = 0$ the inductor starts carrying exactly this current and each capacitor holds its node voltage.

## The DC system a simulator actually builds

SPICE does not delete the reactive elements. It keeps the same matrix structure for DC and transient, so capacitors stamp nothing at DC (their companion conductance is zero), and each inductor keeps a branch-current unknown whose row says $v_+ - v_- = 0$, a zero-volt source. The unknown count is therefore nodes plus voltage sources plus inductors. For the example circuit with nodes in, a, b and c, one voltage source and one inductor, the DC system has $4 + 1 + 1 = 6$ unknowns. A node reachable only through capacitors would have an empty row at DC, making the matrix singular, which is why simulators add a small gmin conductance or report a floating node.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; the ngspice manual, section on the .op analysis.
