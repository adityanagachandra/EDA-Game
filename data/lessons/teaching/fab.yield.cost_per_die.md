# Cost per good die

The economic figure of merit for silicon is the cost of a die that works. It combines wafer cost, dies per wafer and yield in one ratio. Every DTCO, chiplet and die-size decision is ultimately compared through it, and its non-linear dependence on area is why architects care about every square millimetre.

## From wafer cost to die cost

A processed wafer has a cost $C_w$ that covers depreciation, materials, labour and masks amortized over volume; advanced-node wafers cost many thousands of dollars. Before yield, each candidate die costs $C_w/DPW$. A wafer costing 10,000 US dollars with 640 candidate dies gives 15.63 dollars per candidate. Testing, packaging and final test add further costs per die, but the silicon cost is where die area enters most strongly, through both the dies-per-wafer count and the yield. Mask cost, by contrast, is fixed per design and is amortized over volume, so it dominates for low-volume products and vanishes per die for high-volume ones.

## Dividing by good dies

Only good dies can be sold, so cost per good die is $C_{good} = C_w/(DPW \cdot Y)$. With the same wafer, a 100 mm^2 die at 0.2 defects per cm^2 has $Y = e^{-0.2} = 0.819$, so $640 \times 0.819 = 524$ good dies and 19.08 dollars each. A 200 mm^2 die on the same wafer gets about 306 candidates and $Y = e^{-0.4} = 0.670$, so 205 good dies at 48.75 dollars: doubling the area raised cost per good die by a factor of 2.55. Area is penalized twice, once in the count and once in the yield.

## Shrinks and node transitions

A node transition shrinks a design but raises wafer cost. If the die shrinks to 60% of its area while wafer cost rises by 30%, cost per good die changes by $1.3 \times (DPW \cdot Y)_{old}/(DPW \cdot Y)_{new}$. For the 100 mm^2 example, 60 mm^2 gives about 1092 candidates and yield $e^{-0.12} = 0.887$, so 968 good dies at 13.42 dollars, a ratio of 0.70. If defect density on the new node starts higher, the gain shrinks or reverses. This arithmetic, applied to block-level area from DTCO studies, decides whether a technology option is worth its added process cost.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1 on the cost of an integrated circuit; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits", chapter 1.
