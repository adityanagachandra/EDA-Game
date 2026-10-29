# Inertial and transport delay

A gate delay in a simulator is not just a number added to a timestamp. The simulator must decide what happens when an input changes again before the previous output change has matured. The two classic answers are **transport** delay, which shifts every change, and **inertial** delay, which swallows pulses shorter than the delay. Verilog gate and continuous-assignment delays are inertial. The choice decides which glitches a gate-level simulation shows and how many events it processes.

## Transport and inertial delay

With **transport** delay $d$, every input change at time $t$ produces an output change at $t + d$, whatever happens in between: the output waveform is the input waveform shifted right. It models an ideal wire or delay line. With **inertial** delay $d$, a new input change cancels any output change still pending, and schedules the new value for $t + d$. A pulse narrower than $d$ therefore never appears at the output: its second edge cancels the first before the first matures. This models the fact that a real gate needs some energy, sustained for some time, to switch its output. In Verilog, `assign #5 y = a;` and gate primitives use inertial delay; transport-like behavior is written with nonblocking intra-assignment delays such as `y <= #5 a;`.

## Computing the output of a pulse train

Apply the inertial rule mechanically. A buffer has inertial delay 3 ns; its input starts at 0 and toggles at 2, 8, 10 and 18 ns. At 2 ns the input rises; schedule output 1 at 5 ns. Nothing interrupts it, so the output rises at 5 ns. At 8 ns the input falls; schedule output 0 at 11 ns. At 10 ns the input rises again before 11 ns, so the pending fall is cancelled and output 1 is scheduled for 13 ns, but the output is already 1, so nothing changes. At 18 ns the input falls; the output falls at 21 ns. Result: two output transitions, at 5 and 21 ns, because the 2 ns low pulse was swallowed. Under transport delay all four edges would appear, at 5, 11, 13 and 21 ns.

## Consequences for gate-level simulation

Inertial filtering means gate-level simulation hides glitches shorter than gate delays, which is mostly realistic but can mask a hazard on an asynchronous input such as a reset or clock. SDF annotation refines this with pulse limits: a pulse can be rejected (filtered), passed as X (an error pulse), or passed normally, depending on its width relative to reject and error limits, and simulators expose options to set those percentages. At RTL, with zero delays, there are neither delays nor pulse filtering: delta cycles order the changes but every zero-time glitch is still processed. This is one reason gate-level simulation with timing is slow and is reserved for checks like reset sequencing and timing-annotated sign-off of asynchronous paths.

## Going deeper

IEEE Std 1800-2017, clause 28 (gate-level modeling and delays) and clause 30 (specify blocks and pulse filtering); IEEE Std 1497 (SDF).
