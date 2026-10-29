# Symmetry and ambiguity in LVS matching

Refinement cannot separate vertices that are truly symmetric: two identical inverters with private inputs and outputs look the same from every angle. LVS tools then break ties by choice and refine again. Understanding when refinement stalls, why a choice is safe, and when it is not explains both the speed of LVS and its rare pathological cases.

## Classes refinement cannot split

If an automorphism of the graph maps vertex $u$ to $v$, refinement always gives them the same colour, because every signature is preserved. Three identical NAND2 gates with private inputs and outputs, with only VDD and GND named, stabilize with 10 classes: VDD, GND, and eight classes of size three (four device roles and four net roles, one vertex per gate in each). Eight classes still contain more than one vertex. Such ties are not errors: equal-size classes in both netlists mean the matcher needs more information or an arbitrary choice.

## Breaking ties

When refinement stalls, the matcher picks one layout vertex and one schematic vertex from the same ambiguous class, gives them a fresh unique colour, and refines again. The new colour propagates and usually resolves everything connected to it. For $k$ identical independent gates, $k - 1$ choices suffice: after pairing all but one, the last is forced. If the tied vertices are truly symmetric, any choice works, because an automorphism maps one choice onto another. If they are only indistinguishable to refinement but not symmetric, as in some highly regular graphs, a wrong choice leads to a contradiction and must be undone.

## Names, permutations and regular arrays

Port and label names are the cheapest symmetry breakers: each pre-matched net gets a unique colour, and the inverters attached to it separate at once. Names can also create mismatches: if NAND2 inputs A and B are named and the layout swaps which one drives the bottom transistor of the stack, a transistor-level comparison reports a mismatch although the logic is the same, unless the tool is told those inputs may be permuted. Memory arrays are the hardest case, with millions of truly symmetric bit cells, which is why LVS flows name word and bit lines or compare the bit cell hierarchically.
