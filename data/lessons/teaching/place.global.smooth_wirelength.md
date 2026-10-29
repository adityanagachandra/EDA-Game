# Smooth wirelength models

Gradient-based placers need a wirelength that is differentiable, but HPWL is built from max and min, which have kinks wherever two pins tie for an extreme. Smooth wirelength models replace max and min by smooth functions controlled by a parameter $\gamma$. Log-sum-exp, used by APlace, NTUplace and ePlace, was the first widely used choice; weighted-average came later and is more accurate.

## Smoothing the maximum

The function $\gamma \ln \sum_i e^{x_i/\gamma}$ is a smooth approximation of $\max_i x_i$. It is always at least the maximum and at most the maximum plus $\gamma \ln n$, and it approaches the maximum as $\gamma \to 0$. Compute it stably by factoring out the largest value: $\max + \gamma \ln \sum_i e^{(x_i - \max)/\gamma}$, which avoids overflow in the exponentials. For x = 14, 0, 18, 18 and $\gamma$ = 3, the terms are $e^{-4/3}$, $e^{-6}$, 1 and 1, summing to 2.266, so the smooth maximum is $18 + 3\ln 2.266 = 20.454$. Two pins tied at the top push it up by $\gamma \ln 2 = 2.08$; that excess is the price of smoothness.

## Log-sum-exp wirelength

Apply the same idea to the minimum, $\min_i x_i = -\max_i(-x_i)$, and add:
$$W_x^{LSE} = \gamma \ln \sum_i e^{x_i/\gamma} + \gamma \ln \sum_i e^{-x_i/\gamma}.$$
The y part is the same with y coordinates, and the net's estimate is their sum. For pins at x = 10, 18, 13, 1 with $\gamma$ = 4, the smooth maximum is about 19.45 and the smooth negated minimum about -0.37, so $W_x^{LSE} \approx 19.07$ against an exact span of 17. The function is convex, smooth and cheap: each net costs two passes over its pins. Placers sum it over nets with net weights, exactly like HPWL.

## Choosing gamma

Each of the two terms errs by at most $\gamma \ln n$, so LSE overestimates the span by at most $2\gamma \ln n$. A small $\gamma$ gives accuracy but a nearly non-smooth function with large curvature, which forces small optimisation steps; a large $\gamma$ is smooth but inaccurate, and pulls on interior pins that HPWL ignores. Placers therefore tie $\gamma$ to the bin size and shrink it as placement proceeds: ePlace and RePlAce scale it with the density overflow, so it is large when cells are clumped and the exact wirelength matters little, and small near convergence. The x and y coordinates should be in consistent units with $\gamma$, typically bin widths.

## Going deeper

Naylor, Donelly and Sha, US patent 6,301,693 (2001), on log-sum-exp wirelength; Kahng and Wang, "Implementation and extensibility of an analytic placer" (IEEE TCAD, 2005), on APlace.
