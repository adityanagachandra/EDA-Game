# KLayout, Magic and Netgen

After routing, the open flow turns to layout tools: KLayout to view, check and merge GDS, Magic to edit layout, check rules and extract circuits, and Netgen to compare netlists for LVS. These are the open counterparts of commercial physical verification, and a tool builder needs to know what each consumes, what it proves, and how to read a disagreement between them.

## KLayout for viewing, DRC and GDS merging

KLayout is a fast hierarchical layout viewer and editor for GDS and OASIS. GDS stores layers as numbers, so a layer-properties file maps each layer and datatype pair to a name and color; without it a SKY130 layout shows as anonymous numbers. KLayout's DRC engine runs rule decks written as scripts that express layer derivations and checks, such as width, spacing and enclosure, and it also supports LVS decks. In ORFS the finish stage uses KLayout to merge: the routed design references cells only as LEF abstracts, and KLayout replaces each reference with the cell's real GDS layout to write the final stream. Its hierarchical handling is why it opens very large layouts quickly: a cell drawn once and placed ten thousand times is stored once.

## Magic for layout, DRC and extraction

Magic is a VLSI layout editor with built-in, incremental DRC and circuit extraction, and it is widely used for SKY130 work. Its extraction turns layout into a SPICE netlist of devices and, optionally, parasitic resistance and capacitance, which feeds both LVS and post-layout circuit simulation. Many open flows run DRC in both Magic and KLayout, because the decks are written independently and one may encode a rule that the other misses or interprets differently. When they disagree, for example Magic flags a met2 spacing that KLayout passes, do not take the clean result as the answer: inspect the geometry at the marker, read the rule in both decks, and decide which encoding matches the process rule.

## Netgen for LVS

Netgen performs LVS by comparing two netlists: the one extracted from layout (by Magic or KLayout) and the source netlist from the design. It matches devices and connectivity as a graph comparison, iteratively partitioning devices and nets until each class matches one to one or a mismatch is isolated. Reading its report is a skill. Two source nets that appear as one layout net mean a short. One source net split into two layout nets means an open. A mismatch confined to power pins often means naming differences or implicit power connections in one netlist. Device-property mismatches, such as transistor width, point to layout that differs from what was intended. Because Netgen compares netlists, a broken extraction setup makes every comparison fail; checking extraction on a single known-good cell first saves time.

## Going deeper

The KLayout documentation on DRC and LVS scripting; the Magic and Netgen documentation by R. Timothy Edwards; the SKY130 PDK verification setup.
