# Rows, sites and the placement region

Standard cells are not placed anywhere on the die. They sit in horizontal rows inside the core, each row divided into sites, the smallest legal horizontal step. This grid is what makes millions of cells abut cleanly and share power rails, and it is the discrete space that legalization and detailed placement work in. A tool designer must count it, index it and respect it exactly.

## Die, core, rows and sites

The die is the whole chip; the core is the inner region reserved for standard cells and macros, surrounded by the I/O ring or pad area. The core is filled with rows of equal height, the cell height of the library, for example 1.4 um in a 45 nm teaching library or 2.72 um for the SkyWater 130 nm high-density cells. Every standard cell has that height and a width that is an integer number of sites, the site being defined in the technology LEF (0.19 um wide in the 45 nm example, 0.46 um in sky130). Each row starts at an origin, has a site type and a count of sites, and an orientation. Cells in a row share a VDD rail along one edge and a VSS rail along the other, so they connect to power just by abutting.

## Counting rows and sites

With row height $h$ and site width $s$, a core of height $H$ holds $\lfloor H/h \rfloor$ rows and a row of width $W$ holds $\lfloor W/s \rfloor$ sites. A core 140.5 um tall with 1.4 um rows holds $\lfloor 140.5/1.4 \rfloor = 100$ rows, with 0.1 um left over at the top. If the core is 190.1 um wide and sites are 0.19 um, each row holds 1000 sites, and the core holds 100000 sites. A design of 15000 cells averaging 4 sites wide uses 60000 of them, a site utilization of 60%. Floorplanners usually choose the core size so the leftover is zero, but counting with floors is the safe habit, since a partial row or site cannot hold a cell.

## Rows as the legality grid

The rows and sites are the coordinate system of legality. A legal standard cell has its bottom edge on a row boundary, its left edge on a site boundary, and an orientation that matches the row: rows alternate between N and FS (flipped about the horizontal axis) so that the VDD rail of one row coincides with the VDD rail of the next. This makes legalization a discrete problem, an assignment of cells to row and site slots, and it lets a placer store a row as an ordered list of cells with integer site positions. It also defines what the global placer must approximate: a target density below 100% of the sites leaves room for buffers inserted later, for filler and tap cells, and for routing. Rows can be cut by macros and blockages, so a real row is a list of usable segments, not one interval.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3 on chip planning and chapter 4 on placement; the LEF/DEF Language Reference for the SITE and ROW statements.
