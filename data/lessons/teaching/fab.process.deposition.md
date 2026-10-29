# CVD, PVD and ALD

Most layers in a chip are deposited films: dielectrics between wires, barriers and seeds for copper, spacers that set feature widths, high-k gate dielectrics and work-function metals. Three families dominate, and each has a different conformality, thickness control and throughput. A tool designer cares because these properties set minimum widths, spacings and aspect ratios in the rules, and because in self-aligned patterning a film's thickness becomes a printed dimension.

## Three ways to put down a film

Chemical vapour deposition (CVD) reacts gas-phase precursors at or near a heated surface; plasma-enhanced CVD lowers the temperature so it can run in the back end. It is fast and fairly conformal. Physical vapour deposition (PVD), mainly sputtering, knocks atoms off a solid target with argon ions; they fly in nearly straight lines, so the flux is directional. It deposits metals such as barrier and seed layers. Atomic layer deposition (ALD) alternates two precursors, each of whose surface reactions saturates after one monolayer, so each cycle adds a fixed, small thickness regardless of local flux.

## ALD cycles and step coverage

Because ALD is self-limiting, thickness is cycles times growth per cycle (GPC). A 2.0 nm high-k layer at 0.1 nm per cycle needs 20 cycles; at 0.12 nm per cycle it needs $\lceil 16.7 \rceil = 17$. Conformality is measured by step coverage: film thickness on a trench bottom or sidewall divided by thickness on the open field. ALD reaches nearly 100%. CVD varies with chemistry and pressure. PVD in a high-aspect-ratio via can fall to a small fraction, because the bottom sees only a narrow cone of the incoming flux. If a PVD film is 20 nm on the field and 7 nm on the bottom, step coverage is 35%.

## Choosing a method by what the layer must do

A gate dielectric about two nanometres thick wrapped around a fin needs sub-angstrom control and full coverage, so it is ALD. Spacers in self-aligned double patterning become the final line width, so they need conformal, uniform ALD or CVD. A micrometre of inter-layer oxide would take thousands of ALD cycles, so it uses plasma CVD. PVD barriers overhang at the top of narrow openings and can pinch off before fill, leaving voids. These limits travel into design as rules on minimum trench width, spacing and aspect ratio. A fill process with poor conformality is one reason a foundry forbids certain narrow gaps.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 9 on thin-film deposition; George, "Atomic Layer Deposition: An Overview", Chemical Reviews (2010).
