# Dies per wafer and edge loss

How many dies a wafer holds is the denominator of every cost-per-die calculation and an input to yield and DTCO trade-offs. Rectangular dies do not tile a circle, so the count is less than wafer area over die area, and the loss is worst for large dies. A standard closed-form estimate captures the edge loss well enough for early studies, and exact counts come from placing the die grid.

## Area ratio as an upper bound

The crude estimate divides wafer area by die area. A 300 mm wafer has area $\pi \times 150^2 = 70{,}686$ mm^2. A 10 mm by 10 mm die gives a ratio of 706.9. That number is an upper bound: it counts fractional dies at the edge and ignores the area that partial dies waste. It is still useful as a sanity check, and the gap between it and the real count grows with die size because the wasted ring around the edge is roughly one die deep. Larger dies waste a deeper ring.

## The edge-loss formula

A widely used estimate subtracts a term for the edge: $DPW = \pi d^2/(4A) - \pi d/\sqrt{2A}$, with wafer diameter $d$ and die area $A$. The second term approximates the number of dies lost along the circumference: the circumference $\pi d$ divided by the die diagonal $\sqrt{2A}$. For $d = 300$ mm and $A = 100$ mm^2: $706.9 - 942.5/14.14 = 706.9 - 66.6 = 640.2$, so about 640 dies. For a 400 mm^2 die the area term is 176.7 and the edge term 33.3, leaving 143: edge loss is 19% for the large die against 9% for the small one.

## Edge exclusion, scribe lanes and exact counts

Two corrections matter in practice. The edge exclusion removes a ring of width $e$, so the usable diameter is $d - 2e$. Scribe lanes, the saw streets between dies that also carry test structures, enlarge the die pitch to $(w + s)(h + s)$. A 10 by 10 mm die with 0.1 mm scribe on a 300 mm wafer with 3 mm exclusion uses $d = 294$ and $A = 102.01$, giving $665.6 - 64.7 = 600$ dies. Exact counts come from placing the die grid and counting dies entirely inside the usable circle; tools also optimize the grid offset, which can gain several dies. Die aspect ratio matters too, which floorplanners can exploit.

## Going deeper

de Vries, "Investigation of gross die per wafer formulas", IEEE Transactions on Semiconductor Manufacturing (2005); Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1 on the cost of an integrated circuit.
