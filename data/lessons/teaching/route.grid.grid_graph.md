# The global-routing grid graph

The global router's data structure is a grid graph over gcells, in 2D with summed capacities or in 3D with one copy per layer joined by vias. Its size sets memory and runtime, its edges carry capacity and demand, and blockages are expressed by lowering edge capacities. Counting its parts precisely is how a tool designer sizes arrays and estimates cost before writing a line of search code.

## The 2D gcell graph

Tile the die into $X$ columns and $Y$ rows of gcells. The 2D graph has $XY$ vertices, $(X-1)Y$ horizontal edges and $X(Y-1)$ vertical edges. Each edge stores a capacity (summed over the layers of its direction), a demand (wires currently routed across it) and, during negotiation, a history cost.

Example: a 17 by 102 gcell grid has $16 \cdot 102 = 1632$ horizontal and $17 \cdot 101 = 1717$ vertical edges, 3349 in all. A large block may have a grid of thousands by thousands of gcells, so edges number in the millions, and per-edge storage is a few integers. Routers store edges in arrays indexed by position rather than as an adjacency list, because the grid is regular.

## The 3D graph with layers

A 3D global-routing graph keeps one grid per layer. Because each layer has a preferred direction, a horizontal layer has only the $(X-1)Y$ horizontal edges and a vertical layer only the $X(Y-1)$ vertical ones. Via edges join each gcell to the same gcell on the next layer: $(L-1)XY$ via edges for $L$ layers.

Example: $41 \times 16$ gcells on 4 layers, the lowest vertical. Two vertical layers give $2 \cdot 41 \cdot 15 = 1230$ wire edges, two horizontal layers $2 \cdot 40 \cdot 16 = 1280$, and three via levels $3 \cdot 656 = 1968$, so 4478 edges. Via edges usually carry a capacity too, since too many vias in one gcell block tracks on every layer they cross. Many global routers route on the 2D projection and assign layers afterwards because the 2D graph is $L$ times smaller.

## Blockages on the graph

A macro covering an $a \times b$ block of gcells lowers the capacity of every edge it touches, typically to zero on the layers it blocks. Inside the block there are $(a-1)b + a(b-1)$ edges, and if the block sits away from the die edge, $2a + 2b$ more edges cross its boundary. Both sets must be updated, because a wire crossing the boundary needs tracks inside the macro's first gcell.

Example: an 11 by 12 macro touches $10 \cdot 12 + 11 \cdot 11 = 241$ internal edges and $2 \cdot 11 + 2 \cdot 12 = 46$ boundary edges, 287 in all. A partial obstruction instead scales capacity by the fraction of tracks it blocks. Getting this wrong in either direction shows up later as phantom congestion or as overflow the global router never saw.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the FastRoute papers by Pan and Chu for 2D global routing with later layer assignment.
