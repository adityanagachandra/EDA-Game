# Clustered defects and negative binomial yield

Real defects are not spread uniformly: they cluster on some wafers, in some regions and on some dies. Clustering leaves more dies with zero defects than the Poisson model predicts. The negative binomial model captures this with one extra parameter and is the standard industrial yield model for random defects.

## A gamma-mixed Poisson

Suppose each die's defect count is Poisson, but its mean varies from die to die following a gamma distribution with mean $A D_0$ and shape $\alpha$. Averaging the Poisson zero-probability over that gamma gives $Y = (1 + A D_0/\alpha)^{-\alpha}$. Small $\alpha$ means strong clustering; as $\alpha \to \infty$ the gamma collapses to its mean and the model returns to Poisson, since $(1 + x/\alpha)^{-\alpha} \to e^{-x}$. With $A D_0 = 0.6$ and $\alpha = 2$: $Y = (1.3)^{-2} = 59.2\%$. Typical fitted values of $\alpha$ lie between about 0.5 and 5.

## Comparing with Poisson

At the same mean defect count, clustering raises yield because some dies receive several defects while others receive none. For $A D_0 = 0.6$: Poisson gives 54.9%, the negative binomial with $\alpha = 2$ gives 59.2%, a 4.3-point difference; with $\alpha = 0.5$ it gives $(2.2)^{-0.5} = 67.4\%$. The gap grows with die area. A cost model that uses Poisson will overstate the penalty of a large die when defects cluster strongly. Choosing $\alpha$ is a fitting problem: the foundry fits it to wafer maps, often by counting how yield changes with the size of a block of dies.

## Inferring defect density with clustering

Inverting the model gives $D_0 = \alpha(Y^{-1/\alpha} - 1)/A$. For a 200 mm^2 die yielding 80% with $\alpha = 2$: $0.8^{-0.5} = 1.118$, so $D_0 = 2 \times 0.118/2 = 0.118$ per cm^2, slightly higher than the Poisson estimate of 0.112. The two models agree when yield is high and disagree as yield falls. Murphy's and Seeds' models are older alternatives built from other assumed distributions of defect density. For design decisions, what matters is consistency: compare options with the same model and parameters, and note that chiplet trade-offs are sensitive to the model chosen.

## Going deeper

Stapper, Armstrong and Saji, "Integrated circuit yield statistics", Proceedings of the IEEE (1983); Cunningham, "The use and evaluation of yield models in integrated circuit manufacturing", IEEE Transactions on Semiconductor Manufacturing (1990).
