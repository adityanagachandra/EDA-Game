# PCB and package co-design agents

The die, its package and the board are designed by different teams in different tools, yet one decision couples them all: which signal sits on which ball. Ball maps, FPGA pin assignments and channel budgets span the boundary, so optimizing one side alone pushes cost to the other. Co-design agents are a natural fit because the work is search over swappable assignments with deterministic checks on both sides, provided the agent respects each side's rules, models and IP boundaries.

## What crosses the package boundary

Co-design optimizes the **ball map**: the assignment of signals to package balls, so that the package substrate routes from die bumps to balls and the board escapes from balls to its channels, both with few layers and short critical paths. Consistent data must cross the boundary: net names and the die-bump to ball to board-pad mapping, which, if wrong, is an open or short that no single-side check sees. FPGAs make pin assignment a board-side variable: many I/O can swap within a bank and I/O standard, so the board can choose pins that untangle routes. High-speed channels for signoff span the whole cascade, die pad, package, ball, board traces and vias, connector and far device, each represented by S-parameters joined for channel simulation.

## Agent actions and tiers

An agent swapping two FPGA pins in one bank must satisfy the vendor's pin rules (bank voltage, I/O standard, clock-capable pins), enforced by the FPGA tool, and the swap must flow back into the board netlist and schematic. Rerunning a ball-map optimizer with more weight on board layer count is **Tier 1**, a parameter change. Evaluating hundreds of candidate maps needs fast models in the inner loop: escape-routability counts, quasi-static channel models, simple PDN estimates; full-wave extraction is reserved for finalists. Encrypted package models from a supplier can be run by the simulator, but the agent cannot read or export them: the IP boundary is architecture, not policy, and the harness enforces it.

## Verifying a co-design change

A ball-map change is accepted on evidence from both sides. Connectivity checks across die, package and board confirm the mapping; SI re-simulation of affected channels confirms loss and eye margins; PI analysis confirms the power delivery. The last is easy to miss: moving power and ground balls toward the package edge can ease board routing while lengthening the PDN path to the die, raising its inductance and supply noise against the target impedance. The agent's report to the reviewer must contain the ball-map diff, the routing evidence, and SI and PI results against baselines on both sides, the same structure as an agentic ECO loop, where the agent proposes and verifies, and a human approves the change.

## Going deeper

The agent-first EDA notes in this curriculum (agentic signoff and ECO loops, sandboxing and the IP boundary). FPGA vendors' pin-planning and PCB design guides.
