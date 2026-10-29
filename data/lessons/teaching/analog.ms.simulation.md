# Mixed-signal simulation

A mixed-signal chip combines analog blocks, solved as differential-algebraic equations, with digital logic, simulated as discrete events. A mixed-signal simulator couples an analog solver to an event-driven kernel and synchronizes them at the boundaries. Its partitioning, interface elements and synchronization determine both accuracy and speed, and they are what a tool designer must get right for top-level verification of converters, PLLs and SerDes.

## Two engines and their boundary

The analog engine is a SPICE-class solver: Newton iteration on the nodal equations at each timepoint, with timesteps chosen by local truncation error. The digital engine is an event-driven kernel: values change only at events, which are scheduled on a time-ordered queue, and nothing is computed between them. A mixed-signal simulator partitions the design into analog and digital instances. At every net that crosses the partition, Verilog-AMS inserts a connect module: a digital-to-analog element that drives a voltage with a specified rise time and output resistance, or an analog-to-digital element that turns a voltage into a logic value by thresholds, often with hysteresis. Values between the thresholds keep the previous state or become X after a timeout, depending on the module.

## Synchronizing the engines

The two engines advance time differently, so they must agree at every interaction. When a digital event drives an analog net, the analog solver must place a timepoint at the event and resolve the resulting edge, often with small steps. When an analog voltage crosses a connect-module threshold, the solver must locate the crossing time accurately, by shortening or backing up its step, and post a digital event at that time. Errors in crossing time become timing errors in the digital logic. Simulators use lockstep or more relaxed schemes in which the analog engine runs ahead and rolls back when a digital event arrives in its past. Synchronization cost grows with the number of boundary events, which is why partitioning a busy digital bus into the analog side is expensive.

## Choosing abstraction levels

Top-level simulation with every block at transistor level is far too slow: a PLL alone can need hours per microsecond. Practical verification mixes levels: the block under scrutiny at transistor level, other analog blocks as behavioral Verilog-A or real-number models, digital blocks in RTL. Each replacement must be validated: a behavioral model is only as good as the non-idealities it includes, and a model without comparator offset or DAC mismatch will pass regressions that silicon fails. Teams therefore keep model-versus-schematic checks that run the same testbench on both and compare key metrics within tolerance. For a tool designer, the mixed-signal simulator is the integration point where abstraction choices, partitioning and verification coverage meet.

## Going deeper

Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (2004); the Accellera Verilog-AMS Language Reference Manual.
