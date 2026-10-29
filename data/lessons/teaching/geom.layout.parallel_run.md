# Overlap area and parallel run length

Parasitic extraction turns geometry into capacitance and resistance. Its geometric inputs are simple: the overlap area between shapes on different layers (plate capacitance), the parallel run length and spacing between neighbouring wires on one layer (coupling capacitance), and the widths and lengths of wires (resistance). Rule-based extractors compute these quantities with Booleans and sweeps, then look up capacitance per unit from precharacterized tables.

## Parallel run length

Two parallel horizontal wires couple along the length where their x-projections overlap, the parallel run length (PRL):

$$\text{PRL} = \max(0, \min(a_2, b_2) - \max(a_1, b_1)).$$

Example: wires over $x \in [40, 300]$ and $x \in [120, 450]$ have a PRL of $300 - 120 = 180$ nm. The coupling capacitance between them is roughly proportional to PRL divided by a function of their spacing, so an extractor needs both numbers for every neighbouring pair. Wires that do not overlap in projection have PRL 0 and couple only weakly through fringe fields at their ends. PRL also enters design rules: the spacing table at advanced nodes requires more space for long parallel runs, so the same quantity serves both DRC and extraction.

## Overlap between layers

Where a wire on metal 1 crosses a wire on metal 2, the overlap area forms a parallel-plate capacitor through the dielectric between them. The total overlap area between two layers is the area of their AND, which for rectangular wires is the sum of pairwise intersection areas (when shapes on each layer are already merged and disjoint). Example: M1 wires $[0, 50] \times [0, 4]$ and $[5, 60] \times [20, 26]$ crossing an M2 wire $[10, 16] \times [0, 40]$ overlap it in $6 \times 4 = 24$ and $6 \times 6 = 36$, total 60. Plate capacitance is that area times a per-unit-area coefficient for the layer pair, plus fringe terms per unit perimeter of the overlap. A sweep or spatial index finds the crossing pairs without testing every combination.

## Neighbours within a coupling window

Coupling falls off quickly with distance, so extractors consider only neighbours within a window $S$ of the victim wire, and often only the nearest neighbour on each side along each stretch (the others are shielded). For each neighbour with edge-to-edge gap at most $S$, the extractor adds its PRL with the victim, usually binned by spacing so that each bin uses its own capacitance per unit length. Example: a victim over $x \in [0, 200]$ with neighbours at gaps 15, 40 and 22 nm and PRLs 120, 90 and 60 nm, with $S = 30$, counts $120 + 60 = 180$ nm. The search is a window query: grow the victim by $S$ and ask the index for shapes on the same layer.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure and parasitics; the OpenROAD documentation for OpenRCX.
