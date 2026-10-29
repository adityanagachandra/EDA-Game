# Standard-cell layout

A standard cell is a gate laid out on a fixed template: a common height, power rails along the top and bottom, and a width that is a whole number of placement sites. That regularity is what lets a placer treat millions of cells as rectangles in rows and lets a router assume where pins and rails are. Track height, poly pitch and rail sharing are the numbers that tie the transistor world to the placer's grid.

## The template: rails, wells and devices

Every cell in a library shares one height. A VDD rail runs along one horizontal edge and a VSS rail along the other, so when cells abut in a row their rails join into continuous stripes. pMOS devices sit in an n-well on the VDD side and nMOS devices on the VSS side, with gate poly running vertically across both, so an input's nMOS and pMOS share one poly line. Diffusion is shared between neighbouring transistors wherever the netlist allows, which is why cell layout searches for an ordering of transistors that maximizes diffusion sharing (an Euler path through the pull-up and pull-down graphs). In FinFET processes, device widths come in whole fins, so a transistor's strength is quantized.

## Height in tracks, width in poly pitches

Cell height is quoted in tracks: the number of horizontal routing pitches of the lower metal that fit in the height. A 7.5-track library on a 28 nm metal pitch is $7.5 \times 28 = 210$ nm tall. Cell width is quantized by the contacted poly pitch (CPP), the spacing of gate lines: an inverter might be 3 CPP wide and a flip-flop 20 or more. With a 51 nm CPP, a cell 6 CPP wide in that library has area $6 \times 51 \times 210 = 64{,}260$ square nanometres, or 0.0643 square micrometres. Fewer tracks give denser libraries but leave less room for devices and for in-cell routing and pin access, which is a central trade-off in library architecture.

## What the template buys the tools

Because every cell has the same height and rail positions, a placer can model the core as rows of sites and a cell as an integer-width rectangle, and legalization becomes snapping cells to sites without overlaps. Abutting cells share rails automatically, so power routing for standard cells is just connecting the row rails to a grid. Pins sit on known tracks, so a router can find access points from the abstract view without seeing transistors. The cost of regularity is some wasted area inside each cell compared with custom layout, and limited flexibility: a cell that needs more devices than fit in the height must be made wider or built as a double-height cell spanning two rows, which placers must then handle as a special case.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 1.5 and 14.2; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 1.
