# Merging segments in the Manhattan plane

In a real layout, a merge point is not on a fixed wire between two roots. It can be anywhere that has the right rectilinear distance to both children. The set of such points is the merging segment, a diagonal line in the Manhattan plane. Merging segments are the data structure that lets DME defer the choice of tap and so minimise wire; computing them, and distances to them, is the geometric core of every zero-skew router.

## Edge lengths from the zero-skew condition

Take sinks a and b with downstream delays $t_a$ and $t_b$ (path-length model) and Manhattan distance $d = |x_a - x_b| + |y_a - y_b|$. With no detour, the edges to a and b have lengths $e_a + e_b = d$ and zero skew needs $e_a + t_a = e_b + t_b$, so
$$e_a = \frac{d + t_b - t_a}{2}.$$
For $a = (40, 30)$ with $t_a = 0$ and $b = (20, 0)$ with $t_b = 10$, $d = 50$ and $e_a = 30$ µm. This is the linear merge, but in the plane: the edge lengths are fixed, the location is not.

## The shape of the segment

A point $p$ with $|p - a| = e_a$ and $|p - b| = e_b$, where $e_a + e_b = d$, must lie in the bounding box of a and b, since only there does the Manhattan triangle inequality hold with equality. Inside the box, $|p - a| = e_a$ is a line of slope $\pm 1$, so the merging segment is a diagonal (a Manhattan arc), clipped by the box. For $a = (50, 30)$, $t_a = 15$, $b = (20, 20)$, $t_b = 20$: $d = 40$ and $e_a = 22.5$. The segment runs from $(37.5, 20)$ to $(27.5, 30)$, an x extent of 10 µm. It shrinks to a point when the box is a line or when $e_a$ or $e_b$ is zero.

## Distance from a point to a segment

In the top-down pass, a placed parent at $p$ must be connected to the nearest point of the child's segment. Along a straight segment, Manhattan distance to $p$ is a convex piecewise-linear function, so its minimum is at an endpoint or at a point level with $p$ in x or in y. For the segment from $(17.5, 20)$ to $(27.5, 10)$ and $p = (0, 30)$, the endpoints are at 27.5 and 47.5, and no interior point is level with $p$, so the distance is 27.5 µm. When $p$ lies diagonally off the segment, every point of it can be equidistant, and any choice is optimal.

## Going deeper

Chao, Hsu, Ho, Boese and Kahng, "Zero skew clock routing with minimum wirelength", IEEE Transactions on Circuits and Systems II (1992); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
