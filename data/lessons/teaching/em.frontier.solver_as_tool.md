# Field solvers and circuit kernels as agent tools

In an agent-first flow, field solvers and circuit kernels are not GUIs; they are callable tools with typed inputs and outputs that agents select from manifests and invoke. The mathematics stays in the immutable core. What changes is the interface: what an agent must pass in, what it must get back to decide its next step, and how calls are shaped so the kernel stays fast and reproducible.

## Typed inputs and outputs for a solver tool

A solver tool takes structured, validated inputs, never prose. For a capacitance solver: conductor geometry with explicit units, the dielectric stack (thicknesses and permittivities), terminal definitions or conductor voltages, frequency where relevant, and an accuracy target or mode. For an MNA kernel: a netlist, analysis type (operating point, AC sweep, transient) and tolerances. Outputs are keyed by stable names: a Maxwell matrix indexed by net, node voltages and branch currents by netlist name.

Example: an agent calls a 2D solver with widths in um and receives a 3 by 3 matrix in fF per um keyed by net. An input with a length in nm where um was expected should be rejected at the interface, because a factor of 1000 error otherwise flows silently into timing. Choosing between existing fast and accurate modes is a parameter setting, Tier 1 in the four-layer model; the kernels behind those modes are core.

## What the tool must report besides the answer

An agent decides where to recurse to accurate engines from what the tool reports about confidence. A capacitance tool should return, alongside the matrix, an error estimate or convergence indicator, the discretization actually used (panel or mesh count), and a validity flag stating whether the inputs fell inside the model's calibrated range. A transient kernel that fails to converge should return a structured failure: the time point, the step size reached and the offending nodes, not a partial waveform presented as success.

Example: a closed-form capacitance fit calibrated for $0.3 \le W/H \le 30$ receives $W/H = 60$. The right response is the value plus an out-of-range flag, or a refusal, so the agent calls the 3D solver. Silently clamping or extrapolating looks fine in the output and is exactly the failure that misleads a loop. Provenance fields make two results comparable.

## Batching and determinism

Interfaces should be coarse: an agent asking for capacitance of 20,000 nets should send them as a batch, so the engine amortizes setup, reuses factorizations and stays vectorized. A per-panel or per-edge call boundary would drown the kernel in overhead and prevent the compiler from optimizing the hot loop, the same reason extension seams are placed per batch rather than per element.

Determinism matters because agents compare runs to decide what changed. The same inputs and seed must give the same outputs, bit for bit or within a stated tolerance; randomized iterative methods must take a seed and record it, and parallel reductions must be ordered. Example: two identical calls that differ in the fourth digit make every differential report noisy. Finally, the tool runs on-site in a sandbox with resource limits and no network, because geometry and process data are the customer's IP.

## Going deeper

The Model Context Protocol specification for tool definitions; Nabors and White, "FastCap", 1991, for the batch structure of a capacitance solver.
