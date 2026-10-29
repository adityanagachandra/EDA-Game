# KLayout, Magic and Netgen in practice

Open-source physical verification rests on three tools. KLayout is a layout viewer and editor with DRC and LVS engines scripted in a Ruby-based language. Magic is an interactive layout editor built on corner stitching with continuous DRC and extraction. Netgen compares netlists for LVS. Together with open PDKs such as SKY130 and GF180MCU, they form the signoff layer of flows like OpenLane.

## KLayout

KLayout reads and writes GDSII and OASIS and serves as the common viewer for both. Its DRC and LVS engines are driven by scripts in a Ruby-based domain-specific language: input layers, derived layers by Booleans and sizing, checks such as width and space with a choice of metric, and output markers in a results database. Deep mode processes the layout hierarchically; tiled mode splits flat work into tiles for parallelism. Its LVS extracts devices and connectivity from layer rules and compares the result with a SPICE netlist.

## Magic

Magic stores each layer plane as corner-stitched tiles, which cover empty space as well as material and are linked at their corners, so local queries and edits are fast. Its design rules live in the technology file as edge-based rules, and it rechecks the edited region continuously, showing errors while the user draws. Magic also extracts layouts into its own extraction format, converted to SPICE for LVS and simulation. In open flows Magic often serves as a second, independent DRC and as the extractor feeding Netgen.

## Netgen and the open flow

Netgen compares two netlists, typically the extracted layout and the schematic, by iterative partition refinement with tie-breaking for symmetric classes. Its setup file sets comparison policy: property tolerances, permutable pins, and parallel or series merging. OpenLane uses Magic for DRC and extraction, Netgen for LVS, and KLayout for checks such as an XOR between streams written by different tools. Open PDKs ship decks for these tools, which makes the whole signoff flow inspectable, and a usable base for agent benchmarks that cannot use confidential decks.

## Going deeper

The KLayout documentation for DRC and LVS; the Magic and Netgen documentation; the OpenLane documentation.
