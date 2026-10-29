# Fill side effects and fill-aware flows

Fill is metal, so it changes capacitance, and it lives in the same layout as signal wiring, so it must be removed and regenerated when routing changes. Fill-aware flows control where fill goes, how it is tied, and when it is inserted relative to extraction and timing signoff. Ignoring these effects is a common way to fail timing after a clean DRC.

## Coupling added by fill

A first-order sidewall model makes each side's coupling capacitance proportional to the run length and thickness over the spacing, $C \propto t L / s$. A wire with a neighbour at 170 nm on one side and nothing on the other gains, when fill lands 40 nm away on the empty side, a new term proportional to $1/40$, which is $170 / 40 = 425$ percent of the original sidewall term. Real extraction is three-dimensional and fill is often floating, so the actual increase is smaller, but the model shows why fill near critical nets must be controlled.

## Timing-aware and grounded fill

Timing-aware fill keeps a larger spacing from critical nets, or keeps fill out of their neighbourhood, and makes up density elsewhere in the window. Fill can be floating or tied to a supply. Grounded fill adds a definite capacitance to ground; floating fill couples the nets around it through a conductor of unknown potential, which extraction tools must model specially. Keep-out regions protect analog and RF nets. After fill, signoff timing must use parasitics extracted with the fill present.

## Where fill lives in the flow

Fill is inserted after detailed routing, when the gaps are known, and before final extraction, DRC and timing signoff. Front-end layers (poly, diffusion) get dummy fill too, for polishing and etch uniformity. Wide metal that exceeds the maximum density is slotted. When an ECO changes routing, fill near the change is stripped and regenerated, then re-extracted. Open-source flows insert fill with OpenROAD's density fill or KLayout scripts, and commercial flows offer timing-aware fill modes.
