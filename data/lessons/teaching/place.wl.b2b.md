# The bound-to-bound net model

The bound-to-bound (B2B) net model, introduced with the Kraftwerk2 placer, replaces each net by two-pin edges chosen from the current placement so that the quadratic objective reproduces HPWL exactly at that placement. It combines the fast linear solve of quadratic placement with the accuracy of HPWL, and it underlies modern quadratic placers such as SimPL.

## Edges between the bounds

Sort a net's pins by x. The two extreme pins are the bounds. B2B connects the bounds to each other and connects every inner pin to both bounds; inner pins are not connected to each other. A $k$-pin net gets $1 + 2(k-2) = 2k - 3$ edges, linear in $k$ like the star model, with no extra variables. The y dimension is modeled separately, with its own bounds, so a pin can be a bound in x and an inner pin in y. For pins at x = 3, 8, 15 and 21 the bounds are 3 and 21, and the edges are 3-21, 8-3, 8-21, 15-3 and 15-21.

## Weights that reproduce HPWL

Each edge gets weight $w_{ij} = \frac{1}{(k-1)\,|x_i - x_j|}$, computed from the current positions. Then each edge contributes $w_{ij}(x_i - x_j)^2 = |x_i - x_j|/(k-1)$. The bound edge contributes $L/(k-1)$, where $L$ is the net's span, and every inner pin contributes its distances to both bounds, which add to $L$. The total is $(L + (k-2)L)/(k-1) = L$: the quadratic cost equals the x span exactly. Some papers write the weight as $2/((k-1)|x_i - x_j|)$, which only scales the objective. In the example above, $L = 18$ and the edge 8-3 has weight $1/(3 \cdot 5) = 0.0667$. The weights are large for short edges and small for long ones, which is exactly what converts a squared distance back into a linear one.

## Iterating the model

Once the weights are fixed, the quadratic system is solved and cells move, so the bounds and distances change and the model is no longer exact. B2B placers therefore iterate: build the model from the current placement, solve, rebuild, solve again, converging toward a placement that minimises HPWL subject to the spreading forces. Two practical details matter. Pins that coincide make $|x_i - x_j| = 0$, so implementations clamp the distance to a small minimum to avoid infinite weights. And because the edges change every iteration, the matrix structure changes too, so the solver reassembles it; conjugate gradient with a warm start from the previous solution keeps each solve cheap.

## Going deeper

Spindler, Schlichtmann and Johannes, "Kraftwerk2, a fast force-directed quadratic placement approach using an accurate net model" (IEEE TCAD, 2008); Kim, Lee and Markov, "SimPL: an effective placement algorithm" (ICCAD 2010).
