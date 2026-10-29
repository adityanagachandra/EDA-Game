# Dies per wafer

A wafer is a disc, and dies are rectangles, so only a whole number of dies fit, and some area at the edge is always lost. The number of dies per wafer converts a die area from the floorplan into manufacturing cost. It is the first link between an EDA metric, area, and the economics that decide whether a product is viable.

## The area ratio

A wafer of diameter $d$ has area $\pi (d/2)^2$. A 300 mm wafer has $\pi \cdot 150^2 \approx 70{,}686\ mm^2$. Dividing by the die area $S$ gives an upper bound on the number of dies: for a $100\ mm^2$ die, about 706.9. That bound ignores the fact that rectangles cannot tile a circle: dies that cross the wafer edge are incomplete and cannot be used. For a 200 mm wafer the area is about $31{,}416\ mm^2$, less than half that of a 300 mm wafer, which is one reason the industry moved to larger wafers: the cost of processing a wafer grows more slowly than its area, so each die gets cheaper.

## The edge-loss correction

A standard approximation subtracts the dies lost around the circumference:
$$\mathrm{DPW} \approx \frac{\pi (d/2)^2}{S} - \frac{\pi d}{\sqrt{2S}}.$$
The second term counts roughly how many dies fit along the edge, the circumference divided by the die diagonal $\sqrt{2S}$ of a square die. For $d = 300$ mm and $S = 100\ mm^2$: $706.86 - 942.48/14.142 = 706.86 - 66.64 = 640.2$, so about 640 dies, rounded down. Larger dies lose a larger fraction to the edge: at $S = 400\ mm^2$ the formula gives $176.7 - 33.3 = 143.4$, about 143.

## Scribe lanes and die size

Dies are separated by scribe lanes, strips sawn away during dicing that also hold process test structures. The area that counts is the die plus its share of scribe lane: a $9.8 \times 9.8$ mm die with a 0.2 mm scribe lane occupies a $10 \times 10$ mm step, $100\ mm^2$. Because dies per wafer falls faster than linearly as area grows, every square millimetre the floorplan and placement save is worth more on a large die. That is one reason area efficiency of placement and routing tools is directly a cost metric.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1 on the cost of an integrated circuit.
