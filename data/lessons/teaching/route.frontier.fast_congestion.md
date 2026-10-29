# Fast congestion estimators in modeling loops

An agent exploring placement and flow options cannot afford a full global route for every candidate. It needs a fast congestion estimator in its loop and a rule for when to trust it and when to call the accurate router. RUDY and probabilistic estimators are the classic fast modes; knowing what they miss is what makes the loop safe.

## The RUDY estimate

RUDY (rectangular uniform wire density, Spindler and Johannes) spreads each net's expected wire uniformly over its bounding box. For a box $W$ by $H$ and wire pitch $p$, the half-perimeter wire occupies area $p(W + H)$, so the density inside the box is $p(W + H)/(WH)$, the fraction of the box's area taken by the net.

Example: pitch 0.2 um and a 31 by 35 um box give $0.2 \cdot 66 / 1085 = 0.012$. Small boxes are dense, long thin boxes spread thin. RUDY needs no routing and no decomposition, and it changes smoothly when cells move, so a placer can evaluate it, and even differentiate it, inside its own optimisation loop.

## Summing over a tile

The demand map adds the RUDY densities of every net whose box covers a tile. Example: two nets whose boxes contain a 1 by 1 um tile with densities 0.035 and 0.032 give the tile a density of 0.067, and nets whose boxes miss it contribute nothing. Comparing the map with the routing capacity per tile flags likely hot spots.

The estimate is cheap because it ignores everything a router does: it assumes wire spreads evenly over the box, ignores detours, blockages and layer limits, and treats multipin nets as their half-perimeter. It is good at ranking regions and poor at absolute overflow, which is why it serves as a fast mode, not as signoff.

## When to recurse to the accurate router

In a modeling loop the agent flags a tile for the accurate global router only when the fast estimate exceeds a threshold. The dangerous errors are misses: tiles the estimator called safe that the accurate router later finds overflowed. Counting them on calibration designs sets the threshold. Example: if several tiles with fast utilisation below 0.8 overflow under the accurate router, the threshold must drop or the estimator must improve.

The loop should also recurse when the estimate's variance is high, when a candidate is about to be committed, and periodically to recalibrate. This is the fast-versus-accurate pattern of agent-first engines: cheap models explore, accurate engines decide, and the agent watches the disagreement between them as closely as it watches either one.

## Going deeper

P. Spindler and F. M. Johannes, "Fast and Accurate Routing Demand Estimation for Efficient Routability-driven Placement" (2007).
