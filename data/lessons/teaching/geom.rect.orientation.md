# The orientation predicate

Almost every geometric algorithm, from convex hulls to polygon Booleans to segment intersection, branches on one question: given three points, do they make a left turn, a right turn, or lie on a line? The answer is the sign of a 2 by 2 determinant. It is the most frequently evaluated predicate in a geometry kernel, so its exactness and its cost shape the whole engine.

## The sign of a determinant

For points $a$, $b$, $c$ define

$$\text{orient}(a, b, c) = (b_x - a_x)(c_y - a_y) - (b_y - a_y)(c_x - a_x),$$

the z-component of the cross product of $b - a$ and $c - a$, and twice the signed area of triangle $abc$. It is positive when $a \to b \to c$ turns left (counter-clockwise), negative for a right turn and zero when the points are collinear. Example: $a = (1, 2)$, $b = (5, 3)$, $c = (2, 7)$ gives $(4)(5) - (1)(1) = 19 > 0$, a left turn. Everything else follows: a convex hull keeps left turns, a polygon's orientation is the sign of the sum of such terms, and a point is left of a directed edge when the orientation is positive. Because the value is an exact integer on integer input, the sign can be trusted completely when the arithmetic is wide enough.

## Deciding whether segments cross

Segments $pq$ and $rs$ properly cross exactly when $r$ and $s$ lie on opposite sides of line $pq$ and $p$ and $q$ lie on opposite sides of line $rs$:

$$\text{orient}(p, q, r) \cdot \text{orient}(p, q, s) < 0 \quad \text{and} \quad \text{orient}(r, s, p) \cdot \text{orient}(r, s, q) < 0.$$

Four predicate evaluations, no division and no intersection point. Example: $p = (0, 0)$, $q = (6, 6)$, $r = (0, 6)$, $s = (6, 0)$: the first pair gives $+36$ and -36, the second -36 and $+36$, so they cross. When any value is zero an endpoint is collinear with the other segment and the engine must handle touching and overlap explicitly. Computing the actual crossing point needs a division and is deferred until the topology is known.

## Exactness and bit width

With integer coordinates the predicate can be evaluated exactly, but only if the arithmetic is wide enough. If coordinates are signed $b$-bit integers, a difference needs $b + 1$ bits, so its magnitude is at most $2^b - 1$. Each product is at most $(2^b - 1)^2 < 2^{2b}$, and the determinant, twice the area of a triangle inside the coordinate square, is bounded by the same value. A signed type of $2b + 1$ bits therefore always suffices, and $2b$ bits do not. For 32-bit coordinates that is 65 bits, one more than a 64-bit integer. Engines respond by limiting the coordinate range, using 128-bit intermediates, or computing in floating point with an error bound and falling back to exact arithmetic only when the result is near zero.

## Going deeper

Shewchuk, "Adaptive Precision Floating-Point Arithmetic and Fast Robust Geometric Predicates" (Discrete and Computational Geometry, 1997); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 1.
