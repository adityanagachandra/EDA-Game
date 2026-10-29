# Verilog-AMS and VHDL-AMS

Verilog-AMS and VHDL-AMS extend digital HDLs with continuous-time analog behavior: voltages and currents governed by differential-algebraic equations, solved by a SPICE-like engine alongside the event-driven digital kernel. They are how analog blocks are modeled behaviorally, how compact device models are written (in the Verilog-A subset) and how mixed-signal chips are verified. Tool designers meet them in circuit simulators and mixed-signal flows.

## Disciplines and conservative quantities

Verilog-AMS describes analog signals through disciplines. The electrical discipline pairs a potential nature (voltage) with a flow nature (current). A discipline with both is conservative: when modules connect at a node, the simulator enforces Kirchhoff's laws, so currents into the node sum to zero and connected terminals share one potential. A signal-flow discipline has only a potential and models an ideal block diagram with no loading.

VHDL-AMS (IEEE 1076.1) expresses the same idea with terminals, which are analog nodes, and quantities defined across and through them: an across quantity is a voltage between terminals, a through quantity a current. Either way, the model's job is to state the branch equations; the simulator assembles the network.

## Contribution statements and the analog solver

An analog block states branch relations with the contribution operator. I(p, n) <+ V(p, n) / R; contributes a resistor's current, and I(p, n) <+ C * ddt(V(p, n)); contributes a capacitor's, where ddt is the time derivative (idt is the time integral). Contributions to the same branch add, so these two statements together model a resistor and capacitor in parallel.

The simulator combines all contributions with Kirchhoff's laws into a system of differential-algebraic equations $F(x, \dot{x}, t) = 0$ and integrates it in time, solving the nonlinear equations at each time point with Newton iteration, exactly as SPICE does for transistors. VHDL-AMS writes the equations directly as simultaneous statements, such as v == r * i;, which are declarative constraints rather than assignments.

## Mixed-signal boundaries

In a mixed-signal simulation, a digital output that drives an analog input passes through a connect module, inserted automatically according to discipline rules. It converts logic values into a voltage waveform with modeled rise and fall times; in the other direction, a comparator with a threshold turns an analog voltage into digital events, ideally with hysteresis so that noise near the threshold does not produce a storm of spurious toggles.

Every boundary crossing forces the analog solver to stop at that time point and synchronize with the digital kernel, so designs with frequent interaction simulate slowly. Compact device models written in Verilog-A are compiled for circuit simulators, for example by OpenVAF for ngspice, so that the model's currents and their derivatives are evaluated quickly inside every Newton iteration.

## Going deeper

Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (2004); Accellera, the Verilog-AMS Language Reference Manual; IEEE Std 1076.1 for VHDL-AMS.
