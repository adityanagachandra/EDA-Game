# H-trees

The H-tree is the classic symmetric clock topology: an H whose four tips each feed a smaller H, recursively. Because every root-to-sink path has the same shape, nominal skew is zero for uniformly distributed sinks. H-trees are still used for the top level of clock distribution in high-performance designs, and OpenROAD's CTS builds H-tree-style topologies over sink clusters. The geometry also gives clean closed forms for wirelength, which make it a good reference point for any other network.

## The recursive geometry

Take a square die of side $D$. The level-1 H has a horizontal bar of length $D/2$ through the centre, and two vertical bars of length $D/2$ at its ends, so its total length is $1.5D$ and its four tips are the centres of the four quadrants. Each tip is the centre of a level-2 H of half the size, total $0.75D$, and so on. After $n$ levels there are $4^n$ tips, each the centre of a cell of side $D/2^n$, which is where the sinks or the next stage of distribution attach.

## Wirelength level by level

Level $k$ has $4^{k-1}$ H's, each of length $1.5D/2^{k-1}$, so the level's total is $1.5D \cdot 2^{k-1}$: wirelength doubles every level, because the count quadruples while the size halves. Summing the geometric series over $n$ levels gives $1.5D(2^n - 1)$. For $D = 8000$ µm and $n = 2$ that is $1.5 \cdot 8000 \cdot 3 = 36\,000$ µm; level 4 alone of a 9000 µm H-tree is $1.5 \cdot 9000 \cdot 8 = 108\,000$ µm. The doubling tells you that the bottom levels dominate wire and power, which is why the last levels are often replaced by local trees.

## Equal path length and its limits

Within a level-$k$ H, a path runs half a horizontal bar and half a vertical bar, $D/2^{k+1} + D/2^{k+1} = D/2^k$. Over $n$ levels every source-to-sink path has length $D(1 - 2^{-n})$, so for $D = 10\,000$ µm and four levels each path is 9375 µm. Equal length gives zero nominal skew, but only nominally: real sinks are not uniform, loads differ, blockages force detours, and process variation affects each branch differently. That is why H-trees are tuned, combined with meshes, or replaced by algorithms such as DME that adapt to the actual sink locations.

## Going deeper

Bakoglu, "Circuits, Interconnections, and Packaging for VLSI" (1990), the chapter on clock distribution; the OpenROAD documentation for the CTS module.
