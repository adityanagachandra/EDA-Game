# Design for assembly

Design for assembly (DFA) is the set of layout rules that make a board cheap and reliable to build: fiducials and handling edges for the machines, spacing and orientation for placement and inspection, thermal and via rules for good joints, and BOM and side choices that cut setups. These rules are placement constraints, so a board tool checks them during layout, not after. Each rule exists because of a specific machine or process step.

## Rules the machines need

Assembly equipment needs references and handling. **Fiducials** are bare copper marks, global on the board or panel and local near fine-pitch parts, that the placement camera finds to correct for board position, rotation and stretch. Conveyors and clamps grip the board edges, so two edges keep a component-free margin, or the board is panelized with breakaway rails. **Thermal reliefs** connect pads to planes through spokes so the joint reaches soldering temperature instead of losing heat into the plane. **Polarity and pin-1 marks** on the silkscreen and assembly drawing let programmers and inspectors orient diodes, polarized capacitors and ICs; a rotated polarized part passes every electrical check before assembly and fails afterwards.

## Spacing, panels and joints

Spacing rules come from inspection and rework. Tall parts shadow automated optical inspection cameras, which view at an angle, and block rework nozzles, so the required gap grows with component height. Small boards are **panelized**: several copies run through printing, placement and reflow together, raising throughput, and are separated afterwards by routing tabs or V-scoring, which adds keep-outs near the break lines. **Open vias in lands** wick molten solder down the barrel and starve the joint, so vias stay out of lands unless filled and capped. Placing polarized parts in a few consistent orientations reduces programming and inspection errors. Each of these is a geometric check on the placed board that a DRC engine can run with the right rule types.

## Cost and yield trade-offs

DFA is also about cost. Every distinct BOM line needs a feeder, a setup and a kitting step, so consolidating close resistor values where the circuit allows saves machine time. Double-sided SMT needs a second print and reflow pass and possibly adhesive for heavy bottom parts; it buys area, so it should be a deliberate choice. One fine-pitch BGA can raise the assembly quote sharply: finer stencils, tighter placement accuracy, via-in-pad boards and X-ray inspection. Because these are placement decisions, DFA review must happen before routing is complete; moving parts or changing sides after routing means rerouting and re-verifying signal integrity. A tool that scores DFA live during placement turns a late review into an early constraint.

## Going deeper

IPC-7351 and IPC-2221 sections on assembly considerations. Prasad, "Surface Mount Technology: Principles and Practice" (2nd ed., 1997), chapter on design for manufacturability.
