# Fin depopulation

In a FinFET library, each transistor's drive comes from a whole number of fins. Fin depopulation reduces the fin count per device, for example from three to two or from two to one, so that the active area, and the cell height, can shrink. It trades drive strength for density and is one of the main levers that let track heights fall.

## Drive scales with fin count

Each fin contributes roughly the same effective channel width, so drive current scales with the number of fins: going from three fins to two keeps about 66.7% of the drive, and from two to one, 50%. Input capacitance falls in the same proportion, so a gate driving another depopulated gate loses less speed than the drive ratio suggests; wire-dominated loads feel the full loss. Width is therefore quantized: designers cannot size a FinFET continuously, only in fin steps, and libraries offer drive variants by fins and fingers. Multi-finger devices recover drive at the cost of cell width.

## Cell height on the fin grid

Fins sit on a fixed fin pitch, and the cell height must contain the NMOS fins, the PMOS fins, the gap between the two regions, and the space at the rails. Written in fin pitches, a cell with $n$ fins per device and $g$ overhead pitches is $(2n + g)$ fin pitches tall. With a 30 nm fin pitch, two fins per device and four overhead pitches give $8 \times 30 = 240$ nm; depopulating to one fin gives $6 \times 30 = 180$ nm. Each fin removed per device takes two fin pitches out of the cell, which must then match a smaller metal track count.

## Drive density and its limits

Depopulation cuts drive faster than area when overhead is fixed. From two fins to one with four overhead pitches, drive halves and height falls from 8 to 6 pitches, so drive per unit area falls to $(1/6)/(2/8) = 0.667$. For logic dominated by small cells and short wires, the density gain outweighs this; for high-performance paths, designers mix taller, higher-drive cells. Depopulation also raises variability, since fewer fins average out less random variation. Gate-all-around nanosheets restore a continuous width knob, since sheet width can vary, which changes how this trade is made at the next nodes.

## Going deeper

Liebmann et al., DTCO papers on fin depopulation and track-height scaling in Proc. SPIE; Hu, "Modern Semiconductor Devices for Integrated Circuits", chapter on FinFETs.
