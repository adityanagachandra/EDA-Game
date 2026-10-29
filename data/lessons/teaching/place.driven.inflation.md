# Cell inflation for congestion

Cell inflation is the simplest way to make a density-driven placer care about congestion: temporarily enlarge the cells that sit in congested regions, so the density force pushes them apart and leaves more room for wires. RePlAce and OpenROAD's gpl use it in routability mode, with a congestion map from RUDY or a global router.

## The inflation rule

For each cell, look up the congestion ratio $r$ (demand over capacity) of its bin. If $r \le 1$, leave the cell alone. Otherwise multiply its area by an inflation factor that grows with $r$ and is capped, for example
$$f = \min(f_{max}, r^k).$$
With $f_{max} = 2.5$ and $k = 2$, a 5 um$^2$ cell in a bin with $r = 1.2$ becomes $5 \cdot 1.44 = 7.2$ um$^2$. Inflation is virtual: the cell's real size does not change, only the area the density model sees. The exponent and cap are tunable parameters; the inflated areas are recomputed whenever the congestion map is refreshed.

## Total area after inflation

Inflation adds area to the design, which the density model must absorb. For cells of 10, 9 and 4 um$^2$ in bins with ratios 0.77, 1.37 and 0.78 and $f = \min(2, r)$, only the second cell grows, to 12.33 um$^2$, so the total rises from 23 to 26.33 um$^2$. Dividing the inflated total by the region's capacity gives the new utilization; if it exceeds the target density, the placer cannot satisfy both, so implementations cap the total inflation, for example to a fraction of the available whitespace, and scale individual factors down when the budget is exceeded. That budget is what keeps inflation from simply blowing up the design.

## When inflation helps and when it does not

Inflation works on local hot spots caused by too many cells in a small area: spreading the cells spreads their wires. It works poorly on congestion caused by long nets passing through a region, since moving local cells does not reroute those through-nets, and on congestion near macros, where pin access, not cell density, is the limit. It also lengthens wires, so inflation is usually applied after the placement has spread and is reverted, or kept small, for timing-critical cells. A routability-driven placer therefore watches both the congestion map and the wirelength and timing changes, and an agent tuning it treats the cap, exponent and refresh schedule as Tier 1 parameters, validated by a real global route.

## Going deeper

Cheng, Kahng, Kang and Wang, "RePlAce" (IEEE TCAD, 2019); the OpenROAD documentation for gpl's routability-driven mode.
