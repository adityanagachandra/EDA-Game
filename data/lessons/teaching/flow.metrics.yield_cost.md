# Yield and cost per good die

Not every die on a wafer works. Random defects, particles and process faults kill some fraction, and the probability that a die survives falls with its area. Yield turns dies per wafer into good dies per wafer, and dividing the wafer cost by good dies gives the cost of each working chip. These numbers make area a superlinear cost and explain why large chips are split into chiplets.

## The Poisson yield model

If killer defects land independently and uniformly at a density of $D_0$ per unit area, the number of defects on a die of area $A$ is Poisson with mean $D_0 A$. The die works only if it has none, so
$$Y = e^{-D_0 A}.$$
With $D_0 = 0.1$ defects per $cm^2$ and $A = 1\ cm^2$, $Y = e^{-0.1} \approx 0.905$, or 90.5%. For a $4\ cm^2$ die at the same density, $Y = e^{-0.4} \approx 0.670$. More realistic models (Murphy, negative binomial) account for clustered defects and give higher yields for large dies, but all agree that yield falls with area.

## Good dies per wafer

Multiply dies per wafer by yield and round down: 640 candidate dies at 90.5% yield give $\lfloor 640 \cdot 0.9048 \rfloor = 579$ good dies. Mind the units: defect density is usually quoted per $cm^2$ while die area is often given in $mm^2$, and $100\ mm^2 = 1\ cm^2$. Getting this wrong by a factor of one hundred is the most common mistake in yield arithmetic. Good dies per wafer is the number that matters for supply: in the example, a product that needs a million parts needs about $10^6/579 \approx 1728$ wafers, before test and packaging losses.

## Cost per good die

The cost of each good die is the wafer cost divided by good dies. A 10,000-dollar wafer with 579 good dies costs about 17.27 dollars per good die. Doubling the die area roughly halves dies per wafer and also lowers yield, so cost per good die more than doubles. That superlinear growth is why designers fight for every square millimetre, why redundancy (spare memory rows) is built in to repair defects, and why very large systems are increasingly split into several smaller dies in one package. Test and packaging costs then add per good die, not per wafer.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1; Stapper, "Fact and Fiction in Yield Modeling", Microelectronics Journal (1989).
