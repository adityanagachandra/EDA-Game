# Loop and partial inductance

Inductance belongs to a current loop, not to a wire, so it depends on where the return current flows. On a chip that return path is a mesh of power and ground wires that is rarely known in advance. Extractors solve this by computing partial inductances, self and mutual, for wire segments, and leaving the loop to be formed by the circuit simulation. Understanding how partial terms combine into loop inductance is the foundation of on-chip inductance extraction.

## Partial self and mutual terms

The partial inductance method assigns each segment a partial self inductance and each pair of segments a partial mutual inductance, defined as if each segment's current returned at infinity. These quantities depend only on geometry. A loop formed by a signal segment and its return combines them with signs set by current direction: for current out on the signal and back on the return, $L_{loop} = L_s + L_r - 2M$. With $L_s = 230$ pH, $L_r = 340$ pH and $M = 50$ pH, $L_{loop} = 470$ pH. The mutual term subtracts because opposing currents cancel each other's field between them. This decomposition lets an extractor emit a netlist of partial inductors and couplings without choosing the return path.

## Where the return current flows

At low frequency, return current divides by resistance. At high frequency, impedance is dominated by inductance and the current takes the path that minimizes loop inductance, which is the path closest to the signal, because a closer return has a larger mutual term. With a signal of 770 pH and two candidate returns of 710 pH each, an adjacent ground wire with $M = 550$ pH gives $770 + 710 - 1100 = 380$ pH, while a distant strap with $M = 120$ pH gives 1240 pH. Fast edges therefore concentrate return current in the nearest ground wire. This is why clock and high-speed nets are routed with dedicated shields or adjacent ground lines: they define the return path and keep loop inductance small and predictable.

## Non-additivity and dense matrices

Partial inductance does not add over length. A straight wire split into two collinear segments in series has $L = L_a + L_b + 2M_{ab}$, where $M_{ab}$ is the mutual between the halves. Segments of 620 and 230 pH with a 50 pH mutual give 950 pH, not 850. Mutual inductance between segments falls off slowly with distance, roughly logarithmically, unlike capacitive coupling, which is screened by neighbours. Every segment couples to every other, so the partial inductance matrix is dense. Extractors sparsify it with techniques such as shift-truncation, return-limited models or the inverse-inductance (susceptance) matrix, which is much more local, while preserving passivity of the resulting model.

## Going deeper

Ruehli, "Inductance calculations in a complex integrated circuit environment", IBM Journal of Research and Development (1972); Paul, "Inductance: Loop and Partial" (2010).
