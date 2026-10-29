# Epitaxy and strained layers

Epitaxy grows a crystalline film that continues the lattice of the crystal beneath it. Modern transistors depend on it: raised source and drain regions, silicon-germanium channels and stressors, and the alternating Si and SiGe stack from which gate-all-around nanosheets are released. A tool designer meets epitaxy through strain-dependent device models, merged source-drain regions in extraction, and spacing rules around fins and sheets.

## Single-crystal growth on a template

In epitaxy, precursor gases such as silane or dichlorosilane, plus germane or dopant gases, decompose at the hot surface and the atoms settle into lattice sites continuing the substrate crystal. The result is single-crystal rather than polycrystalline or amorphous. Selective epitaxy is tuned so growth nucleates only on exposed silicon and not on oxide or nitride, which lets a flow grow source and drain only in the cavities beside the gates. In situ doping during growth reaches high active concentrations without implant damage. Gate-all-around flows start with an epitaxial superlattice of alternating Si and SiGe layers; the SiGe is later etched away selectively to leave stacked silicon sheets.

## Lattice mismatch and strain

Germanium's lattice constant is 5.658 angstrom against 5.431 for silicon. By Vegard's law, $\mathrm{Si}_{1-x}\mathrm{Ge}_x$ has $a \approx 5.431 + 0.227x$ angstrom. The mismatch with a silicon substrate is $f = (a_{film} - a_{sub})/a_{sub}$. For $x = 0.3$: $a = 5.499$ angstrom and $f = 0.068/5.431 = 1.25\%$. A thin film that stays pseudomorphic adopts the substrate's in-plane spacing, so a SiGe film on silicon is under compression. Strain is a design tool: compressive stress in the channel raises hole mobility, so PMOS source-drain regions use SiGe, while NMOS benefits from tensile stress.

## Critical thickness and layout consequences

Strain energy grows with film thickness. Past the critical thickness it becomes cheaper for the film to relax by forming misfit dislocations, which destroy the benefit and can short junctions; higher germanium content lowers that limit. Strain also depends on geometry: a stressor relaxes near free edges, so mobility depends on active length and on distance to diffusion breaks, effects that PDK models capture as layout-dependent parameters. Laterally overgrown epitaxy on adjacent fins can merge, changing source-drain resistance and parasitic capacitance, so extraction and spacing rules must know whether neighbouring fins merge.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter on epitaxy; Matthews and Blakeslee on misfit dislocations and critical thickness.
