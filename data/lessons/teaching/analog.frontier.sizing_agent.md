# Sizing agents with SPICE as the golden judge

A sizing agent proposes device sizes for a fixed topology, runs the simulator through a harness, reads the results and proposes again. It can reason about specs the way a designer does, but it can also be confidently wrong. The architecture that makes it safe is simple: the agent proposes, SPICE with a fixed testbench judges, and only what SPICE certifies at every required corner is accepted.

## The harness loop for sizing

The harness automates the designer's own cycle. A request arrives ("meet 60 degrees of phase margin at all corners with less than 1 mW"). The agent edits parameter values in the netlist: widths, lengths, bias currents, the compensation capacitor. In the four-layer model this is a parameter-level change: the topology, the simulator and the PDK models are untouched. The harness runs the simulations, measures every spec and returns structured feedback: which specs failed, by how much, at which corner, and operating-point warnings such as devices leaving saturation. The agent uses that to choose its next move. Golden regressions, known circuits with known-good sizings, are re-run whenever the harness, models or testbenches change, so the loop itself stays trustworthy.

## Acceptance by the golden judge

Acceptance must be mechanical and owned by the harness. The rule is: every spec holds at every required corner, evaluated by SPICE with the frozen testbench. A sizing whose phase margin is 62 degrees at TT but 57 at SS fails a 60-degree spec, however good its nominal numbers look; the first acceptable iteration is the first whose worst corner passes. Non-converged corners count as failures, never as passes. The testbench, measurement scripts and specs are read-only to the agent: an agent that can edit the load capacitance or the measurement window can pass by changing the question. After acceptance, the result is shown to a designer as a reviewable diff of parameter values with the SPICE evidence attached.

## Failure modes and what the loop must still check

Agents optimize what the loop measures. If power is not in the objective, an agent meeting GBW at the slow corner will push bias currents to the maximum. If only PVT corners are checked, offset and yield under mismatch are unverified, because global corners do not capture local mismatch, so a Monte Carlo check belongs in acceptance for matching-sensitive blocks. If the agent may change topology (adding a cascode, say), the change is structural, beyond a parameter edit, and needs stronger verification and human review. And because netlists and PDK models are protected IP, the simulator and data stay in the company's sandbox, with the agent receiving only the results it needs. These rules turn an agent from a fluent suggester into a reliable sizing tool.

## Going deeper

The agent_first_eda lessons on the harness loop and on golden and differential testing.
