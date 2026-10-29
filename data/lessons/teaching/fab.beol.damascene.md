# Damascene copper and barriers

Copper cannot be plasma-etched cleanly, so copper wires are made by damascene: etch trenches and via holes in dielectric, line them with a barrier, fill with copper, and polish the excess away. The barrier keeps copper from diffusing into the dielectric and silicon, but it is poorly conducting and does not scale down with the wire, so it takes an ever larger share of the cross-section. That share is a direct input to wire resistance in extraction.

## The dual-damascene sequence

In dual damascene, one level's via holes and line trenches are both etched into the dielectric, in a via-first or trench-first order. A thin barrier, typically tantalum nitride with a tantalum or cobalt liner, is deposited by PVD or ALD, then a copper seed by PVD. Electroplating fills trench and via together from the bottom up. CMP removes the copper and barrier overburden from the field and leaves copper only in the trenches, flush with the dielectric. A dielectric cap then seals the top. Because filling and polishing define the line, CMP dishing and erosion directly change line thickness, which is why metal density rules exist.

## Copper area after the barrier

The barrier lines both sidewalls and the bottom. For a trench of width $w$ and height $h$ with barrier thickness $b$, copper occupies $(w - 2b)(h - b)$. A 24 nm by 48 nm trench with 3 nm barrier holds $18 \times 45 = 810$ nm^2 of copper out of 1152, which is 70.3%. With copper resistivity 5 micro-ohm cm (50 ohm nm, including size effects) and negligible barrier current, resistance per micrometre is $50 \times 1000/810 = 61.7$ ohm. Note that 1 micro-ohm cm equals 10 ohm nm, a conversion worth memorizing.

## A barrier that does not scale

At the same copper resistivity, the barrier raises resistance by $A_{trench}/A_{Cu}$: $1152/810 = 1.42$ in the example. Halve both dimensions to 12 nm by 24 nm with the same 3 nm barrier, and copper is $6 \times 21 = 126$ nm^2 of 288, a factor of 2.29. Copper's resistivity also rises sharply in narrow lines from surface and grain-boundary scattering. Together these effects explain why foundries study thinner barriers, barrier-free metals such as ruthenium or molybdenum, and alternative schemes such as subtractive etch for the tightest layers, and why extraction tools use width-dependent resistivity tables rather than a single bulk value.

## Going deeper

Andricacos et al., "Damascene copper electroplating for chip interconnections", IBM Journal of Research and Development (1998); Gall, "The search for the most conductive metal for narrow interconnect lines", Journal of Applied Physics (2020).
