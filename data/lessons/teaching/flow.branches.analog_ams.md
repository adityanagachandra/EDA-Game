# Analog and mixed-signal design flows

Analog and mixed-signal (AMS) circuits, such as amplifiers, data converters, PLLs, regulators and RF front ends, do not fit the RTL-to-GDSII flow. Their behaviour is continuous, so there is no gate abstraction to synthesize to and no static timing to sign off with. Their flow is schematic-driven, simulation-heavy and largely manual, and automating it is one of EDA's open problems.

## A schematic-driven flow

An analog designer draws a transistor-level schematic, sizes each device (width, length, number of fingers) and simulates it with a SPICE-class simulator across corners: DC operating point, AC gain and phase margin, transient settling, noise, distortion. The layout is then drawn largely by hand or with assisted placement, using the PDK's parameterized cells, and checked with DRC and LVS. Post-layout extraction adds parasitics, and the circuit is re-simulated, because a few femtofarads on a sensitive node can change bandwidth or stability. Cadence Virtuoso with Spectre is the dominant commercial environment; Xschem, Magic, KLayout and ngspice form an open-source one.

## Why analog resists automation

Digital design has a clean abstraction: a gate is a Boolean function plus a delay table, and correctness reduces to logic and timing. Analog performance depends directly on device physics: matching between transistors, parasitic capacitance, layout symmetry, substrate coupling, process variation. The specification is a list of continuous metrics (gain, bandwidth, noise, power, offset) with tight trade-offs, and simulation is the only reliable evaluator. Layout constraints such as common-centroid placement for matching and symmetric routing for differential pairs are hard to state to a general placer. So analog blocks take a small share of a chip's area and a large share of its design time.

## Mixed-signal integration

A mixed-signal chip combines analog blocks with large digital logic. The analog blocks are usually hard macros in the digital floorplan, with their own power domains, guard rings and keep-out zones to limit noise coupling. Verification mixes abstractions: analog blocks are modeled behaviourally (Verilog-AMS or real-number models) so that a chip-level simulation can run in reasonable time, with transistor-level simulation reserved for the blocks themselves. Interfaces such as level shifters, the analog-digital boundary of a converter, and supply sequencing are frequent sources of first-silicon bugs because they fall between the two flows.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017); Rutenbar, Gielen and Roychowdhury, "Hierarchical Modeling, Optimization, and Synthesis for System-Level Analog and RF Designs", Proceedings of the IEEE (2007).
