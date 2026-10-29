# Netlist handoff and annotation

The netlist is the contract between schematic capture and layout. It lists each net as a set of pins, written as reference designator and pin number, and binds each designator to a footprint. Annotation assigns the designators; forward annotation pushes schematic changes into a layout in progress; back annotation returns layout-side changes such as gate swaps. Mathematically the netlist is a hypergraph. A tool designer must keep it consistent across two databases that are edited concurrently, which is a data-synchronization problem.

## Nets, labels and annotation

A schematic becomes a netlist by resolving connectivity. Wires join pins on one sheet. **Net labels** join wires by name, with scope: a local label connects only within its sheet, a global label everywhere, and power symbols such as GND or +3V3 act as global labels, so one misspelled power name silently creates a separate floating net. Hierarchical sheets connect through sheet pins, which lets one sheet be instantiated several times without its local nets shorting together. **Annotation** assigns unique reference designators (R12, U3). The result is a hypergraph: components are vertices and each net is a hyperedge on the pins it joins, the same structure chip tools use. A net with a single pin is suspicious, since it connects nothing.

## Forward and back annotation

Design does not stop when layout starts. **Forward annotation** pushes schematic edits (a changed value, an added capacitor, a reconnected pin) into the layout as an incremental update, keeping placed and routed parts. **Back annotation** returns layout decisions: swapping two equivalent NAND inputs or two gates of a quad package shortens routes, and the schematic must be updated so it still describes the board. Libraries declare which pins and gates are swappable; anything else changes the circuit. To survive renumbering, tools match components by a hidden unique identifier rather than by designator. If they matched by designator, renaming R7 to R12 would rip up the placed R7 and drop a new unplaced R12, or worse, attach R7's placement to a different part.

## Consistency checks across the boundary

The schematic and the board are two databases that must agree. After an import, a ratsnest line to an unrelated part usually means a symbol and footprint number their pins differently. After routing, a **connectivity check** extracts which lands the copper joins and compares those nets with the netlist, reporting opens (a net split in copper) and shorts (two nets merged), the board analogue of LVS. At fabrication, an **IPC-D-356** netlist goes with the Gerbers so the bare-board tester can probe every net for continuity and isolation before parts are mounted. When schematic and layout come from different tools, the netlist is exported in a format the layout tool imports, and the fragile part is mapping footprint names between the two libraries.

## Going deeper

The KiCad documentation on schematic hierarchy, labels and the update-PCB-from-schematic flow. Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 1, for the hypergraph view of netlists.
