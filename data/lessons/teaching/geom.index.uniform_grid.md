# Uniform grids and bins

The simplest spatial index divides the die into equal bins and keeps, for each bin, the shapes that overlap it. Lookup is arithmetic, insertion and deletion are trivial, and on evenly spread layout data it is hard to beat. Placers, routers and DRC engines all use bin grids. Their weaknesses, duplicated large shapes and uneven bin loads, are equally simple to quantify.

## Bin index arithmetic

With square bins of side $b$, $C$ columns and $R$ rows, the bin holding point $(x, y)$ has column $\lfloor x/b \rfloor$ and row $\lfloor y/b \rfloor$, and in row-major numbering its index is $\text{row} \cdot C + \text{col}$. Bins are half-open so that every point belongs to exactly one bin. Example: bins of 50 nm, 12 columns, point $(437, 268)$: column 8, row 5, index $5 \cdot 12 + 8 = 68$. No tree descent, no comparisons beyond two divisions; with $b$ a power of two the divisions become shifts. This is why global placers and routers use grids of bins for density and congestion: the bin of every cell or wire segment is known in constant time.

## Shapes in several bins

A rectangle is registered in every bin it overlaps. With half-open coordinates its columns run from $\lfloor x_1/b \rfloor$ to $\lfloor (x_2 - 1)/b \rfloor$ and its rows likewise, and it occupies (columns) times (rows) bins. Example: with $b = 25$ the shape $[40, 90) \times [10, 30)$ spans columns 1 to 3 and rows 0 to 1, six bins. Total storage is the sum over shapes, so a few long wires or large blocks can dominate memory. A query that collects shapes from several bins sees a long shape more than once and must de-duplicate, for example by reporting a pair only from the bin containing a reference corner of their intersection. Choosing $b$ near the typical shape size keeps duplication low.

## Expected candidates per query

If $n$ small shapes are spread uniformly over $C \times R$ bins, each bin holds $n/(CR)$ shapes on average, and a query window touching $q$ bins fetches about $q n/(CR)$ candidates, which are then tested exactly. Example: 50,000 shapes over $10 \times 10$ bins and a query touching 4 bins: $4 \cdot 500 = 2000$ candidates. Making bins smaller reduces candidates per query but increases duplication and the number of bins a query touches; the sweet spot keeps a constant number of shapes per bin. The average hides the weakness: layouts are clustered, so some bins hold thousands of shapes while others are empty, and hierarchical structures such as quadtrees and R-trees adapt where a grid cannot.

## Going deeper

Samet, "Foundations of Multidimensional and Metric Data Structures" (2006), chapter 1; the OpenROAD documentation for its global placer and global router, which use bin grids for density and congestion.
