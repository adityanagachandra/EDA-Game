# Hierarchical DRC

A modern chip repeats a few thousand cells billions of times. Hierarchical DRC checks each cell once and reuses the result wherever the surroundings are the same, checking only the regions where instances interact. Done well, it cuts runtime and memory by orders of magnitude; done naively, it reports false errors or misses real ones. The central difficulty is the context problem.

## Reporting errors once per cell

A flat run sees every copy of every shape, so an error inside a cell is reported once per placement. If cell A has 5 internal violations and is placed 47 times, cell B has 5 and is placed 30 times, and the top cell adds 6, a flat run reports $5 \cdot 47 + 5 \cdot 30 + 6 = 391$ violations. A hierarchical run reports each cell's errors once, in the cell's coordinates: $5 + 5 + 6 = 16$. The smaller report is also the more useful one, because one edit to the cell master fixes all placements, provided the error does not depend on context.

## Interaction regions between instances

Instances interact where their shapes come within the largest rule distance of each other. In an array of cells with metal along all four boundary edges, horizontal neighbours interact if the x gap between cells is below the spacing rule $S$, vertical neighbours if the y gap is, and diagonal neighbours if $\sqrt{g_x^2 + g_y^2} < S$. A $c \times r$ array has $(c-1)r$ horizontal, $c(r-1)$ vertical and $2(c-1)(r-1)$ diagonal neighbour pairs. Each interacting pair is a region the engine checks with shapes from both instances. Regular arrays produce few distinct interaction types, which the engine checks once each.

## The context problem

A cell clean in isolation can violate a rule in context: a top-level wire passes 20 nm from a cell shape, or two overlapping instances bring shapes together. The engine must therefore check each instance's boundary region together with its surroundings. When placements see different surroundings, engines create variants of the cell or process interaction regions separately. Hierarchy pays off when cells interact only through thin boundary regions; layouts with many overlaps or tiny cells placed once gain nothing and can run slower than flat. KLayout's deep mode and commercial hierarchical engines differ mainly in how they manage contexts.
