# Calibration as a Tier 1 change surface

An extractor's results depend on calibration data: the process stack, resistivity tables, fit coefficients and pattern tables. In the four-layer engine model that data is the parameter layer, so most customer and foundry requests about extraction accuracy should land in Tier 1, configuration with no recompile. Recognizing which requests do not, and proving a recalibration, is the core skill of an agent maintaining an extractor.

## What calibration data an extractor carries

An extractor's technology file and tables include layer thicknesses and heights, dielectric permittivities per layer, conductor resistivities with width and temperature dependence, via resistances, coefficients of fringe and coupling fits, and pattern-matched capacitance tables generated offline by field solvers for that stack. All of it is data the engine reads.

What is not calibration: the field-solver kernel (core), the algorithm that decomposes a net into patterns (strategy or core), and the design's netlist and layout (input). Example: when a foundry revises a dielectric from $\varepsilon_r = 2.9$ to 2.75, the response is to update the stack, rerun the calibration solves and ship new tables, with no code change. Keeping this boundary clean is what lets the common change stay at Tier 1, and it is the extractor's equivalent of a timer keeping its derating choices out of the propagation core.

## Placing a request in the right tier

Tier 1 is a value or table change; Tier 2 is new behavior at a defined extension point; Tier 3 is a core edit with human approval. Examples for an extractor:

A customer supplies a measured copper resistivity table versus line width: Tier 1, because the engine already reads width-dependent resistivity. A customer wants a correction for wires near a particular dummy-fill pattern, which needs code to detect that neighbourhood: Tier 2, a plug-in at the pattern-correction seam. A request to change the boundary-element quadrature for near-singular panel integrals: Tier 3, because it is the numerical kernel and can shift every result.

When a request arrives as a symptom, such as capacitance 3% low against silicon on metal 4, the agent should look first at the smallest safe tier: metal 4 stack thickness, permittivity and table coefficients. A systematic, layer-specific bias is almost always calibration.

## Proving a recalibration

A recalibration is an engine change and is verified like one. Field-solver goldens on reference structures must match the regenerated tables; a differential report on a reference design must show that only the intended layers moved, by the intended amounts; and the calibration version must be recorded with every result so that any number can be traced and rolled back.

Example: after recalibrating metal 2, the differential report shows metal 6 capacitances also shifted by 1%. That is out of scope and blocks promotion until explained, often by a shared stack parameter edited by mistake. A subtler failure: new tables match the field solver on every golden but correlate worse with silicon. Goldens prove table-to-solver consistency, not solver-to-silicon truth, so the stack data fed to the solver is the suspect, and correlation against measured test structures is the check.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8; the OpenRCX documentation in OpenROAD on extraction rule files and calibration.
