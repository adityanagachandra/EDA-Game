# Width checks by shrinking and regrowing

A minimum-width rule says no part of a shape may be narrower than $w$. Width is not a property of a rectangle but of every neck and arm of a merged polygon, which makes it harder than it sounds. Morphology gives a clean implementation: an opening (shrink then grow by $w/2$) removes exactly the parts too narrow to hold a $w \times w$ square, and the removed area marks the violations. The same sizing trick, with AND, implements spacing.

## Opening removes narrow shapes

Shrink a shape by $d$ and then grow the result by $d$. A rectangle with both sides greater than $2d$ shrinks to a smaller rectangle and grows back exactly; one with a side below $2d$ vanishes in the shrink and nothing grows back. So for isolated rectangles the area lost to an opening by $d$ is the total area of those narrower than $2d$, and those are exactly the width violations for a rule of $2d$. Example: with $d = 3$, rectangles of $10 \times 4$, $8 \times 9$ and $5 \times 12$ lose the first and third (4 and 5 are below 6), area $40 + 60 = 100$. Engines then XOR the original with the opened result to produce violation markers. A side of exactly $2d$ is the boundary case, which the engine's open or closed convention decides.

## Narrow arms of one polygon

For a merged polygon, the opening keeps every point that some $2d \times 2d$ square inside the polygon covers. A thick body survives whole, while an arm narrower than $2d$ cannot hold the square anywhere, not even straddling the junction, because the square would stick out of the arm; so the arm disappears and leaves the body. An arm wider than $2d$ survives entirely. Example: a $20 \times 16$ body with a $12 \times 3$ arm, opened with $d = 2$: the arm is 3 wide, below 4, so it vanishes, and the area left is 320. This is why width must be checked on merged polygons: two overlapping drawn rectangles can each be wide while their union has a narrow neck, or each be narrow while their union is wide.

## Spacing by growing and AND

To find spacing violations with rule $s$, grow every shape by $s/2$ and AND shapes from different polygons: a separated pair with gap $g < s$ now overlaps. For two shapes facing each other across an x-gap $g$, with y-ranges overlapping by $L$, the grown shapes overlap by $s - g$ in x and by $L + s$ in y (each y-range grew by $s/2$ at both ends), so the AND has area $(s - g)(L + s)$. Example: $s = 20$, $g = 12$, $L = 30$ gives an $8 \times 50$ marker of area 400. Growing with square corners checks the square metric; a Euclidean check either grows with rounded corners or confirms corner candidates exactly. The marker's size is not the violation's size, so tools report the facing edges as well.

## Going deeper

Serra, "Image Analysis and Mathematical Morphology" (1982); the KLayout DRC reference on width, space and sized.
