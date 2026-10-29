# Layout netlist extraction

After devices are recognized, the extractor builds the nets: which shapes are electrically connected through overlaps and vias. The result is a netlist of devices and nets that LVS compares with the schematic. Extraction is a connected-components problem over a huge geometric graph, and its mistakes, real or false, are what LVS reports as shorts and opens.

## Connectivity by union-find

The deck declares which layers connect: shapes on the same conducting layer connect where they overlap or touch, and a via connects the shapes above and below it that contain it. Crossing without a via does not connect. The extractor runs union-find over shapes: one union per overlapping same-layer pair (found by a sweep) and one per via. The number of disjoint sets is the number of nets. Seven wire shapes with three vias may collapse into three nets. Union-find with path compression runs in nearly linear time, so connectivity is cheap compared with finding the overlaps.

## Labels, shorts and net names

Text labels attach names to the shapes they sit on. After extraction, each net collects the names of the labels on it. A net carrying two different names, such as A and Y, is reported as a short: somewhere shapes that should be separate are connected. A schematic net whose name appears on two different layout nets suggests an open. The same name repeated on one net is fine. Label attachment follows the deck's rules for text layers, so a label on the wrong text layer names nothing.

## Critical vias and opens

A via is critical if removing it alone splits a net: it is a bridge of the connectivity graph. Removing a critical via creates an open, while a via with a parallel path through other shapes can fail without changing the netlist. Counting critical vias is a bridge computation: remove each via, recount components, or use a linear-time bridge-finding search. The same analysis motivates redundant vias for yield and explains how a single missing via in a large net shows up in LVS as one schematic net matched to two layout nets.
