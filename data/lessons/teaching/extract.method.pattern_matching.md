# Rule-based extraction by pattern matching

Full-chip extraction cannot solve a field problem for billions of wire segments. Production extractors instead fracture the layout into pieces, classify each piece's local cross-section as a known pattern (a wire with a neighbour at some spacing, over some layer, under another), and look up its capacitance per length in tables precomputed by a field solver. This rule-based, or pattern-matching, approach is what makes signoff extraction feasible, and its accuracy is exactly the accuracy of its tables and its interpolation.

## Table lookup along one index

A pattern table stores capacitance per length for a pattern class, indexed by geometric parameters such as wire width and spacing to the neighbour. When the query lies between table points, the extractor interpolates. Along a single index this is linear: at width 0.02 um, if coupling is 75 aF/um at spacing 0.03 um and 55 aF/um at 0.04 um, then at 0.0325 um the fraction is $t = 0.25$ and the value is $75 + 0.25 \times (55 - 75) = 70$ aF/um. The lookup is the same operation a timer performs on an NLDM table, and it carries the same caveat: linear interpolation is only as good as the table's spacing relative to the curvature of the real function, which for coupling versus spacing is strongly convex.

## Two indices and a run length

With two indices, width and spacing, the extractor interpolates bilinearly over the four surrounding table points, then multiplies by the segment's parallel run length. For a table with 97 and 82 aF/um at width 0.06 um (spacings 0.06 and 0.08 um), and 102 and 87 aF/um at width 0.08 um, the point (0.075, 0.075) has fractions 0.75 and 0.75, giving 89.5 aF/um. A 80 um run then contributes 7.16 fF of coupling. Real pattern libraries index more parameters: the layers above and below, whether they are dense or sparse, the neighbour on the other side. Each added index multiplies table size, so library design is a trade between accuracy and the cost of characterization and lookup.

## Fracturing a varying net

Real nets change environment along their length: neighbours jog away, layers above end, spacing steps. The extractor fractures a wire into segments over which the cross-section is constant, looks each up separately and sums. A 70 um segment at spacing 0.05 um (103 aF/um) followed by 130 um at 0.06 um (86 aF/um) gives $7.21 + 11.18 = 18.39$ fF of coupling. Fracturing is where 3D effects slip through: line ends, corners and vias do not match any 2D cross-section, so libraries add special 3D patterns or corrections for them. A shape that matches no pattern at all, for example a new via-rail structure, is extracted with the nearest available pattern and an unknown error, which is a strong argument for calibration checks on new layout styles.

## Going deeper

Kao, Lo, Basel and Singh, "Parasitic extraction: current state of the art and future trends", Proceedings of the IEEE (2001); the documentation of the open-source OpenRCX extractor in OpenROAD.
