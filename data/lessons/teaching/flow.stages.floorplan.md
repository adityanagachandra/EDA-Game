# Floorplanning and power planning

Floorplanning is the first physical step. It chooses the die and core size, places the large blocks (memories, analog macros, hard IP), arranges the I/O, and lays down the power distribution network. Every later physical stage works inside the floorplan, so a bad floorplan, such as a memory in the wrong corner or a starved power grid, cannot be fixed by a better placer or router.

## Core area from utilization

The core is the region that holds standard cells and macros, inside the I/O ring. Its size starts from the cell area reported by synthesis and a target utilization, the fraction of core area occupied by cells. Placement and routing need white space for buffers, wires and later fixes, so typical targets are well below 100%. If the cells and macros cover $A_{cell} = 1{,}200{,}000\ \mu m^2$ and the target utilization is $u = 0.7$, the core area is $A_{core} = A_{cell}/u \approx 1{,}714{,}286\ \mu m^2$. Choosing too high a utilization saves area but makes routing and timing closure fail later.

## Shape, rows and the die

With an aspect ratio $r = H/W$, the core width is $W = \sqrt{A_{core}/r}$ and height $H = rW$. For a square core in the example, $W = H \approx 1309.3\ \mu m$. The core is divided into standard-cell rows whose height equals the library's cell height. The die adds a margin for the I/O ring and seal: with a margin $m$ on each side, the die is $(W + 2m)(H + 2m)$. With $m = 100\ \mu m$ the die is about $1509.3 \times 1509.3\ \mu m$, or about $2.278\ mm^2$.

## Macros and the power grid

Macros are placed before standard cells. Good practice puts them near the edges, leaves channels between them for routing, and orients their pins toward the logic they talk to. A poor macro placement creates long wires and routing congestion that no later stage can fix. Power planning builds the grid: power rings around the core, straps on upper metal layers, and rails along each row. The grid must keep the voltage drop (IR drop) at every cell within budget while consuming as little routing resource as possible, so its density is a trade between signal routability and supply integrity.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3; the OpenROAD documentation for ifp, the macro placer and pdngen.
