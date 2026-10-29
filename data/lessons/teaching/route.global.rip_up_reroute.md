# Rip-up and reroute

Rip-up and reroute is the outer loop of every practical router: remove the nets that cause overflow and route them again on updated costs. Which nets to rip, in what order to reroute them, and how many iterations to run are policy choices that strongly affect quality and runtime, and they are the natural place to customise a router.

## Which nets to rip up

After an initial routing, compute demand on every edge and find the overflowed edges. The standard choice rips up every net that uses at least one overflowed edge; a gentler one rips up only enough of them to clear each edge, preferring nets with cheap alternatives. Nets on uncongested edges are left in place, which keeps each iteration fast and preserves good routes.

Example: with eight edges and five nets, if edges $e_4$ and $e_8$ are overflowed, every net whose route includes $e_4$ or $e_8$ is ripped up, here perhaps three of the five. Ripping more nets explores more but wastes effort on nets that were fine; ripping fewer converges slowly.

## Reroute ordering

The ripped nets are rerouted one at a time, so their order matters: earlier nets see less congestion. Common rules sort by bounding-box size. Smallest first lets short nets with few alternatives claim the resources they need before long nets, which can detour, take them. Largest first settles the long nets that block many others while space remains. Other rules use the amount of overflow a net touches, its criticality or its pin count.

Example: ripped nets with half-perimeters 27, 15 and 13 rerouted largest first go in the order 27, 15, 13. Negotiation reduces the sensitivity to order, because prices adjust across iterations, but order still affects runtime and the final wirelength.

## Measuring a reroute

Each reroute removes the net's demand from its old edges and adds it to the new ones. The effect on total overflow is the sum of changes on every edge touched. Example: a net moved from $e_1, e_3$ onto $e_4, e_5, e_8$ may relieve $e_1$ but push $e_4$ and $e_8$ over capacity, raising total overflow from 1 to 2. A good reroute uses edge costs that rise steeply near capacity, so the search avoids exactly that.

Loops stop when overflow reaches zero, stops improving, or an iteration budget runs out. Production global routers combine rip-up and reroute with negotiation (history and present costs), escalating from pattern routing to monotonic to bounded maze routing for the hardest nets, and they report the residual overflow map when they cannot clear it.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the NTHU-Route 2.0 and NCTU-GR 2.0 papers.
