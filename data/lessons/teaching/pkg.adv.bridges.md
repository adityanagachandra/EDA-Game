# Embedded silicon bridges

An embedded silicon bridge is a small piece of silicon with fine wiring, placed only under the facing edges of two dies, inside an organic substrate or a mold layer. It gives die-to-die wire density close to a full silicon interposer where it is needed, while the rest of the package stays organic. The trade-offs are silicon area, assembly tolerance and the number of edges a bridge can serve.

## Silicon area saved

A full interposer covers the whole multi-die footprint, while bridges cover only the edge regions where dies talk. For a 1600 mm^2 interposer replaced by four bridges of 8 mm by 3 mm, the silicon area falls from 1600 to $4 \times 24 = 96$ mm^2, a ratio of about 16.7. Less silicon means lower cost and no reticle-size limit on the package, because each bridge is small. What is lost: the bridge provides no through-silicon vias, so power and other I/O reach the dies through the organic substrate at coarse pitch, and the full interposer's freedom to route between any two points on the footprint is gone. Bridges suit dies that talk mostly to neighbours.

## Wires per bridge

A bridge carries wires across the gap between two dies. If it spans an edge length $E$ under both dies with wiring pitch $p$ and $N$ layers, it carries $N \lfloor E/p \rfloor$ wires, and $n$ bridges carry $n$ times that. With a 6 mm edge, 2 um pitch, 2 layers and three bridges, total wires are $3 \times 2 \times 3000 = 18000$. The pitch is coarser than a full interposer's best, because the bridge is often built with relaxed rules to lower cost, but still far finer than organic wiring. As with any interposer, the microbumps over the bridge, not the bridge wires, usually set the usable count.

## Bandwidth through bridges

Bandwidth is wires times data rate. If 4000 signal wires per bridge each run at 4 Gb/s, one bridge carries 16 Tb/s, and three bridges 48 Tb/s. The design question is whether that matches the traffic between the dies. Bridges also impose assembly constraints: the bridge sits at a slightly different height and position than the substrate pads around it, so the bump heights and placement tolerances must absorb both, and the die must bond to the fine bridge bumps and the coarse substrate bumps in the same reflow. An EDA flow must model the bridge as a separate routing layer set with its own rules, and floorplan the dies so that the communicating edges face each other over a bridge.

## Going deeper

Mahajan and others, "Embedded multi-die interconnect bridge", IEEE ECTC (2016); Lau, "Semiconductor Advanced Packaging" (2021), chapter on bridges.
