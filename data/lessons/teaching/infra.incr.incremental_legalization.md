# Incremental legalization

After global placement, a legalizer snaps every cell to a row and a site without overlap. During optimization, though, the design changes a little at a time: a resizer widens a cell, a buffer is inserted, a cell is moved. Re-legalizing the whole design for each change would scramble a carefully optimized placement and take far too long. Incremental legalization repairs only the neighbourhood of the change, moving as few cells as possible by as little as possible, which keeps timing and wirelength close to what the optimizer saw.

## Push-right shifting in a row

The simplest incremental repair works inside one row. Cells are kept sorted by x, each with a left edge $x_i$ and width $w_i$ in site units. When cell $k$ grows by $g$ sites with its left edge fixed, its right edge moves from $x_k + w_k$ to $x_k + w_k + g$. Walk right: each later cell must start at or after the right edge of its left neighbour, so its new position is $x'_i = \max(x_i,\ x'_{i-1} + w_{i-1})$. The walk stops changing anything at the first cell whose gap absorbs the push.

Example: cells at $x = 0$ (width 4), $x = 4$ (width 3), $x = 9$ (width 2) and $x = 11$ (width 5). Widen the first by 4. It now ends at 8, so the second moves to 8 and ends at 11. The third, at 9, must move to 11 and ends at 13. The fourth, at 11, moves to 13. Three cells move.

## Where a pushed cell lands

Each pushed cell lands exactly against its left neighbour, because the rule moves it by the least amount that removes the overlap. So positions in a pushed run are a running sum: the first pushed cell starts at the new right edge of the edited cell, and every following pushed cell starts at the right edge of the one before it, until a gap at least as large as the remaining push appears.

In the example, the pushed run is 8, 11, 13: the second cell at 8, the third at $8 + 3 = 11$, the fourth at $11 + 2 = 13$. Gaps before the edit were 0 between the first two cells, 2 between the second and third, and 0 between the third and fourth, so the push of 4 shrank to 2 after the first gap and stayed 2. A useful check: the shift of each pushed cell equals the push it received minus the gap in front of it.

## Total displacement and when to stop

The cost of an incremental repair is usually measured as total displacement, the sum of $|x'_i - x_i|$ over moved cells, because every moved cell perturbs its nets and the timing the optimizer just evaluated. In the example the shifts are 4, 2 and 2, a total of 8 sites.

Pushing in one direction is cheap but not minimal: a run might be better split, with some cells moving left into earlier gaps, which is what a dynamic-programming legalizer such as Abacus does when it re-solves the cluster positions that minimize quadratic displacement. A push can also run off the end of the row, and then cells must move to a neighbouring row. Production legalizers, such as OpenROAD's dpl, therefore search for nearby legal sites within a window and fall back to wider repairs only when a window has no room. The engineering rule is to bound the region touched, so the update stays local and the placement keeps its quality.

## Going deeper

Spindler, Schlichtmann and Johannes, "Abacus: Fast Legalization of Standard Cell Circuits with Minimal Movement" (ISPD 2008); the OpenROAD documentation for dpl (detailed placement).
