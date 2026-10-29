# Compact-model plug-ins at the simulator seam

A circuit simulator's most important seam is the boundary between the solver and the device models. The solver's Newton loop, matrix assembly and time stepping are settled mathematics; device physics varies by technology and by customer. Placing a narrow, stable interface there lets new or customized models arrive as plug-ins, Tier 2 changes, without touching the core. Designing that seam well is a concrete exercise in the agent-first engine model.

## What crosses the seam

The solver asks one question per device evaluation: at these terminal voltages, temperature and instance parameters, what are the terminal currents, the terminal charges, and their derivatives with respect to the voltages? The answer fills the residual vector and the Jacobian. That is the whole interface.

The immutable core is matrix assembly, the Newton loop with its convergence checks, time-step control and the linear solver. The change surface is the device models: BSIM equations, a customer's self-heating variant, a process-specific stress term, a learned model.

Real models need internal nodes, for series resistance, a floating body or a thermal node, so the seam must let a model declare internal nodes and the equations for them. A seam that passes only terminal values would force self-heating into the core.

A working example exists in open source: OSDI, the binary interface through which ngspice loads compiled Verilog-A models. Verilog-A itself is the authoring language, and a compiler such as OpenVAF generates the evaluation code and derivatives.

## Rules that keep the seam fast and correct

Derivatives must be accurate and consistent with the currents. Finite-difference derivatives carry truncation and roundoff error that slows or breaks Newton convergence; compilers that generate analytic derivatives exist for this reason.

Calls should be coarse. Model evaluation is the hot loop of circuit simulation, and a virtual call per device with scattered data defeats vectorization. Evaluating many instances of one model type per call keeps overhead low and lets the model vectorize.

Plug-ins must be deterministic, with no hidden global state, so parallel evaluation gives the same result on every run regardless of thread scheduling.

A plug-in that passes its own unit tests has proved only its code. The harness must also run golden I-V and C-V sweeps and benchmark circuits through the real simulator, because what matters is behavior inside Newton across bias and temperature.

## Robustness, versioning and new model types

Newton iterations visit wild voltages before converging. A plug-in that returns NaN or overflows at $V_{GS} = -2$ V, a bias no normal test covers, will crash simulations unpredictably. The harness should sweep stress biases beyond the operating range and check for finite, bounded outputs and derivatives.

The extension ABI should be frozen and versioned. If it changes, existing plug-ins may break, so changing it is itself a core change with a migration cost.

The seam does not care how a model computes. A neural-network device model plugs in at the same place, but must meet the same contract: smooth currents and charges, consistent derivatives, charge conservation and sensible behavior outside its training data.

Agent-written plug-ins are iterated quickly, so they run in a sandbox with CPU, memory and time limits and no network, which contains failures and protects the IP boundary.

## Going deeper

The OSDI specification and OpenVAF documentation; the repository's docs/eda/agent-first-concepts.md, sections 1, 5 and 6; McAndrew and others, "Best practices for compact modeling in Verilog-A", IEEE Journal of the Electron Devices Society (2015).
