# Stage delay d = gh + p

The logical-effort stage delay $d = gh + p$ condenses the RC model into one line: a topology term, a load term and a self-loading term. It is the unit of reasoning for sizing, buffering and path design, and it is how a timer-free optimizer estimates delay across thousands of candidate choices in milliseconds.

## The three terms

Write the RC delay of a gate sized $k$ times a reference, driving load $C_{out}$, in units of $\tau = R_{inv} C_{inv}$, the delay of an ideal inverter driving an identical one without parasitics. The result is
$$d = g h + p, \qquad h = \frac{C_{out}}{C_{in}}.$$
Logical effort $g$ captures the topology, electrical effort $h$ (sometimes called fanout) captures how much bigger the load is than the gate's own input, and parasitic delay $p$ captures self-loading. The product $f = gh$ is the stage effort. Because $h$ is a ratio of capacitances, sizing the gate up increases $C_{in}$ and reduces $h$, which is the only knob the gate's own size gives you.

## A worked stage

A NAND2 with $\gamma = 2$ has $g = 4/3$ and $p = 2$. Sized with 4 fF of input capacitance and driving 24 fF, it has $h = 6$, so
$$d = \tfrac{4}{3} \cdot 6 + 2 = 10.$$
If $\tau = 5$ ps, that is 50 ps. A NOR2 in the same position has $g = 5/3$, so $d = 12$, 60 ps: the topology penalty is 2 units of $\tau$. If the NAND2 is instead sized to 8 fF of input, $h = 3$ and $d = 6$, but the gate driving it now sees twice the load. That tension between a gate and its driver is what path-level sizing resolves. Two stages in series simply add: a NAND2 at $h = 3$ followed by an inverter at $h = 4$ gives $(4 + 2) + (4 + 1) = 11$ units.

## What the linear model hides

The stage-delay model is linear in load and has no input-slew term; $g$, $p$ and $\tau$ are fitted to a process at typical slews, often from simulated or characterized delays. It treats rising and falling edges as equal because the gate was sized for equal drive, and it ignores wire resistance. Those approximations are acceptable when choosing sizes and topologies, because errors are similar across choices, but not for signoff. Real tools therefore use logical effort or similar linear models inside optimization loops (gate sizing by Lagrangian relaxation, buffer insertion, early synthesis estimates) and then verify with a timer reading Liberty tables, which are the measured nonlinear version of the same delay versus load relation.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapter 1; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.5.
