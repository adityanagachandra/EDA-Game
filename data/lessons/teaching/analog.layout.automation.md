# Analog layout automation: MAGICAL and ALIGN

Analog layout is still largely hand-drawn, but research and open-source flows now generate layouts from netlists. MAGICAL and ALIGN are the best-known open efforts: both extract analog constraints, place devices under symmetry and matching constraints, route with analog-aware rules and output GDS that is checked by DRC, LVS and post-layout simulation. Their structure shows what an analog layout engine must contain.

## Constraint extraction

A netlist carries no layout intent, so the first stage recognizes structures and emits constraints. Differential pairs (shared source, complementary gates, equal sizes) become symmetry pairs about a common axis. Current mirrors become matched groups, often common-centroid. Sensitive nets such as inputs and high-impedance nodes are marked for short, shielded, symmetric routing. Recognition is pattern matching on the device-net graph, by subgraph rules, graph similarity or learned classifiers. Over-constraining is costly (every match costs area and flexibility) and under-constraining loses performance, so designer annotations are usually allowed to override or add constraints. The constraint set is the interface between circuit knowledge and geometric optimization.

## Placement under symmetry

Placement minimizes area and wirelength subject to non-overlap and the extracted constraints. Mirror symmetry about a vertical axis $x = a$ requires paired devices to satisfy $y_1 = y_2$ and $x_1 + x_2 = 2a$, and self-symmetric devices to be centred on the axis. Two families of placers exist. Topological placers search sequence-pair or B*-tree encodings with simulated annealing; every encoding decodes to a non-overlapping placement, and symmetric-feasible encodings guarantee symmetry. Analytical placers, as in MAGICAL, solve a smooth wirelength and density objective with symmetry as equality constraints or penalties, then legalize. ALIGN builds matched primitives (mirrors, pairs) with parameterized generators first, so the placer handles fewer, larger, already-matched blocks.

## Routing and the verification loop

Analog routing honours symmetric net pairs (mirror-image routes with equal lengths, layers and vias), shielding for sensitive nets, wide wires for current-carrying nets and limits on coupling. The output must pass DRC and LVS, but those prove only legality and connectivity. Acceptance requires extraction and post-layout simulation against the specs. When a spec fails, the flow must trace the failure to parasitics and feed it back as tighter constraints or re-sizing, which is the hardest part to automate. Published results compare generated layouts with manual ones on post-layout performance; for a tool designer the lesson is that an analog layout engine is a loop around constraints and simulation, not a one-pass placer and router.

## Going deeper

Xu et al., "MAGICAL: Toward fully automated analog IC layout leveraging human and machine intelligence" (ICCAD, 2019); Kunal et al., "ALIGN: Open-source analog layout automation from the ground up" (DAC, 2019).
