# Bills of materials

The bill of materials (BOM) lists what must be bought and placed to build one board: one line per manufacturer part number, with its quantity, reference designators and approved sources. It is generated from the schematic, consumed by purchasing and the assembler, and must agree with the layout's footprints and the centroid file. A BOM tool groups placements, rolls up quantities for a build, and manages variants that leave some parts unpopulated. Each of these is a small, exact computation, and mistakes in it are expensive.

## One line per part number

A BOM line is keyed by the **manufacturer part number** (MPN), not by value. Placements R1 to R12 that all use one 10 kohm 0402 resistor MPN collapse into one line with quantity 12 and the designator list. The same value in a different package is a different MPN and a different line: a 100 nF 0402 capacitor and a 100 nF 0603 capacitor need different reels and feeders. Example: placements R1, R2 (10 kohm 0402), R3 (10 kohm 0603), C1, C2, C3 (100 nF 0402) and U1 (an LDO) make four lines. The grouping is a hash on the MPN field. Its correctness depends on the library, so a library with two MPN entries for one physical part silently splits a line and doubles the setup cost.

## Rolling up an order

Purchasing multiplies each line's per-board quantity by the build quantity and adds attrition, the parts lost to machine pick-up failures, splicing and setup. Small passives on reels are commonly ordered with a few percent extra, or with a fixed minimum per line. Example: a line with 12 placements per board, a build of 250 boards and 5% attrition needs $12 \cdot 250 \cdot 1.05 = 3150$ parts. When the product is fractional the order rounds up, since half a part cannot be bought. Real BOM tools then round further to the supplier's reel or cut-tape increments and check the approved alternate sources. The arithmetic is trivial; the failure mode is a stale BOM that no longer matches the schematic revision being built.

## Variants and DNP

One bare board often serves several products. A **variant** marks some lines **DNP** (do not populate): the land patterns remain in the copper, but the assembler does not place the part. A low-cost variant might drop an optional sensor and its decoupling, so its placement count is the sum of the populated lines only. Example: lines with quantities 12, 11 (DNP), 4 and 5 (DNP) give $12 + 4 = 16$ placements per board. Variants must flow consistently to every assembly output: the BOM, the centroid file and the assembly drawing must agree on which designators are fitted. A variant that is DNP in the BOM but present in the centroid file makes the machine try to pick a part that is not loaded, which stops the line.

## Going deeper

The IPC-2581 specification's BOM and component sections, for how a neutral format carries BOM, variants and placement in one file.
