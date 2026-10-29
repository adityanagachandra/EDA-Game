# Device recognition from layer Booleans

LVS compares the circuit drawn in the layout with the schematic. The layout has no devices, only shapes, so the first step is recognizing devices from layer combinations: a transistor is wherever poly crosses diffusion, typed by the well and implant layers around it, with W and L measured from the gate region. Recognition is a set of derived-layer rules in the LVS deck, and every later step depends on getting it exactly right.

## Gates from poly and diffusion

The gate layer is GATE = POLY AND DIFF. Each connected piece of it is one transistor channel, and the diffusion on its two sides becomes source and drain. A vertical poly line that crosses two diffusion rows inside their x-ranges forms two transistors, one per row, while a poly line passing beside a row forms none there. A poly line that only partly overlaps a diffusion edge makes a malformed device, which the deck flags. Counting gates is counting connected components of the GATE layer, and splitting diffusion by GATE is what separates the source and drain terminals.

## Channel width and length

For a gate region, L is measured along the direction of current flow, across the poly line, and W along the poly, across the diffusion. For a vertical poly $[34, 64] \times [-66, 984]$ crossing diffusion $[0, 240] \times [0, 920]$, the gate is $[34, 64] \times [0, 920]$, so L = 30 nm and W = 920 nm. Rotating the layout swaps which coordinate gives W, so extractors measure from the edges where gate meets diffusion (their total length is 2W for a rectangular gate) rather than from x or y. Non-rectangular gates, such as bent poly, need W and L formulas defined by the deck.

## Device types from wells and implants

The device type comes from the layers around the gate. In a p-substrate process, a gate whose channel lies inside NWELL with p-plus implant is a PMOS; outside the well with n-plus implant it is an NMOS. Thick-oxide, low-threshold and other variants are recognized by additional marker layers. If a well or implant is missing or misdrawn, the device is recognized as a different type or not at all, and LVS later reports a type mismatch even though the poly and diffusion are perfect. Recognition rules therefore encode the process's device catalogue.
