# Routing grids as graphs

A router does not search the plane; it searches a graph. Global routers use a coarse grid of routing bins, and detailed routers use a fine grid of tracks on each metal layer, with vias between layers. Every routing algorithm's cost is set by the size of that graph, so a router designer must be able to count its vertices and edges and know what each modeling choice, preferred directions in particular, does to them.

## Vertices of a layered grid

A detailed-routing grid places a vertex at every intersection of a horizontal and a vertical track, on every metal layer. With $W$ columns, $H$ rows and $L$ layers there are $W \cdot H \cdot L$ vertices. A global-routing grid is the same structure at coarser granularity, where each vertex is a bin (a g-cell) holding many tracks, and edges carry capacities equal to the number of tracks crossing the bin boundary.

A block with 7 columns, 21 rows and 6 layers has $7 \times 21 \times 6 = 882$ vertices; a real detailed grid has billions, which is why detailed routers work in windows and global routers coarsen first. Vertices that fall on blockages or pins of other nets are removed or given infinite cost.

## Preferred directions and wire edges

Each metal layer has a **preferred direction**, alternating between horizontal and vertical from one layer to the next. Wire edges on a layer join neighbouring vertices only along its preferred direction: a horizontal layer has $(W - 1) H$ wire edges and a vertical layer $W (H - 1)$. Wrong-way wiring is either forbidden or heavily penalized.

With $W = 21$, $H = 6$ and three layers (horizontal, vertical, horizontal), the two horizontal layers contribute $2 \times 20 \times 6 = 240$ edges and the vertical layer $21 \times 5 = 105$, for 345 wire edges. Preferred directions halve the wire edges compared with a full grid on every layer, and they make the layers' wires parallel, which simplifies design-rule checking and keeps coupling predictable.

## Vias and edge costs

Via edges connect the same $(x, y)$ point on adjacent layers, so there are $W H (L - 1)$ of them. For a grid with 10 columns, 9 rows and three layers, the two horizontal layers give $2 \times 9 \times 9 = 162$ wire edges, the vertical layer gives $10 \times 8 = 80$, and vias add $10 \times 9 \times 2 = 180$, for 422 edges in total.

The router's real model sits in the edge costs: wirelength per unit, a larger cost per via (vias add resistance and hurt yield), congestion penalties on edges near capacity, and history costs in negotiated routing. Changing costs changes the routes without changing the graph. The graph's size, about $3 W H L$ edges, is what a maze search must explore, which motivates the bounded windows and A* guidance that real routers use.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 5; the TritonRoute and FastRoute documentation in OpenROAD.
