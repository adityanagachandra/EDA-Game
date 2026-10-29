# Buried power rails

In a conventional cell, the power and ground rails run along the top and bottom edges on the lowest metal layer and take up track space. A buried power rail (BPR) moves them below the transistors, into the substrate or shallow trench isolation, and connects to the devices from beneath. The freed tracks can shrink the cell or be returned to signal routing, and the rail can be made taller, which lowers its resistance.

## Rails that consume tracks

Front-side rails are wide wires on the M0 or M1 layer at the cell boundary, shared half-and-half with the abutting rows. Together they typically occupy one to two tracks' worth of a cell's height. A 6-track cell with 1.5 tracks of rails has 4.5 tracks for signals. Burying the rails removes them from the routing layer: at a 24 nm pitch, a cell keeping its 4.5 signal tracks drops from $6 \times 24 = 144$ nm to $4.5 \times 24 = 108$ nm, a 25% reduction, if nothing else limits height.

## Lower rail resistance

A buried rail does not compete for routing tracks, so it can be narrow but tall. A front-side M0 rail 30 nm wide and 40 nm thick has cross-section 1200 nm^2; a buried rail 20 nm wide and 120 nm tall has 2400 nm^2, half the resistance per unit length at the same resistivity. Lower rail resistance reduces IR drop along the row, which is an increasing problem as rails narrow. The rail metal must survive the remaining front-end processing, which limits the materials, since copper cannot be used so early in the flow; refractory metals such as tungsten or ruthenium are candidates.

## What buried rails change in the flow

Buried rails connect to source and drain through a via from below, and to the power grid either from the front through vias that pass beside the devices or from the back through nano-through-silicon vias. That second option is a step toward full backside power delivery. For design tools, BPR changes the cell template, the extraction stack, the power-grid model and possibly the placement rules near rail taps. A DTCO study weighs the area and IR gains against new process steps, reliability risks of a metal near the transistors, and the cost of the new masks.

## Going deeper

Ryckaert et al., papers on buried power rails and scaling boosters from imec in IEDM and VLSI Symposium proceedings; Gupta et al., "Buried power rail integration with FinFETs for ultimate CMOS scaling", IEEE Transactions on Electron Devices (2020).
