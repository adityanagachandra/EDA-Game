# Hierarchy with SREF and AREF

A layout is a directed acyclic graph of cells: a cell holds shapes and placements of other cells. GDSII expresses placements as SREF (one instance) and AREF (a regular array). Hierarchy is compression: a bit cell stored once may appear a million times. Every verification engine must decide what to flatten and what to keep, and that decision starts with counting what the hierarchy expands to.

## Flattening a cell tree

Flattening replaces each placement by a transformed copy of the placed cell's contents, recursively. Counts compose bottom-up. If LEAF holds 25 shapes, MID holds 11 shapes and 9 SREFs of LEAF, and TOP holds 4 shapes, 10 SREFs of MID and an AREF of LEAF with 7 columns and 14 rows, then MID expands to $11 + 9 \cdot 25 = 236$ shapes and TOP to $4 + 10 \cdot 236 + 98 \cdot 25 = 4814$. The hierarchy stores only $25 + 11 + 4 = 40$ shapes. A post-order traversal of the cell DAG computes every cell's flat count once, in time linear in the number of cells and placements.

## Placements as paths in the cell DAG

Each flat occurrence of a cell corresponds to one path from the top cell to it in the cell DAG, with edges weighted by instance counts. In the example, LEAF appears $10 \cdot 9 = 90$ times through MID and $7 \cdot 14 = 98$ times through the AREF, 188 placements in all. The count is a sum over paths of the product of edge weights, the same computation as counting paths in a DAG. It is what a flat DRC run pays for each LEAF-level check, and what a hierarchical run can avoid when contexts repeat.

## AREF versus SREF

An SREF costs SREF 4 + SNAME (4 plus the padded name) + XY 12 + ENDEL 4 bytes. An AREF costs AREF 4 + SNAME + COLROW 8 + XY 28 + ENDEL 4 bytes, independent of the array size. For DFFR (4 characters) placed in a $17 \times 22$ grid, 374 SREFs cost $374 \cdot 28 = 10472$ bytes while one AREF costs 52, saving 10420. AREFs also tell a tool that the placements are regular, which a hierarchical engine can exploit: interior instances share one context. The cost is rigidity: an AREF must be a perfect lattice, so one displaced instance breaks it into pieces.
