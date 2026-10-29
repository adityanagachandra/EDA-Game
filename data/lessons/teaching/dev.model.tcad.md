# TCAD device simulation

Technology computer-aided design (TCAD) simulates semiconductor fabrication and device behavior from physics. Process TCAD predicts structures and doping from implant, diffusion, oxidation and etch steps; device TCAD solves carrier transport in those structures to get I-V and C-V curves. TCAD is far too slow for circuits, but it is how a new device is explored, and how compact models are calibrated before silicon exists.

## Process and device simulation

Process TCAD models fabrication: ion implantation profiles, dopant diffusion and activation during anneals, oxidation, deposition and etch. Its output is a geometry with doping and material regions, the input to device simulation.

Device TCAD solves the drift-diffusion equations (or more detailed transport models) on a mesh of that structure, with physical models for mobility, recombination, tunneling and quantum corrections. Biasing the contacts and solving yields terminal currents and charges.

The mesh must be refined where quantities change over nanometres: near junctions, depletion edges and the gate-dielectric interface, where the inversion layer is about a nanometre thick. Most device simulators use a finite-volume (box) discretization with Scharfetter-Gummel fluxes along mesh edges, which conserves current exactly and remains stable for large potential steps between nodes. A 2D simulation of a planar device has thousands to tens of thousands of nodes; a 3D FinFET or nanosheet may have hundreds of thousands.

## Solving and calibrating

The discretized equations are a large sparse nonlinear system. Simulators solve them with Newton's method, either fully coupled or with Gummel-style decoupled iterations at low bias. Each Newton step needs a sparse linear solve; in 3D, fill-in makes direct factorization expensive, and iterative solvers with preconditioners are common.

Newton converges only from a good initial guess, so bias is ramped in steps, each converged solution seeding the next. A sweep of $I_D$ against $V_{GS}$ is a sequence of such continuation steps.

Small-signal capacitances come from an AC analysis that linearizes the system around a DC operating point, or from differentiating terminal charge with respect to voltage.

Calibration fits the physical model parameters, such as mobility coefficients, recombination lifetimes and doping-profile details, until simulated curves match measured silicon. A calibrated deck can then predict the effect of a change before it is fabricated, which is the purpose of TCAD in a process or design-technology loop.

## TCAD in the design flow

A single 3D TCAD bias point can take minutes; a full I-V family can take hours. A compact model evaluates in microseconds. This speed gap decides TCAD's role: it never runs inside circuit simulation. Instead it feeds compact-model extraction for new technologies, design-technology co-optimization (DTCO), and variability studies.

In DTCO, device options (fin height, sheet width, spacer thickness, contact scheme) are simulated in TCAD, turned into compact models, used to characterize standard cells and evaluated on circuit metrics such as ring-oscillator delay and power. Statistical TCAD samples random dopant positions, edge roughness or metal grains to produce variability distributions that statistical compact models must reproduce.

Commercial TCAD suites dominate industrial use. In open source, DEVSIM is a programmable finite-volume device simulator with a scripting interface, useful for learning and for research flows. For an EDA tool designer, TCAD is a slow, accurate engine at the device end of a loop whose fast engines are compact models and surrogates.

## Going deeper

Vasileska, Goodnick and Klimeck, "Computational Electronics: Semiclassical and Quantum Device Modeling and Simulation" (2010); the DEVSIM documentation; Selberherr, "Analysis and Simulation of Semiconductor Devices" (1984).
