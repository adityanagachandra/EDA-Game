# Area of a rectilinear polygon by the shoelace formula

A layout polygon is stored as a ring of integer vertices. Its area, needed for density, capacitance estimates, antenna ratios and fill, comes from the shoelace formula: a sum of cross products of consecutive vertices. The same sum also reveals orientation, which is how databases tell outer boundaries from holes. It is exact on integer data, linear time and the first thing to compute on any polygon.

## Summing cross terms

For a simple polygon with vertices $v_0, \dots, v_{n-1}$ listed in order,

$$A = \frac{1}{2} \sum_{i=0}^{n-1} (x_i y_{i+1} - x_{i+1} y_i),$$

with indices taken modulo $n$. Each term is twice the signed area of the triangle from the origin to an edge, and the triangles outside the polygon cancel. It is the discrete form of Green's theorem, $A = \frac{1}{2}\oint (x\,dy - y\,dx)$. Example: the L-shape $(0,0), (4,0), (4,2), (2,2), (2,5), (0,5)$ gives terms $0, 8, 4, 6, 10, 0$, sum 28, area 14, matching $4 \times 2 + 2 \times 3$. On integer coordinates the sum is an integer, so the area of a rectilinear polygon is computed exactly with no rounding at all. The cost is one multiply-and-subtract per vertex, so the area of every polygon in a layer comes out of a single linear pass.

## The sign gives orientation

The formula returns a signed area: positive when the vertices run counter-clockwise, negative when they run clockwise. A tool that stores outer boundaries counter-clockwise reads a negative result as a hole or as a ring that must be reversed. Reversing the vertex list flips the sign of every term. Example: listing the L-shape above as $(0,5), (2,5), (2,2), (4,2), (4,0), (0,0)$ gives -14. Orientation is not cosmetic. Boolean engines use it to decide which side of each edge is inside, fill rules depend on it, and two shapes with the same vertices but opposite orientation mean different things to an engine that counts winding. Normalizing orientation on input is one of the first steps of any geometry pipeline.

## Holes and multi-ring polygons

A polygon with holes is an outer ring plus inner rings, conventionally with opposite orientation: outer counter-clockwise, holes clockwise. Then the shoelace sum over all rings gives the net area directly, because each hole contributes its area with a negative sign. Example: an outer ring of area 60 with a clockwise $3 \times 2$ hole gives $60 - 6 = 54$. If a file stores every ring counter-clockwise, the tool must subtract hole areas explicitly. The same per-ring sum applies to perimeter (sum of edge lengths over all rings, holes included), which feeds sizing and fringe-capacitance estimates. Because the formula is linear in the rings, it also works incrementally: editing one ring changes only that ring's contribution.

## Going deeper

de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 3; O'Rourke, "Computational Geometry in C" (2nd ed., 1998), chapter 1 on polygon area.
