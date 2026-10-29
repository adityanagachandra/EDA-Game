# Weighted-average wirelength

The weighted-average (WA) model, proposed by Hsu, Chang and Balabanov, estimates a net's maximum as an exponentially weighted average of its pin coordinates rather than a log-sum. It has a smaller approximation error than log-sum-exp for the same $\gamma$, and it underestimates rather than overestimates. ePlace, RePlAce and DREAMPlace use it by default.

## The weighted-average maximum

Weight every pin by its softmax weight and average the coordinates:
$$\max_i x_i \approx \frac{\sum_i x_i e^{x_i/\gamma}}{\sum_i e^{x_i/\gamma}}.$$
Because it is a weighted average of the coordinates, it can never exceed the true maximum, and it approaches the maximum as $\gamma \to 0$. For x = 17, 14, 5, 7 with $\gamma$ = 4, the weights favour 17 and 14 and the estimate is 15.233, below 17. The minimum is estimated the same way with weights $e^{-x_i/\gamma}$, giving a value at or above the true minimum. Like LSE, it is computed stably by shifting exponents by the maximum (or minimum) first.

## The WA wirelength

The net's x wirelength is the estimated maximum minus the estimated minimum:
$$W_x^{WA} = \frac{\sum_i x_i e^{x_i/\gamma}}{\sum_i e^{x_i/\gamma}} - \frac{\sum_i x_i e^{-x_i/\gamma}}{\sum_i e^{-x_i/\gamma}}.$$
For x = 18, 11, 16 with $\gamma$ = 3, the maximum estimate is about 16.94 and the minimum estimate about 12.26, giving 4.68 against an exact span of 7. The worst-case error of WA is provably smaller than LSE's. Ties show the difference most clearly: when $k$ pins share the maximum, LSE overshoots by about $\gamma \ln k$, while WA averages the tied pins and returns the maximum exactly. The gradient is more complex than for LSE but still closed form, computed from the same per-net sums plus the weighted coordinate sums.

## WA versus LSE in practice

LSE always overestimates HPWL and WA always underestimates it, and at the same $\gamma$ WA is usually closer. For pins at 1, 7, 11 and 12 with $\gamma$ = 4, WA underestimates the span of 11 by 3.46 while LSE overestimates it by 4.28. WA is not convex in general, unlike LSE, but in practice nonlinear placers are not solving convex problems anyway once the density term is added. The accuracy gain translates into better final HPWL at the same runtime, which is why electrostatic placers adopted WA. For a tool designer the two are interchangeable behind one interface: a wirelength operator returning value and gradient, a natural plug-in seam.

## Going deeper

Hsu, Chang and Balabanov, "TSV-aware analytical placement for 3D IC designs" (DAC 2011), which introduced the weighted-average model; Hsu, Chen, Chang and Balabanov, IEEE TCAD 2013, for the error analysis.
