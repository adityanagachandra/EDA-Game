# Layout parasitics and post-layout performance

A schematic that meets every spec can fail once it is laid out. Wires add capacitance and resistance, junctions add capacitance, and adjacent routes couple. Extraction turns these into a post-layout netlist, and re-simulation shows the damage. For a tool designer, parasitics are where layout and sizing meet: a router needs constraints derived from circuit sensitivity, and a sizer needs margins for what layout will add.

## Parasitics move poles

Every node's pole is set by its resistance and capacitance, so extra capacitance moves poles down. In a two-stage op-amp the non-dominant pole is $g_{m6}/(2\pi C_{out})$. If the schematic load is 2 pF and the output route adds 0.5 pF, the pole falls by a factor $2/2.5$: from 159 MHz to 127 MHz for $g_{m6} = 2$ mS. Internal high-impedance nodes are even more sensitive, because their schematic capacitance is small, so a few femtofarads of wiring is a large relative change. Extraction tools (Magic in the open flow, and commercial RC extractors) report these capacitances per net, and the designer or a tool compares pre- and post-layout pole positions to find the nets that matter.

## Phase margin after layout

Lower non-dominant poles erode phase margin. With a dominant pole contributing -90 degrees, $\mathrm{PM} = 90 - \arctan(\mathrm{GBW}/f_{p2})$. Take GBW = 50 MHz and $f_{p2} = 159$ MHz in the schematic: PM = 72.5 degrees. After layout, $f_{p2} = 127$ MHz gives 68.5 degrees, a loss of 4 degrees from half a picofarad. Parasitic resistance in series with compensation capacitors or in supply lines adds further shifts, and coupling capacitance between input and output routes can create feedback paths of its own. This is why analog signoff is a post-layout simulation across corners, not a DRC and LVS check, and why phase margin specs carry a margin for layout.

## From sensitivity to routing constraints

The useful output of post-layout analysis is a constraint the router can obey. Suppose the margin must stay at least 60 degrees with GBW = 50 MHz. Then $f_{p2} \ge 50/\tan 30 = 86.6$ MHz, so $C_{out} \le g_{m6}/(2\pi \cdot 86.6\ \mathrm{MHz})$, which is 3.68 pF for $g_{m6} = 2$ mS. With 2 pF of device load, 1.68 pF is left for wiring; at 0.2 fF per micrometre that is a maximum route length of 8400 micrometres. A constraint-driven flow derives such limits from circuit sensitivities, hands them to placement and routing as net-length or capacitance bounds, and re-verifies after extraction. That loop of sensitivity, constraint, layout and verification is the backbone of parasitic-aware analog layout automation.

## Going deeper

Hastings, "The Art of Analog Layout" (2nd ed., 2005); the Magic VLSI layout tool documentation on extraction.
