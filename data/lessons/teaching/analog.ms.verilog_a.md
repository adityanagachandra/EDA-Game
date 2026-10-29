# Verilog-A behavioral models

Verilog-A is the standard language for describing analog behaviour as equations. Simulators compile its contributions into the same nodal equations they build from transistors, so a Verilog-A model can stand for anything from an ideal op-amp to a full compact transistor model such as BSIM-CMG. Writing models that are both correct and friendly to Newton iteration and timestep control is a core skill for anyone building mixed-signal flows.

## Contributions build the nodal equations

A Verilog-A module declares electrical ports and, in its analog block, contributes branch quantities with the operator <+. The statement I(p, n) <+ V(p, n)/r; adds a current through branch (p, n) equal to its voltage over r, exactly a resistor's stamp. Contributions to the same branch add, so a second contribution I(p, n) <+ c*ddt(V(p, n)); puts a capacitor in parallel. A voltage contribution, V(out) <+ expr;, defines a voltage-controlled branch and adds a current unknown, as a voltage source does in modified nodal analysis. ddt is the time derivative, which the simulator discretizes with its integration method in transient and replaces by $j\omega$ in AC; idt is the integral. The simulator also needs derivatives of every contribution with respect to node voltages for its Jacobian, which model compilers such as ADMS and OpenVAF generate automatically.

## Evaluating a behavioral model

A behavioral model is a function from node voltages and time to branch quantities, so evaluating it at an operating point is ordinary arithmetic. An amplifier modelled as V(out) <+ vsat*tanh(gain*V(in)/vsat); with vsat = 1 V and gain = 40 gives, at $V(in) = 10$ mV, an output of $\tanh(0.4) = 0.380$ V, not the 0.4 V an ideal linear model predicts. A branch I(p, n) <+ V(p, n)/r + c*ddt(V(p, n)); with r = 10 kilohms and c = 2 pF, at 0.5 V with a slope of 5 V per microsecond, carries $50 + 10 = 60$ microamperes. Checking such values by hand against simulator output is the simplest test that a model compiles to what its author intended.

## Writing models that converge

Newton iteration and timestep control assume smooth equations. A model with a hard discontinuity, such as an if statement that switches an output between rails, gives the solver an infinite derivative and causes tiny timesteps or convergence failure. Smooth functions (tanh, softplus limiters) or the transition() filter, which gives an output finite rise and fall times, fix this. Threshold-driven behaviour belongs in events: @(cross(V(clk) - vth, +1)) tells the solver to locate the crossing and run the event code there; @(timer(...)) schedules events in time. Expressions must stay finite at unphysical points Newton may visit, for example by guarding divisions. Finally, AC analysis linearizes at the DC operating point, so event-driven internal state does not appear in a small-signal result, a frequent source of confusion in behavioral PLL and converter models.

## Going deeper

Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (2004); McAndrew et al., "Best practices for compact modeling in Verilog-A" (IEEE Journal of the Electron Devices Society, 2015).
