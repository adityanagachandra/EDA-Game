# Design for manufacturability and yield

Design for manufacturability (DFM) changes a legal layout to make it more robust: fewer shapes sensitive to random defects, fewer patterns that print poorly, more redundancy where a single failure would kill the chip. The changes are optional and cost area, capacitance or runtime, so DFM tools score candidates and apply fixes where they are cheap. Their inputs are critical area, recommended rules, hotspot patterns and via counts.

## Layout fixes that buy yield

The standard post-route DFM moves target random defects. Via doubling adds a second via beside each single one where space allows, so a location fails only if both fail. Wire spreading moves wires off minimum spacing into free tracks, lowering short critical area. Wire widening raises the width of long wires to lower open critical area. Recommended rules (longer line-end extensions, larger enclosures, wider spaces near line ends) reduce sensitivity to lithography and overlay variation. Each fix adds capacitance or area, so routers and DFM engines treat them as soft costs to apply where slack and space exist.

## Redundant vias by the numbers

If each via fails open independently with probability $p$ and a design has $N$ single vias, via-limited yield is $(1 - p)^N \approx e^{-Np}$. With $N = 5 \times 10^7$ and $p = 2 \times 10^{-8}$: $Np = 1$ and $Y = 36.8\%$. Doubling every via makes each location fail only with probability $p^2 = 4 \times 10^{-16}$, so $e^{-N p^2} \approx 1 - 2 \times 10^{-8}$, essentially 100%. Linearizing as $1 - Np$ would give zero here, which shows why the exponential form matters. Real redundant vias are not perfectly independent, since a defect large enough can take out both, but the improvement is still large.

## Systematic versus random yield

Random-defect yield follows defect statistics and critical area. Systematic yield loss comes from specific layout patterns that fail under process variation: litho hotspots, CMP-sensitive density regions, stress-sensitive structures. Total yield is often modelled as $Y = Y_{sys} \cdot Y_{rand}$, which assumes the mechanisms are independent. Early in a node's life systematic loss dominates, so flows use litho simulation and pattern matching against libraries of known-bad patterns to find and fix hotspots. As the process matures, random defects become the larger share and critical-area-driven fixes matter more. Separating the two in silicon data takes care, since both appear as failed dies.

## Going deeper

Gupta and Kahng, "Manufacturing-aware physical design", ICCAD (2003); Kahng, "Key directions and a roadmap for electrical design for manufacturability", ESSCIRC (2007).
