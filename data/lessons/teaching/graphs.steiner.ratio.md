# RSMT versus RMST and half-perimeter

Placers and timers pick a wirelength estimate per net: half-perimeter, minimum spanning tree, or Steiner tree. Each is cheaper and less accurate than the next. Knowing the exact bounds between them, and how they behave as pin count grows, is what lets a tool designer choose an estimate and know how much error it can carry into the next stage.

## MST over Steiner

For any rectilinear net, $\text{RSMT} \le \text{RMST} \le \tfrac{3}{2}\,\text{RSMT}$. The lower inequality holds because a spanning tree is a Steiner tree with no extra points; the upper is Hwang's theorem (1976). The ratio $\text{RMST}/\text{RSMT}$ measures how much a Steiner point saves.

A four-pin net with RMST 17 and RSMT 16 has ratio $17/16 \approx 1.06$. Most real nets have ratios close to 1, but the bound is tight: pins at $(0, 1)$, $(1, 0)$, $(2, 1)$, $(1, 2)$ have RSMT 4 (a cross through $(1, 1)$) and RMST 6, ratio exactly $3/2$.

## Half-perimeter as a lower bound

Any tree connecting the pins must span the bounding box horizontally and vertically, so $\text{HPWL} \le \text{RSMT}$, with equality for two and three pins. For larger nets the ratio $\text{RSMT}/\text{HPWL}$ grows: four pins spread around a box may need a trunk plus branches, so the tree is longer than the half-perimeter.

A four-pin net may have RSMT equal to HPWL when the pins lie along a monotone staircase, and noticeably more when they occupy all four sides of their bounding box. Analytical placers optimize HPWL or a smooth approximation of it, and some correct for the pin-count effect by scaling a net's HPWL by an empirical factor that grows with pin count, as in Cheng's RISA model.

## Using the bounds

The bounds convert one estimate into a guaranteed range for another. If a net's RSMT is 40, its RMST is at most 60. If its RMST is 39, its RSMT is at least $2/3 \times 39 = 26$, so an MST-based estimate overshoots the best Steiner tree by at most 50 percent. HPWL gives a lower bound on both.

In a flow, these ranges tell a designer where an estimate is safe. HPWL is exact for the large majority of nets, which have two or three pins, so its total error is small even though its worst case is large. For timing, length is not enough: a Steiner tree minimizing length can still give a long path from the driver to a critical sink, so timing-driven tree construction trades total length against source-sink path lengths.

## Going deeper

Hwang, "On Steiner minimal trees with rectilinear distance" (SIAM J. Applied Mathematics, 1976); Cheng, "RISA: accurate and efficient placement routability modeling" (Proc. ICCAD, 1994).
