# Design centering for yield

Given a variation model and a spec, the nominal design point can be moved to maximize the fraction of dies that pass. That is design centering, the classical form of yield optimization. It needs the yield as a function of the nominal point, which the parametric yield formula provides, and the sensitivities that tell how a design change moves the performance distribution.

## Centering the mean

For a Gaussian performance with fixed sigma and a two-sided spec from $L$ to $U$, yield is maximized with the mean at the spec centre $(L + U)/2$, because the Gaussian density is symmetric and unimodal. Off-centre designs lose yield on the near side faster than they gain it on the far side. With a spec from 90 to 110, $\sigma = 4$ and the mean at 95, the yield is $\Phi(3.75) - \Phi(-1.25) = 0.9999 - 0.1056 = 89.4$ percent; centred at 100 it is $\Phi(2.5) - \Phi(-2.5) = 98.8$ percent, a gain of 9.3 points with no change to the variation. With several correlated performances the acceptance region is a polytope in performance space and the optimal centre is no longer simply the midpoint.

## The sigma a spec can afford

Turn the question around: for a centred design, how much spread does a yield target allow? Each tail may hold $(1 - Y)/2$, so the half-width of the spec must be $z = Q^{-1}((1 - Y)/2)$ sigmas and
$$\sigma_{\max} = \frac{U - L}{2z}.$$
For a window of 20 units and $Y = 99$ percent, $z = 2.576$ and $\sigma_{\max} = 3.88$. This is the budget a designer hands back to the variation analysis: if the predicted sigma exceeds it, centering alone cannot meet the target, and the design must reduce sensitivity (bigger devices, a more robust topology) or the spec must be renegotiated.

## Moving the nominal with sensitivities

The designer does not move the performance mean directly; it changes a design parameter $x$ (a width, a bias current), and the performance moves by about $s\,\Delta x$ where $s = \partial P/\partial x$. To move a mean from 95 to 100 with $s = 2.5$ units per unit of $x$ takes $\Delta x = 2$. In general the change also alters the spread, because sensitivities to process parameters depend on the design point. Yield optimizers therefore iterate: estimate yield and its gradient with respect to design parameters (from sensitivities or Monte Carlo with likelihood-ratio gradients), take a step, and re-estimate. The objective is noisy when Monte Carlo estimates it, so the step must be judged against the confidence interval of the yield.

## Going deeper

Director and Hachtel, "The simplicial approximation approach to design centering," IEEE Transactions on Circuits and Systems, 1977; Graeb, "Analog Design Centering and Sizing" (2007).
