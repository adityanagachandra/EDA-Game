# Board design rule checking

Board DRC checks the layout against the fabricator's geometric rules and the design's electrical constraints: clearance between copper on different nets, trace widths, annular rings, drill sizes, hole-to-hole spacing, mask slivers, courtyard overlaps, plus connectivity against the netlist. The geometric core is the same as chip DRC, distances between shapes found through a spatial index, but the rules are per net class and per region, and the shapes include arcs and arbitrary pads.

## Clearance as edge-to-edge distance

Clearance is the shortest distance between two copper shapes on different nets on one layer. For axis-aligned rectangles $a$ and $b$, compute the axis gaps $d_x = \max(0, b_{x1} - a_{x2}, a_{x1} - b_{x2})$ and likewise $d_y$. If one gap is zero, the rectangles overlap in that axis and the clearance is the other gap. Example: a pad from (13, 6) to (39, 46) and a trace from (43, 8) to (54, 36), in mil, overlap in y and have $d_x = 4$, so the clearance is 4 mil. If both gaps are zero, the shapes overlap, which on different nets is a short. Real shapes include rounded trace ends, circles and polygons, so production engines measure polygon-to-polygon distance, but the rectangle case is the core of many checks.

## Corner-to-corner distance

When both gaps are positive, the closest points are two corners, and the clearance is the Euclidean distance $\sqrt{d_x^2 + d_y^2}$, not the larger gap. Example: gaps of 3 and 7 mil give $\sqrt{9 + 49} \approx 7.62$ mil. A check that used $\max(d_x, d_y)$ would report 7 and pass a 7.5 mil rule that the shapes actually satisfy; one that used $d_x + d_y = 10$ would wrongly pass a 9 mil rule they violate. Many boards route 45-degree traces, whose clearance to pads is a point-to-segment distance. Fabricators measure clearance in the true Euclidean metric, so the DRC must too. Rule values depend on the net classes of both shapes, so the engine looks up the rule per pair before comparing.

## Scaling the check

Testing every pair of shapes is quadratic, and a large board has hundreds of thousands of shapes. Engines instead insert shapes into a **spatial index**, an R-tree or a grid, and query each shape's bounding box grown by the largest applicable clearance, testing only those candidates exactly. Same-net pairs are skipped, since copper on one net may touch. The check runs **incrementally**: when a trace moves, only shapes near its old and new positions are re-checked, which is what lets the interactive router stay DRC-clean as it pushes traces. Beyond clearance, board DRC covers annular ring (pad minus drill), drill-to-copper distance, minimum width, silkscreen over pads, and the connectivity check against the netlist, each a separate rule type in the engine.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification. The KiCad documentation on the DRC custom rules language.
