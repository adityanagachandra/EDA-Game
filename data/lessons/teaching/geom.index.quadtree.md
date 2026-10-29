# Quadtrees over the die

A quadtree recursively divides a square region into four equal quadrants, stopping where a cell is simple enough: empty, or holding few objects. Over a die of side $W = 2^k$ the cell sizes are powers of two, so every cell index is a bit-shift of a coordinate. Quadtrees index points, rectangles and rasterized layers, and the same hierarchy underlies multi-resolution density maps and hierarchical routing grids.

## Cell size by depth

At depth $d$ (the root at depth 0) a quadtree over a $W \times W$ die has cells of side $W / 2^d$, and there are $4^d$ of them if the tree is complete. To get cells no larger than a target size $s$, the depth must satisfy $W / 2^d \le s$, so $d = \lceil \log_2 (W/s) \rceil$. Example: a die of 4096 units with features of 100 units needs $d = 6$, because $4096/32 = 128 > 100$ while $4096/64 = 64 \le 100$. In practice quadtrees are adaptive: only cells that contain too many objects subdivide, so empty regions stay coarse. That makes them compact on clustered layouts where a uniform grid would waste memory on empty bins.

## Separating two points

In a point-region quadtree, a cell subdivides until each point has its own cell. At depth $d$ the cell index of a point is $(\lfloor x / c \rfloor, \lfloor y / c \rfloor)$ with $c = W/2^d$; two points share a cell until those indices first differ. With $W = 2^k$, the index is the top $d$ bits of each coordinate, so the separation depth is set by the first leading bit in which $x$ or $y$ differ. Example: in $[0, 256)$, points $(100, 40)$ and $(110, 45)$ share cells while $c \ge 16$, since $\lfloor 100/16 \rfloor = \lfloor 110/16 \rfloor = 6$, and separate at $c = 8$, depth 5, where $12 \ne 13$. Close points force deep trees, which is why point quadtrees limit depth or store several points per leaf.

## Rectangles across cells

A rectangle stored in a quadtree either goes into the smallest cell that contains it entirely (the MX-CIF scheme) or is registered in every cell it overlaps at a chosen depth. For half-open cells of side $c$ and a rectangle $[x_1, x_2) \times [y_1, y_2)$, the overlapped columns run from $\lfloor x_1/c \rfloor$ to $\lfloor (x_2 - 1)/c \rfloor$ and similarly for rows. Example: cells of side 32 and the rectangle $[20, 75) \times [64, 70)$ span columns 0 to 2 and row 2 only, so 3 cells. Storing in every overlapped cell duplicates large shapes; storing in the smallest enclosing cell avoids duplication but puts shapes that straddle a major boundary near the root, where every query meets them.

## Going deeper

Samet, "Foundations of Multidimensional and Metric Data Structures" (2006), chapters 1 and 3; Finkel and Bentley, "Quad trees: a data structure for retrieval on composite keys" (Acta Informatica, 1974).
