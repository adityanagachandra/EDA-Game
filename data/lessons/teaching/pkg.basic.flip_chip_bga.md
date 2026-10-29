# Flip chip and BGA packages

In flip chip the die is turned face down and joined to the package by solder bumps spread over its whole area. The package then reaches the board through a ball grid array (BGA): an area array of solder balls on its underside. Both are area arrays, so I/O count grows with area rather than perimeter, and the short bumps have far less inductance than bond wires. Bump and ball counts are the first numbers a package planner computes.

## Area-array bump count

Flip-chip bumps sit on a regular grid at pitch $p$, often 130 to 200 um for solder bumps and copper pillars on mainstream parts. If the first bump is placed $p/2$ from each edge, a die of width $W$ and height $H$ holds $\lfloor W/p \rfloor \times \lfloor H/p \rfloor$ bumps. A 10 mm by 8 mm die at 150 um pitch holds $\lfloor 10000/150 \rfloor \times \lfloor 8000/150 \rfloor = 66 \times 53 = 3498$ bumps. Not all are signals: typically more than half carry power and ground, spread across the die so the on-chip power grid is fed from directly above every region. That spread is the electrical reason for an area array, beyond the raw count.

## Area array against perimeter

Take a square die whose side holds $n$ sites at a given pitch. An area array uses all $n^2$ sites, while a perimeter ring at the same pitch uses only the outer ring, $4(n-1)$ sites. For $n = 50$ that is 2500 against 196, a ratio of about 12.8, and the ratio grows linearly with $n$. Flip chip also shortens the electrical path: a bump is tens of micrometres tall with tens of picohenries of inductance, against a nanohenry or more for a bond wire. The costs are a more expensive substrate that must escape-route the dense array, underfill to protect the bumps from CTE strain, and a die whose back side, not its face, now carries the heat path.

## BGA ball maps

A BGA package of side $L$ with ball pitch $p$ has an $\lfloor L/p \rfloor$ by $\lfloor L/p \rfloor$ grid; common pitches are 1.0, 0.8, 0.65 and 0.5 mm. Many packages depopulate a central $k \times k$ block, for example to leave room under the die for thermal vias or to ease board routing, so the ball count is $\lfloor L/p \rfloor^2 - k^2$. A 35 mm package at 1.0 mm pitch with a 9 by 9 hole has $35^2 - 81 = 1144$ balls. The ball map is a co-design object: the package designer wants power balls under the die, the board designer wants signal balls on the outer rows where they escape without vias, and the ball pitch sets which board via and trace technology is needed.

## Going deeper

Lau, "Flip Chip Technologies" (1996); JEDEC publication 95 for BGA outlines; Tummala, "Fundamentals of Microsystems Packaging", chapter on flip chip.
