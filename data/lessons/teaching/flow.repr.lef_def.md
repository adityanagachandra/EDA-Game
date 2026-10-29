# LEF and DEF physical exchange

LEF (Library Exchange Format) and DEF (Design Exchange Format) are the text formats in which physical-design tools exchange abstract layout. LEF describes the technology and the library: routing layers, pitches, vias and design rules, plus each cell's outline, pins and obstructions. DEF describes one design: die area, rows, placed components, pins, and routed nets. Placement and routing engines, including OpenROAD, read and write them.

## LEF: abstracts of cells and layers

A technology LEF lists each routing layer with its preferred direction, track pitch, minimum width and spacing, plus via definitions. A cell LEF describes each macro by its class (core cell, block, pad), its size, the site it occupies, its pin shapes on specific layers, and obstructions (OBS) where a router must not draw. It deliberately omits the cell's internal transistors and most internal wiring. That abstraction is exactly what a placer and router need: where a cell may sit, how big it is, and where its pins can be reached.

## DEF: a design in database units

DEF stores coordinates as integers in database units, declared by a line such as `UNITS DISTANCE MICRONS 1000`, meaning 1000 units per micron. A `DIEAREA ( 0 0 ) ( 1500000 1200000 )` is therefore a 1500 by 1200 micron die, an area of $1.8\ mm^2$. Each component appears with its master and location, for example `- u42 NAND2_X1 + PLACED ( 84000 212800 ) N`, which is at $(84, 212.8)\ \mu m$ with north orientation. Rows are declared with their site and repeat count, and nets carry routed segments once routing is done.

## Rows, sites and capacity

Standard cells sit on a site grid. If the core is 1360 $\mu m$ tall and the site height is 2.72 $\mu m$, the core holds $1360 / 2.72 = 500$ rows. If each row is 1200 $\mu m$ wide with a site width of 0.19 $\mu m$, a row has $\lfloor 1200 / 0.19 \rfloor = 6315$ sites. A cell's width in LEF is a whole number of sites. These counts bound placement capacity, and they show why LEF and DEF must agree: a DEF written with a different site or unit declaration than the LEF it is read with silently misplaces every cell.

## Going deeper

The LEF/DEF Language Reference (Cadence, distributed through Si2); the OpenDB documentation in OpenROAD.
