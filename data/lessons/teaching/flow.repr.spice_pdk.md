# SPICE netlists and PDK decks

A process design kit (PDK) is the foundry's package that lets a design team use a process: device models, design rules, verification decks, extraction data, and often basic cells. SPICE netlists are the transistor-level circuit descriptions those models simulate. Together they are the lowest-level representations in the flow and the ground truth that every higher abstraction, from Liberty tables to LEF rules, was derived from.

## SPICE netlists

A SPICE netlist lists devices one per line with their nodes and parameters: `M1 out in vdd vdd pmos W=0.4u L=0.05u`, resistors, capacitors, sources, and subcircuits that group them. Analysis commands request DC operating points, transient waveforms or AC responses. The simulator writes Kirchhoff's laws over the nodes with nonlinear device equations from the model files and solves them numerically. SPICE netlists appear in three places in a digital flow: library characterization, analog and custom design, and LVS, where a netlist extracted from layout is compared with the schematic netlist.

## What a PDK contains

A PDK typically includes compact device models (for example BSIM for MOSFETs) with corner files for process variation; the design-rule manual and the DRC deck that encodes it for a checker such as Calibre or KLayout; the LVS deck that defines how devices are recognized from layers; extraction technology files with the layer stack's thicknesses and dielectric constants; the layer map that gives meaning to GDSII layer numbers; and parameterized cells (PCells) for transistors and passives. Open PDKs such as SkyWater SKY130 and GlobalFoundries GF180MCU let open-source tools such as Magic, KLayout, Netgen and ngspice run real flows.

## The PDK as the root of every abstraction

Every higher representation was computed from the PDK. Liberty tables come from SPICE simulations of each cell with the PDK's models. LEF rules and pitches come from the design-rule manual. Extraction accuracy comes from the technology files. So a PDK update, such as revised models or a new rule, can invalidate a library's Liberty and LEF, the extraction calibration, and every signoff result that depended on them. Flows track PDK versions for that reason, and a mismatch between the PDK used for characterization and for signoff is a real and costly failure mode.

## Going deeper

The SkyWater SKY130 PDK documentation; Nagel, "SPICE2: A Computer Program to Simulate Semiconductor Circuits" (UC Berkeley, 1975).
