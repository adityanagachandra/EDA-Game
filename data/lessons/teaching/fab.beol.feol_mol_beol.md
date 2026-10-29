# Process integration: FEOL, MOL and BEOL

Process integration is the art of ordering hundreds of unit steps so that each one leaves the wafer in a state the next can use. Foundries and designers divide the result into three segments: the front end builds devices, the middle of line connects to them, and the back end wires them together. Each segment has its own materials, rules, models and failure modes, and each maps to different engines in a design flow.

## Three segments of the flow

The front end of line (FEOL) forms the transistors: shallow trench isolation, wells, fins or nanosheet stacks, the gate stack, spacers, and epitaxial source and drain. The middle of line (MOL) forms the local contacts: trench contacts to source and drain, the gate contact, and the first via level, often called V0, that reaches the first metal. The back end of line (BEOL) is the stack of copper or other metal wiring levels separated by low-k dielectric and joined by vias, ending in thick top metals and pads. Thermal budget enforces the order: FEOL steps reach about 1000 C, while once silicides and copper exist the flow must stay near 400 C.

## Why the middle of line became its own segment

At older nodes, contacts were a single small step between devices and metal 1. As gate pitch shrank, the space between gate and contact became a few nanometres, and contact resistance became a large share of a transistor's total resistance. The MOL gained its own materials (cobalt or tungsten fills, thin liners), self-aligned processes that use the gate cap as an etch stop, and its own rules. Contact over active gate (COAG) lets the gate contact land above the channel region rather than on isolation, removing a constraint on cell height. For tools, MOL parasitics are now first-order inputs to cell characterization and extraction.

## Consequences for design tools

The segments map onto the flow. FEOL and MOL shapes are mostly fixed inside library cells and drawn by layout engineers or generators; BEOL shapes are drawn by routers. Overlay error between a via and the line below is a large fraction of spacing at tight pitch, so self-aligned via processes use the line trench to confine the via, which in turn changes via rules. Changing an MOL dielectric thickness changes gate-to-contact capacitance and therefore every cell's timing. Moving power delivery to the wafer back removes the power network from the front-side BEOL. Knowing which segment a change touches tells a tool designer which models and rules must be regenerated.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 11 on process integration; the IEEE IRDS "More Moore" roadmap chapters on interconnect and device integration.
