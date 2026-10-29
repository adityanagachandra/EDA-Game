# The interconnect stack

Above the transistors sit ten to twenty wiring levels. The lowest are thin, dense and resistive, for local connections inside and between cells; higher levels are progressively wider and thicker, for long routes, clocks and power. The stack's pitches and directions define the routing resource every placer and router reasons about, and its cross-sections set the resistance and capacitance that timing sees.

## Layers, pitches and directions

Each metal layer has a pitch, a preferred direction and a thickness. Adjacent layers alternate horizontal and vertical so that a route can turn by dropping through one via. Layers are grouped by pitch: the tightest local layers (often called 1x) near the cells, then intermediate and semi-global layers at two or four times that pitch, then thick global layers for power and clock. The number of routing tracks a layer supplies across a block is the block dimension divided by the pitch: a 40 um wide block on a vertical layer at 32 nm pitch offers $\lfloor 40000/32 \rfloor = 1250$ tracks.

## Routing supply across the stack

A router's capacity model counts tracks per layer, per direction, per region. Horizontal layers supply tracks across the block's height; vertical layers across its width. For a 20 um by 30 um block with M2 horizontal at 28 nm, M3 vertical at 36 nm and M4 horizontal at 40 nm, supply is $\lfloor 30000/28 \rfloor + \lfloor 20000/36 \rfloor + \lfloor 30000/40 \rfloor = 1071 + 555 + 750 = 2376$ tracks. Some of that supply is consumed by power rails, straps, pin access and blockages, so routers work with effective capacity. Global routers divide the block into tiles and compare demand against these per-tile capacities to estimate congestion.

## Why long wires go up

At fixed aspect ratio and width equal to half the pitch, the wire cross-section scales with pitch squared, so resistance per unit length scales as $1/p^2$. A 24 nm pitch layer has $(96/24)^2 = 16$ times the resistance per length of a 96 nm pitch layer. Capacitance per length changes far less because it depends on ratios of dimensions. Long routes therefore belong on upper layers, at the cost of via stacks to reach them and of fewer tracks. Layer assignment, deciding which nets use which layers, is a timing-driven optimization in routers and in buffering, and the stack definition in the PDK's technology files is its input.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapters on global routing; the IRDS interconnect tables for typical stack pitches.
