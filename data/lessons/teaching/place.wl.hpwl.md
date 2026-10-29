# Half-perimeter wirelength

Half-perimeter wirelength (HPWL) is the standard estimate of a net's routed length: half the perimeter of the smallest axis-aligned box containing all its pins. Placers minimise it, compare results with it and report it, because it is exact for two- and three-pin nets, a lower bound for all others, and cheap enough to evaluate millions of times per second.

## The bounding box of a net

A net is a hyperedge connecting a set of pins. With pin coordinates $(x_i, y_i)$, its HPWL is
$$\mathrm{HPWL} = \left(\max_i x_i - \min_i x_i\right) + \left(\max_i y_i - \min_i y_i\right).$$
For pins at (2, 7), (10, 3) and (5, 12) the box spans x from 2 to 10 and y from 3 to 12, so HPWL = 8 + 9 = 17. Only the extreme pins matter: moving an interior pin anywhere inside the box leaves HPWL unchanged. The x and y parts are independent, which is why placers often optimise the two dimensions separately. The computation is a single pass over the pins, keeping a running minimum and maximum in each dimension, so it costs time linear in the pin count of the net.

## Total and weighted wirelength

The placement objective is the sum over nets, often weighted: $W = \sum_e w_e \,\mathrm{HPWL}_e$. Weights express priorities: timing-critical nets get larger weights, clock and high-fanout nets may be excluded or down-weighted because they are built later by CTS or buffering. Computing the total is linear in the number of pins. For three nets with HPWL 20, 35 and 12 and weights 1, 3 and 2, the objective is $20 + 105 + 24 = 149$. Because HPWL is computed per net, the total parallelises trivially, which GPU placers exploit.

## How good an estimate it is

For two- and three-pin nets, HPWL equals the length of the rectilinear Steiner minimum tree, the shortest possible rectilinear routing. For more pins it is a lower bound, and the Steiner tree can be noticeably longer for nets with many pins spread across the box; empirical correction factors that grow with pin count are sometimes applied. HPWL also ignores detours caused by congestion and blockages. Its other weakness is mathematical: max and min are not differentiable when two pins tie for the extreme, so gradient-based placers replace it by smooth approximations, while detailed placers use it directly because they evaluate discrete moves.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.2 on wirelength estimation.
