# Total capacitance of a net

A driver's load is the sum of everything it must charge: wire capacitance to ground, coupling to every neighbour, and the input capacitance of every sink pin. Extractors report these components separately so that downstream tools can treat them differently. A timer without signal-integrity analysis grounds the couplings; a crosstalk analyser keeps them as couplings; nothing should drop them. Assembling the total, and knowing what share of it is coupling, is the first step of any delay or noise reasoning about a net.

## Summing the components

A net's capacitance is wire ground capacitance (area plus fringe) plus coupling to each neighbour plus pin capacitance: $C_{tot} = L(c_g + c_{c,1} + c_{c,2}) + \sum C_{pin}$, using per-length values for a uniform wire. A 40 um net with 112 aF/um to ground, 22 and 61 aF/um of coupling to its two neighbours, and two pins of 0.5 and 1.1 fF has $40 \times 195 = 7800$ aF of wire, 7.8 fF, plus 1.6 fF of pins: 9.4 fF. Counting coupling at face value, as if the neighbours were grounded, is the quiet-neighbour assumption, Miller factor 1. Pin capacitance comes from the Liberty library, not from extraction, and the timer adds it when it combines the two.

## The coupling share

The fraction of a net's capacitance that is coupling tells how sensitive it is to its neighbours. A 200 um net with 79 aF/um to ground, 64 and 48 aF/um of coupling, and 4.6 fF of pins has 22.4 fF of coupling out of 42.8 fF: 52%. At that share, an aggressor switching against the victim can almost double the effective coupling load, which moves delay by tens of percent. Extractors and routers track this ratio: a high coupling share flags a net for spacing, shielding or a crosstalk-aware timing check. Long parallel runs at minimum spacing on the same layer, typical of buses, push the share highest.

## Grounded versus dropped coupling

Without SI analysis, a timer grounds each coupling capacitor: it keeps the capacitance but forgets which net is on the other side. That preserves the total load. A flow that instead drops coupling, keeping only ground and pin capacitance, under-reports load by the full coupling: for a 260 um net with 96 aF/um of coupling, that is 24.96 fF out of 40.82 fF, a huge optimistic error. Real flows lose coupling by accident through format conversions, filtering thresholds that discard small couplings without regrounding them, or incomplete SPEF reads. A robust extractor offers decoupling with regrounding as an explicit mode, and a robust timer checks that each net's total matches its header.

## Going deeper

The Liberty and SPEF chapters of the OpenSTA documentation; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure.
