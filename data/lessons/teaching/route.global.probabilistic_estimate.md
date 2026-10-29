# Probabilistic congestion estimation

Before routing, a placer or an early flow step needs to know where congestion will be. Probabilistic estimators guess each net's route distribution without searching: a two-pin connection is assumed to take each of its L shapes with equal probability, and demand on an edge is the expected number of nets crossing it. The estimate is cheap, smooth and good enough to steer placement.

## Expected demand of one net

Take a two-pin connection from $a$ to $b$ that differs in both coordinates. Its horizontal-first L runs along $a$'s row from $a$'s column to $b$'s, then vertically; its vertical-first L runs vertically first, then along $b$'s row. With each L chosen with probability 1/2, a horizontal edge in $a$'s row between the two columns gets expected demand 1/2 from this net, and so does one in $b$'s row. A net in a single row is straight and puts demand 1 on every edge it crosses.

Example: a net from (6, 1) to (1, 0) crosses the horizontal edge between (5, 0) and (6, 0) only on the L that runs along row 0, so its expected demand there is 0.5. Edges outside the pins' bounding box get nothing.

## Summing many nets

By linearity of expectation, the expected demand on an edge is the sum of the per-net expected demands, whatever the dependence between nets. Example: three nets contributing 0.5, 0 and 0.5 to an edge give expected demand 1. Doing this for every edge and net gives an expected-demand map in time proportional to the total bounding-box perimeter, with no search at all.

Variants spread demand over Z shapes too, or uniformly over all monotone routes, which smears demand into the interior of each bounding box. RUDY goes further and spreads each net's wire area uniformly over its box. All share the property that matters for a placer: they are cheap and change smoothly when a cell moves, so they can be evaluated inside an optimisation loop.

## Expected overflow is not overflow of the expectation

Capacity checks need overflow, not demand, and overflow is a convex function of demand: $o = \max(0, d - c)$. By Jensen's inequality, $E[\max(0, D - c)] \ge \max(0, E[D] - c)$, often strictly. Treating the nets' choices as independent, the exact expected overflow comes from enumerating (or convolving) the distribution of demand.

Example: four nets each cross an edge with probability 1/2, so $E[D] = 2$, and capacity is 2. Plugging the mean gives zero overflow, but $D$ is binomial: $P(D = 3) = 4/16$ and $P(D = 4) = 1/16$, so the expected overflow is $0.25 + 2/16 = 0.375$. An estimator that thresholds the mean map underestimates trouble exactly at edges near capacity, which are the ones that matter.

## Going deeper

J. Lou, S. Thakur, S. Krishnamoorthy and H. S. Sheng, "Estimating Routing Congestion Using Probabilistic Analysis" (2002); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
