# Symbols, footprints and land patterns

Every board part exists in two views. The schematic symbol is logical: pins with names, numbers and electrical types. The footprint is physical: copper lands, solder mask and paste openings, silkscreen, a courtyard and often a 3D body. The land pattern is the copper part of the footprint that the terminations solder to. The binding between the two views is the pin number. Library data is the most reused data in board design, and its errors are invisible to every automated check, which is why tool designers treat libraries as verified, versioned artifacts.

## Two views of one part

A **symbol** carries pin names (VDD, SDA), pin numbers (1, 2, ...) and pin electrical types (input, output, bidirectional, passive, power input, power output). The types exist for the schematic ERC. A multi-unit symbol splits one package into units, so a dual op-amp appears as U1A and U1B on different sheets but maps to one footprint. A **footprint** carries lands numbered to match, plus layers for solder mask openings, solder paste apertures, silkscreen legend, an assembly outline and a **courtyard**, the keep-out envelope used by placement checks. A **land pattern**, in IPC terms, is the set of lands. Netlisting writes each connection as designator plus pin number, for example U3.14, and layout attaches it to land 14 of U3's footprint. Names are only for humans.

## Pin-to-pad errors and footprint layers

The classic library error is a pin-number mismatch. A transistor symbol numbered 1 = B, 2 = E, 3 = C, paired with a package whose vendor numbers 1 = B, 2 = C, 3 = E, yields a board with collector and emitter swapped. DRC passes, ERC passes, connectivity matches the netlist, and the board fails on the bench. The footprint's other layers each feed a different process. The **paste** layer becomes stencil apertures. The **mask** layer opens solder resist around lands, with a small expansion (often a few mils) so registration error between the copper and mask imaging steps never covers a land. The **courtyard** includes body, lands and an assembly margin, so two parts whose bodies clear can still collide in the courtyard check. Each layer answers a question for one downstream consumer.

## Land patterns from package tolerances

Land patterns are computed, not drawn. IPC-7351 derives land size from the package's dimensions and their tolerances, the fabrication tolerance and the placement tolerance, so that a worst-case part still forms a good **toe, heel and side fillet**. It defines three density levels, most, nominal and least material, which trade board density against joint robustness. Special features need care: a QFN's exposed thermal pad gets a paste window pane of several small openings covering roughly half its area, because one full opening floats the part and opens perimeter joints. A 3D body model lets the tool check clearance against the enclosure. Because one footprint is reused on every board that uses the part, organizations verify it once and freeze it in a shared library rather than letting designers draw parts per project.

## Going deeper

IPC-7351, "Generic Requirements for Surface Mount Design and Land Pattern Standard". The KiCad library conventions (KLC), for a published, checkable set of symbol and footprint rules.
