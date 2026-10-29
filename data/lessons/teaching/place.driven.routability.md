# Routability-driven placement

A placement with short wires can still be unroutable if too many wires must pass through one region. Routability-driven placement estimates routing demand during placement and spreads cells out of congested regions. RUDY, Spindler and Johannes's rectangular uniform wire density, is the cheapest widely used estimate: it needs only net bounding boxes, no routing.

## RUDY wire density

RUDY assumes each net's wire is spread uniformly over its bounding box. A net with box $w \times h$ has estimated wirelength $w + h$, its HPWL, so its wire density inside the box is
$$d = \frac{w + h}{wh},$$
in units of wire length per unit area. A net spanning 21 by 5 um has density $26/105 = 0.2476$ um per um$^2$. Small boxes have high density: a local net with pins close together concentrates its wire in a small area, which is exactly the kind of local crowding that creates routing hot spots. Degenerate boxes with zero width or height are clamped to a minimum size, typically one bin or one track pitch, so the density stays finite.

## Demand per bin

The routing demand in a bin is the sum over nets of their density times the area of their box inside the bin:
$$D_b = \sum_e d_e\, |B_e \cap b|.$$
For a 10 by 10 bin overlapped 70 um$^2$ by a net of density 0.2054 and not at all by another, the demand is 14.38 um of wire. Summed over all bins, each net contributes exactly its HPWL, so RUDY redistributes wirelength rather than inventing it. The computation is a rectangle-overlap scatter over nets and bins, linear in the total number of net-bin overlaps, and on GPUs it uses the same accumulation tricks as the density map.

## From demand to congestion and action

Capacity is the wire length the bin can carry: with $L$ usable layers at track pitch $p$, a $B \times B$ bin offers roughly $L B^2/p$ um of wire, less what power stripes and blockages take. The congestion ratio is demand over capacity, and bins above 1 are predicted to overflow. RUDY ignores detours, pin access and layer direction, so it is a fast proxy; RePlAce and OpenROAD can instead call a global router for a more accurate congestion map at milestones. Placers act on the map by inflating cells in congested bins, lowering the local target density, or adding the congestion as a penalty, so cells spread where wires crowd.

## Going deeper

Spindler and Johannes, "Fast and accurate routing demand estimation for efficient routability-driven placement" (DATE 2007); Cheng, Kahng, Kang and Wang, "RePlAce" (IEEE TCAD, 2019), on routability mode.
