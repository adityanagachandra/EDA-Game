# PCB and package design

A chip is useless until it is packaged and mounted on a board. Package design connects the die to the outside world through bumps, substrate traces and balls; printed circuit board (PCB) design connects packages, passives and connectors through layers of copper. Both are EDA branches with their own tools and their own physics: signal integrity, power integrity and thermal behaviour at millimetre scale rather than nanometre scale.

## The board flow

A PCB flow starts from a schematic of components and their connections, which is compiled into a netlist with physical footprints. Layout places components on the board and routes traces on several copper layers, with vias between them. Rules come from the board fabricator's capabilities (trace width, spacing, drill sizes) and from electrical needs: controlled-impedance traces, length matching for buses, differential pairs. The outputs are fabrication data, traditionally Gerber files plus drill files, or richer formats such as ODB++ and IPC-2581, and an assembly bill of materials. Cadence Allegro, Siemens Xpedition and PADS, Altium Designer and Zuken CR-8000 are commercial tools; KiCad is open source.

## Packages and advanced packaging

A package fans out the die's dense bumps to a pitch that a board can use, distributes power and removes heat. Wire-bond packages connect pads by thin wires; flip-chip packages solder the die face-down onto a substrate. Advanced packaging puts several dies in one package: side by side on a silicon interposer or bridge (2.5D), or stacked with through-silicon vias (3D). This turns package design into a chip-like problem, with routing of thousands of die-to-die connections, and it ties the chip floorplan to the package: bump positions constrain where I/O and power enter the die.

## Signal, power and thermal integrity

At board and package scale, wires are transmission lines. A fast signal reflects at impedance discontinuities, couples into neighbours and loses amplitude over distance, so high-speed links are designed with channel simulation and eye diagrams. Power delivery from regulator to die must keep impedance low across frequency, using decoupling capacitors placed on the board, in the package and on the die. Heat must flow out of the package. These analyses are coupled across die, package and board, which is why co-design across all three is a major direction for EDA tools.

## Going deeper

Bogatin, "Signal and Power Integrity, Simplified" (3rd ed., 2018); Tummala (ed.), "Fundamentals of Microsystems Packaging" (2001).
