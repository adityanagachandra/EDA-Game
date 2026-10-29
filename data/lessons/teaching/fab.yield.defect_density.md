# Defect density and yield models

Random particles and process defects land on wafers with some average density. A die that receives a killer defect fails. Turning a defect density into a yield needs a statistical model, and the simplest, the Poisson model, already shows why yield falls exponentially with die area, which is the main economic pressure on die size and on chiplets.

## Expected defects per die

Let $D_0$ be the density of killer defects per unit area, and $A$ the die area. The expected number of killer defects per die is $\lambda = A D_0$. Units must match: die area is usually given in mm^2 and defect density in defects per cm^2, and 100 mm^2 is 1 cm^2. A 150 mm^2 die at 0.2 defects per cm^2 has $\lambda = 1.5 \times 0.2 = 0.3$. A mature process may run at around 0.1 per cm^2 or lower; an early process runs higher, and that early defect density is a major reason large dies are expensive on new nodes.

## The Poisson yield model

If defects land independently and uniformly, the count on a die is Poisson with mean $\lambda$, and the die works only with zero defects: $Y = e^{-\lambda} = e^{-A D_0}$. For $\lambda = 0.3$, $Y = 74.1\%$. Doubling the die area to $\lambda = 0.6$ gives 54.9%: the yield of the larger die is the square of the smaller one's. That exponential dependence is the core of yield economics. Combined with dies per wafer, it means a die twice the area yields well under half as many good dies per wafer.

## Inverting yield to learn defect density

Foundries and product teams run the model backwards: from measured yield infer an effective defect density, $D_0 = -\ln Y / A$. A 200 mm^2 die yielding 80% implies $D_0 = -\ln 0.8/2 = 0.112$ per cm^2. That inferred density lumps together all mechanisms, so it is only meaningful if systematic losses have been separated out. The Poisson model also tends to be pessimistic for large dies because real defects cluster; clustered models, discussed next, give higher yield at the same mean density. Product teams track the inferred density over time as a learning curve.

## Going deeper

Stapper, "Fact and fiction in yield modeling", Microelectronics Journal (1989); Koren and Koren, "Defect tolerance in VLSI circuits: techniques and yield analysis", Proceedings of the IEEE (1998).
