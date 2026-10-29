# Slab decomposition versus minimum partition

The simplest way to cut a rectilinear polygon into rectangles is to slice it with horizontal lines from every reflex vertex: a sweep produces it directly, it is easy to make robust, and mask fracturers and trapezoid-based Boolean engines use it. It is not minimal. Knowing how far it is from the minimum, and when a vertical slicing does better, is the practical side of rectangle decomposition.

## Horizontal slabs

Extend a horizontal cut from every reflex vertex into the interior until it meets the boundary. Between consecutive vertex y-values the polygon is a set of horizontal runs, and runs with the same x-extent in consecutive bands that are not separated by a cut stack into one rectangle. A sweep from bottom to top produces the pieces directly: start a rectangle when a run appears or changes, end it when the run ends or changes. Example: an H-shape (two posts joined by a crossbar) gives five rectangles: the two posts are cut above and below the crossbar, and the crossbar band is one full-width piece. A U-shape gives three, the same as the minimum. Each reflex vertex adds at most one cut, so at most $r + 1$ pieces without holes.

## The gap to the minimum

Slab cuts resolve each reflex vertex separately and ignore chords, especially vertical ones, so slabs can use far more rectangles than the minimum $r - g - h + 1$. Example: the H-shape has $r = 4$ and two vertical chords along the crossbar's sides; the minimum is $4 - 2 + 1 = 3$ (two posts and the crossbar), while slabs give 5. A comb with $t$ teeth on a vertical spine gives $2t - 1$ slabs against a minimum of $t + 1$. For mask writing the extra pieces cost write time, and for extraction or fill they cost memory and run time. The slab result is still valuable: it is canonical (unique for a given polygon), trivial to verify, and compatible with sweep-based engines.

## Choosing the slicing direction

Slicing vertically instead (cuts up and down from every reflex vertex) gives a different count, and a tool can compute both and keep the better. Example: the H-shape sliced vertically gives the two posts and the crossbar, 3 rectangles, which here equals the minimum, while horizontal slicing gives 5. A comb with teeth to the right of a vertical spine behaves the same way. For a plus shape both directions give 3. Some shapes defeat both directions, needing a mix of horizontal and vertical cuts, which is exactly what the matching-based minimum finds. Production fracturers add more criteria than count: they avoid slivers narrower than the writer's minimum figure size, prefer cuts on existing edges and keep runs aligned between neighbouring polygons.

## Going deeper

Eppstein, "Graph-theoretic solutions to computational geometry problems" (2009); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), material on mask data preparation and fracturing.
