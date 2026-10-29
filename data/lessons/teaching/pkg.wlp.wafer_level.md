# Wafer-level and fan-out packaging

Wafer-level packaging builds the package on the wafer itself, before or instead of dicing into a separate substrate. In fan-in WLP the balls sit under the die, so the package is the size of the die. In fan-out WLP the dies are embedded in mold compound to form a larger reconstituted wafer or panel, and redistribution layers carry the I/O out beyond the die edge. The choice between them is a count of balls against area.

## Fan-in capacity

A fan-in WLP places solder balls directly on the die at a pitch such as 0.4 or 0.35 mm, routed from the pads by a redistribution layer on the wafer. The ball count is limited by die area: a 4 mm by 3 mm die at 0.4 mm pitch holds $\lfloor 4/0.4 \rfloor \times \lfloor 3/0.4 \rfloor = 10 \times 7 = 70$ balls. Fan-in is the smallest and cheapest package, ideal for small dies such as power management and radio parts. It has two limits. A large die puts far-from-centre balls under high thermal-mismatch strain against the board, and a small die with many I/O simply has no room. The second limit is what fan-out removes.

## Fan-out area expansion

If a die needs $B$ balls at pitch $p$ in a square full array, the package side must be at least $\lceil \sqrt{B} \rceil p$. A 5 mm by 5 mm die needing 400 balls at 0.4 mm pitch needs a $20 \times 0.4 = 8$ mm package, so the package area is $64/25 = 2.56$ times the die area. Fan-out WLP provides that extra area with mold compound around the die, and the RDL routes from die pads into it. Because the RDL is built with wafer-style lithography at a few micrometres line and space, fan-out also supports multi-die packages with fine die-to-die wiring and no laminate substrate, which is why it is used for application processors and for chiplet packages.

## Reconstituted wafers and panels

Fan-out starts by placing known-good dies on a carrier and overmolding them into a reconstituted wafer or a rectangular panel. Units are laid out at a step equal to package size plus a saw street, so a panel of $W_p \times H_p$ holds $\lfloor W_p/s \rfloor \times \lfloor H_p/s \rfloor$ units at step $s$. Panels raise throughput because a rectangle wastes less area on rectangular packages than a round wafer. Two process effects matter to a designer. Die shift: dies move slightly during molding, so the RDL lithography must align to each die's actual position, sometimes by adaptive patterning. Warpage: mold compound and silicon have very different CTEs, so the reconstituted wafer bows after cure, which limits panel size and RDL alignment.

## Going deeper

Lau, "Fan-Out Wafer-Level Packaging" (2018); Fan, "Chip-first and chip-last fan-out", in the IEEE Heterogeneous Integration Roadmap.
