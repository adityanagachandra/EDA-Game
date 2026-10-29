# Markowitz pivot selection

Circuit matrices from modified nodal analysis are unsymmetric, and their pivots must also be numerically safe. Markowitz pivoting chooses, at each step, the pivot that promises the least fill, using only row and column counts. It is the strategy of the original SPICE sparse solver and still a reference point for every unsymmetric sparse LU in EDA.

## The Markowitz count

At a step of elimination, let $r_i$ be the number of nonzeros in row $i$ of the active submatrix and $c_j$ the number in column $j$. The Markowitz count of a candidate pivot $a_{ij}$ is
$$M_{ij} = (r_i - 1)(c_j - 1).$$
Pivoting on $a_{ij}$ updates every position $(p, q)$ with $p$ a row holding a nonzero in column $j$ and $q$ a column holding a nonzero in row $i$, excluding the pivot row and column. There are exactly $(r_i - 1)(c_j - 1)$ such positions.

Example: row 2 has nonzeros in columns 1, 2 and 4 (so $r_2 = 3$); column 4 has nonzeros in rows 2 and 3 ($c_4 = 2$). Then $M_{24} = 2 \cdot 1 = 2$. A singleton row or column gives $M = 0$: pivoting there creates no fill at all.

## Choosing the pivot with least predicted fill

A Markowitz search scans candidate entries and takes the one with the smallest count. Searching every entry is expensive, so implementations keep rows and columns bucketed by count and look first at the shortest ones, stopping early when a zero or small count appears.

For the pattern row 1: columns 1, 2; row 2: columns 2, 3, 4; row 3: columns 1, 3; row 4: column 4, entry (4, 4) has $r_4 = 1$, so its count is 0 and it is chosen first. After each pivot the counts change, because fill adds entries and the pivot row and column leave the active matrix, so the search is repeated on the updated submatrix.

Kundert's Sparse 1.3, used by SPICE3 and ngspice, combines this search with a stability threshold, described in threshold partial pivoting.

## Markowitz count as a bound on actual fill

The count is an upper bound. Fill occurs only at those of the $(r_i - 1)(c_j - 1)$ updated positions that were zero; positions already nonzero are just modified. So actual fill $\le M_{ij}$.

Example: pivot $(1, 1)$ with row 1 holding columns 1, 2, 3 and column 1 holding rows 1, 3, 4. The updated positions are $(3, 2), (3, 3), (4, 2), (4, 3)$, four of them, matching $M = 2 \cdot 2 = 4$. If $(3, 3)$ and $(4, 2)$ were already nonzero, the step creates only 2 fill entries.

Computing exact fill for every candidate costs much more than the counts, which is why the bound is used. It is a local, greedy rule: the pivot with least fill now can still lead to more fill later, as with every greedy ordering.

## Going deeper

Markowitz, "The elimination form of the inverse and its application to linear programming", Management Science (1957); Duff, Erisman and Reid, "Direct Methods for Sparse Matrices" (2nd ed., Oxford, 2017), chapter 7.
