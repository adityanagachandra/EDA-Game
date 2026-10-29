# Poisson defect yield

Random particles and pattern defects kill dies independently of process variation. If killer defects land at random with a uniform density, the number on a die is Poisson, and the probability of a clean die is the zeroth Poisson term. This model is the starting point of every defect-limited yield calculation and of critical-area analysis, and it is the formula against which better models are compared.

## The basic formula

Let $D_0$ be the density of killer defects per unit area and $A$ the die area. If defects are placed independently and uniformly, the count on a die is Poisson with mean $\lambda = AD_0$, and the die is good with probability
$$Y = P(K = 0) = e^{-AD_0}.$$
Units must agree: with $A = 150$ mm$^2 = 1.5$ cm$^2$ and $D_0 = 0.2$ per cm$^2$, $\lambda = 0.3$ and $Y = e^{-0.3} = 74.1$ percent. Doubling the die area squares the yield (to 54.9 percent), which is the exponential penalty on large dies. Mixing mm$^2$ with defects per cm$^2$ is the most common error, a factor of 100 in the exponent.

## Inferring defect density

Fabs run the formula backwards. If dies of 200 mm$^2$ yield 70 percent, then $D_0 = -\ln Y/A = -\ln 0.7/2 = 0.178$ per cm$^2$. That inferred density is what a design team uses to predict the yield of a new die in the same process, and what a fab tracks over time as its process matures. The inference assumes all yield loss is from random defects; parametric and systematic losses inflate the apparent $D_0$, so fabs separate them with test structures and failure analysis before quoting a density.

## Critical area

Not every defect kills. A particle of a given size shorts two wires only if it lands where it touches both, and that region is the **critical area** $A_c$, which depends on the layout's spacing and widths. Each mechanism (metal shorts, metal opens, via failures) has its own critical area and density, and expected faults add:
$$Y = \exp\!\left(-\sum_j A_{c,j} D_j\right).$$
With 40 mm$^2$ of short-critical area at 0.3 per cm$^2$ and 25 mm$^2$ of via-critical area at 0.2 per cm$^2$, $\lambda = 0.12 + 0.05 = 0.17$ and $Y = 84.4$ percent. Critical-area extraction is a layout-analysis engine in its own right, and wire spreading or via doubling in design-for-manufacturing flows works by shrinking $A_c$.

## Going deeper

Stapper and Rosner, "Integrated circuit yield management and yield analysis: development and implementation," IEEE Transactions on Semiconductor Manufacturing, 1995; Koren and Koren, Proceedings of the IEEE, 1998.
