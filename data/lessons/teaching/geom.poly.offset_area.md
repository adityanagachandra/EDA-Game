# Area after offsetting a rectilinear polygon

Fill engines, keep-out generators and density estimators need the area of a grown or shrunk polygon without always building it. For rectilinear polygons there is an exact closed form in terms of area, perimeter and corner counts, valid as long as the offset does not make distant parts of the shape interact. It turns sizing into arithmetic, and the conditions under which it breaks are exactly the conditions under which sizing changes topology.

## The A plus P d plus 4 d squared law

Grow a rectilinear polygon with area $A$ and perimeter $P$ by $d$ with square corners. Each edge of length $\ell$ sweeps an outward strip of area $\ell d$. Each convex corner adds a $d \times d$ square; at each reflex corner two strips overlap in a $d \times d$ square that was counted twice. A rectilinear polygon always has exactly 4 more convex than reflex corners, so

$$A_{\text{grown}} = A + P d + 4 d^2.$$

Shrinking by $d$ gives $A - P d + 4 d^2$ by the mirror argument, provided no part is narrower than $2d$. Example: a polygon with $A = 120$, $P = 52$ grown by 2 has area $120 + 104 + 16 = 240$. A $10 \times 6$ rectangle checks the law: $(14)(10) = 140 = 60 + 32 \cdot 2 + 16$. The law needs only two numbers per polygon, so a fill or keep-out estimator can predict grown areas without building any geometry.

## Applying it to a vertex list

Given vertices, compute $A$ by the shoelace formula and $P$ as the sum of edge lengths $|\Delta x| + |\Delta y|$, then apply the law. Example: the polygon $(0,0), (5,0), (5,2), (9,2), (9,7), (0,7)$ has area $5 \cdot 7 + 4 \cdot 5 = 55$ and perimeter $5 + 2 + 4 + 5 + 9 + 7 = 32$; grown by 1 it covers $55 + 32 + 4 = 91$. The law needs the offset not to make separate parts meet. Outward, that means no pocket of the exterior is narrower than $2d$: inside a narrower notch the two wall strips overlap and the law overcounts. Inward, no neck or arm may be narrower than $2d$. These are exactly the cases where a closing or opening would change the shape, and an engine should check them rather than trust the formula.

## Round corners and the Euclidean offset

A Euclidean offset adds every point within distance $d$ of the polygon. Edges still sweep strips of area $\ell d$ and reflex corners still remove a $d \times d$ overlap, but each convex corner now adds a quarter disc of area $\pi d^2/4$ instead of a square. With $c$ convex and $r$ reflex corners ($c = r + 4$),

$$A_{\text{grown}} = A + P d + \left(\frac{\pi c}{4} - r\right) d^2.$$

For a convex polygon this is $A + Pd + \pi d^2$, the classical Steiner formula. Example: $A = 55$, $P = 32$, $r = 1$, $c = 5$, $d = 1$ gives $55 + 32 + 3.927 - 1 = 89.93$. The law also needs every edge next to a reflex corner to be at least $d$ long, or a quarter disc spills into a neighbouring strip.

## Going deeper

Santalo, "Integral Geometry and Geometric Probability" (1976), on Steiner's formula for parallel bodies; the KLayout documentation on sizing modes and corner handling.
