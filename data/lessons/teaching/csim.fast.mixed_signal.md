# Mixed-signal simulation and co-simulation

Mixed-signal simulation runs an analog solver and an event-driven digital simulator together on one design. The two engines think in different currencies, continuous voltages and discrete events, and in different notions of time, adaptive steps and event queues. Making them agree at every boundary and at every instant is the whole problem, and its solutions recur in every co-simulation, from RTL with SPICE to chip with package.

## Crossing the analog-digital boundary

Where an analog net drives a digital input, an A-to-D connect module decides when the voltage counts as 0, 1 or unknown, usually with two thresholds (hysteresis) so that noise on a slow edge does not produce a burst of spurious events. Where a digital output drives an analog net, a D-to-A module presents a voltage or current source with finite rise time and output resistance, because an ideal step would be both unphysical and a numerical discontinuity. Connect modules are inserted automatically where nets of different disciplines meet. Their thresholds, rise times and impedances change event timing, so they are part of the model, not plumbing.

## Synchronizing two engines

A digital event that changes an analog source must land on an analog time point, so it becomes a breakpoint for the analog solver. In the other direction, an analog waveform may cross a threshold inside a step; the crossing is located by interpolation or bisection, the analog side backs up or cuts its step, and the digital engine is notified at the crossing time. In lockstep schemes neither engine advances past the next point where it might need input from the other. Waveform relaxation instead lets each subsystem simulate a whole window from the other's previous waveform and iterates, which converges quickly when coupling is weak or one-directional.

## Co-simulation architectures

In RTL-with-SPICE co-simulation, the analog engine usually dominates runtime, since it solves Newton systems at small steps while the digital side processes events cheaply. Real-number models move analog behavior into the digital simulator as discrete-time real values passed as events, trading fidelity for orders of magnitude in speed in chip-level regressions. ngspice's XSPICE extension provides event-driven code models, including digital ones, inside the analog simulator. When a co-simulation disagrees with an all-SPICE reference by an event, suspect the interface thresholds first, because they decide when the analog waveform becomes a logic transition.

## Going deeper

Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (Kluwer, 2004); the ngspice manual, chapters on XSPICE code models.
