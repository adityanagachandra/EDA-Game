# What a synthesized netlist preserves and abstracts

A synthesized netlist is a list of library cell instances and the nets that connect them. It is the hand-off from synthesis to placement, timing and verification, so a tool designer must know exactly what it guarantees. It preserves function at registers and ports, and it discards most of what the RTL author wrote: expressions, internal names, word structure and any notion of delay.

## What the netlist guarantees

Synthesis guarantees that the netlist is **functionally equivalent** to the RTL at the registers and primary outputs: each register input and each output computes the same Boolean function of the register outputs and primary inputs. Registers and ports survive with recognizable names, so an equivalence checker can match them. Nothing else is guaranteed. Internal combinational signals can be merged, duplicated or removed; multi-bit variables are bit-blasted into single-bit nets; and an X assigned in the RTL is a don't care that synthesis resolves to whatever value simplifies the logic, which is one source of RTL versus gate-level simulation mismatches.

## What the netlist adds and what it leaves out

Mapping binds each node to a library cell, and through the cell name the netlist inherits everything the library knows: area, pin capacitances and timing arcs from Liberty, and footprint and pins from LEF. The netlist itself holds no delays and no coordinates. A timer computes delays from the netlist plus Liberty tables plus wire parasitics, and a placer assigns coordinates. This separation is deliberate: the same netlist can be timed at several corners and placed many ways. It also means synthesis decisions are made against estimates of wires that do not exist yet, which is the root of the synthesis-to-placement correlation problem.

## Consequences for verification and debug

Because only register correspondence is preserved, **combinational equivalence checking** works only when registers match one to one. Retiming breaks that correspondence and needs sequential equivalence checking. Debug taps, clock-domain-crossing synchronizers and other structures that must survive are protected with keep or dont_touch attributes, since otherwise optimization owns them. Two netlists that both pass equivalence can still differ wildly in area, timing and power: equivalence says nothing about quality. Finally, a flow that resynthesizes after a small RTL change can produce a very different netlist, so late changes are made as ECO patches that keep the rest of the netlist intact.

## Going deeper

Kuehlmann and van Eijk, "Combinational and sequential equivalence checking", in Hassoun and Sasao (eds.), "Logic Synthesis and Verification" (2002); the Liberty and LEF reference manuals for what cell views contain.
