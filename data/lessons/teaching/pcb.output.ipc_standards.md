# IPC standards for design, fabrication and assembly

IPC standards are the shared vocabulary of the board industry: design rules, land patterns, current capacity, acceptance criteria for bare boards and assemblies, soldering requirements, and data formats. They are organized by stage and by product class, which sets how strict every criterion is. A board tool's rule decks, footprint generators and output writers implement pieces of these standards, so a tool designer needs to know which document governs which decision.

## Design standards

**IPC-2221** is the generic design standard: conductor spacing tables by voltage, conductor sizing, hole and land relationships. Its sectionals add detail: IPC-2222 for rigid organic boards, **IPC-2223** for flex and rigid-flex, **IPC-2226** for HDI and microvias. **IPC-7351** gives the land-pattern calculation and naming convention for SMT footprints. **IPC-2152** gives conductor current capacity from newer test data than the IPC-2221 chart, including the effect of planes and board thickness. **IPC-2141** is the design guide for controlled impedance. A rule deck derived from these documents is portable across fabricators and auditable, because each value traces to a table.

## Fabrication, assembly and data standards

Performance and acceptance follow design. **IPC-6012** is the qualification and performance specification for rigid boards, the spec a board is built and tested to. **IPC-A-600** shows acceptable and rejectable conditions of bare boards, such as plating voids and annular ring breakout. **J-STD-001** sets requirements for soldered assemblies, the process side, and **IPC-A-610** gives visual acceptance criteria for assemblies, the inspection side. Data exchange has its own documents: **IPC-2581** is an open XML format for the whole manufacturing job, and **IPC-D-356** is a netlist format for bare-board electrical test. Moisture-sensitivity classification of packages is in J-STD-020, a joint IPC and JEDEC standard.

## Product classes

IPC defines three product classes. **Class 1** is general electronics, where function matters but cosmetic imperfections are accepted. **Class 2** is dedicated-service electronics, where extended life and uninterrupted service are desired. **Class 3** is high-reliability electronics, such as life-critical medical or aerospace equipment, where downtime is unacceptable. The class selects criteria throughout the standards: minimum annular ring, plating thickness, allowed voids and fillet requirements all tighten from Class 2 to Class 3, which raises cost and lowers fabrication yield. A board's drawing states its class, and the tool's rule deck, footprint density level and output notes should all follow from it, so one setting propagates consistently instead of being re-entered per check.

## Going deeper

IPC-2221, IPC-6012, IPC-A-600 and IPC-A-610, available from IPC. The IPC standards tree, which maps documents by product life-cycle stage.
