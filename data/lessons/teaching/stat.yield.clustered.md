# Clustered defects and negative binomial yield

Real defects are not spread uniformly: they cluster on some wafers and some regions. Clustering means some dies collect several defects while others get none, so for the same average density more dies are clean than the Poisson model predicts. The negative binomial model captures this with one extra parameter and is the standard industrial yield formula.

## The negative binomial model

Let the defect density itself vary from die to die, following a gamma distribution with mean $D_0$ and shape $\alpha$. Mixing a Poisson count over a gamma-distributed rate gives a negative binomial count, and the probability of zero defects is
$$Y = \left(1 + \frac{AD_0}{\alpha}\right)^{-\alpha}.$$
The clustering parameter $\alpha$ is small for strong clustering and tends to infinity for none. With $AD_0 = 1$ and $\alpha = 2$, $Y = 1.5^{-2} = 44.4$ percent. As $\alpha \to \infty$, $(1 + \lambda/\alpha)^{-\alpha} \to e^{-\lambda}$, recovering Poisson. Typical fitted values are of order 1 to a few; the parameter is fitted from wafer maps, not assumed.

## Comparing with Poisson

For the same average density, clustering always raises yield: concentrating defects on fewer dies leaves more dies clean. With $AD_0 = 1$, Poisson gives $e^{-1} = 36.8$ percent while the negative binomial with $\alpha = 2$ gives 44.4 percent and with $\alpha = 0.5$ gives $3^{-0.5} = 57.7$ percent. The gap grows with $AD_0$. A yield model that uses Poisson for large dies therefore underestimates yield, and a cost model built on it over-penalizes large dies. Other historical models, such as Murphy's $Y = ((1 - e^{-AD_0})/(AD_0))^2$, sit between the two and reflect other assumptions about the density distribution.

## Scaling the die

Yield prediction for a new product scales area in the same process. Under Poisson, multiplying area by $m$ raises the yield to the power $m$. Under the negative binomial, $Y(mA) = (1 + mAD_0/\alpha)^{-\alpha}$, which falls more slowly. With $AD_0 = 0.5$ and $\alpha = 2$, the original die yields $1.25^{-2} = 64$ percent and a die four times larger yields $2^{-2} = 25$ percent; Poisson would predict $e^{-2} = 13.5$ percent for the larger die. Clustering is measured at a block scale, so the model is calibrated for die sizes near those it was fitted on, and very large dies or chiplet assemblies need care.

## Going deeper

Stapper, "Fact and fiction in yield modeling," Microelectronics Journal, 1989; Koren and Koren, "Defect tolerance in VLSI circuits," Proceedings of the IEEE, 1998.
