# Etching: anisotropy, selectivity and over-etch

Etching transfers a lithographic pattern into the film below by removing exposed material. Its directionality, its selectivity to the mask and underlying layers, and the over-etch needed to clear thickness variation decide how faithfully the drawn shape becomes a physical one. Many design rules, from minimum spacing to enclosure and end-of-line margins, are etch rules in disguise, and etch bias is one of the effects OPC must model.

## Isotropic and anisotropic etch

Wet chemical etches remove material at the same rate in every direction: they are isotropic and undercut the mask by about the etch depth. Plasma, or reactive ion etching, combines reactive radicals with ions accelerated toward the wafer, and sidewall passivation protects vertical faces, so material is removed mainly downward. Anisotropy is $A = 1 - r_{lat}/r_{vert}$. A wet etch has $A = 0$; a plasma etch removing 300 nm/min vertically and 9 nm/min laterally has $A = 1 - 9/300 = 0.97$. Nanometre features with vertical sidewalls need $A$ very close to 1, which is why critical layers use plasma etch.

## Over-etch and etch time

A film is never perfectly uniform across a wafer, and etch rate varies too. To be sure every region clears, the recipe runs past the nominal endpoint by an over-etch, often 10 to 50 percent. The total time is $t = (h/r)(1 + OE)$. A 200 nm film at 250 nm/min takes 48 s nominally; with 30% over-etch it runs 62.4 s. During over-etch, the regions that cleared first are attacking the layer below, so the next number that matters is selectivity. Endpoint detection, which watches plasma emission for the moment the film clears, lets recipes keep over-etch as short as the uniformity allows.

## Selectivity and mask budget

Selectivity is the ratio of etch rates of the target film to another material, such as the mask or the stop layer below. A mask must survive the whole etch: consumed mask thickness is the equivalent film etched divided by selectivity. Etching 200 nm with 30% over-etch through a mask with 10:1 selectivity consumes $260/10 = 26$ nm of mask. Low selectivity forces thicker resist or hard masks, and thick resist limits resolution and aspect ratio. Etch also has pattern dependence: narrow openings etch more slowly than wide ones (aspect-ratio-dependent etching, or RIE lag), and dense and isolated features etch differently, effects that feed OPC models and density rules.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 10; Lieberman and Lichtenberg, "Principles of Plasma Discharges and Materials Processing".
