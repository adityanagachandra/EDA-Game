# Validating a surrogate model

A surrogate replaces an expensive simulation with a cheap fitted model inside an optimization or a Monte Carlo loop. Its training error says almost nothing about whether it can be trusted on the next design point. Validation measures prediction error on data the fit did not see, and it is what decides whether the fast loop can act on the surrogate or must call the accurate engine.

## Hold-out error

Set aside some simulated points that the fit never uses, predict them, and compute the root-mean-square error $\sqrt{\frac{1}{n}\sum (y_i - \hat{y}_i)^2}$. Errors of $+0.5, -1.2, 0.8, -0.3$ ps give $\sqrt{(0.25 + 1.44 + 0.64 + 0.09)/4} = 0.78$ ps. Compare that with what the surrogate is used for: a 0.78 ps error is harmless when ranking candidates that differ by 20 ps, and decisive when judging a path with 1 ps of slack. Report the error in the region where decisions are made, because a surrogate that is accurate on average can be poor near a spec boundary or in a tail.

## Leave-one-out without refitting

With few simulations, holding data out is wasteful. Leave-one-out cross-validation refits without each point in turn and predicts it, and for least-squares models it needs no refits at all:
$$e_{(i)} = \frac{e_i}{1 - h_{ii}},$$
where $e_i$ is the ordinary residual and $h_{ii}$ the leverage, the $i$-th diagonal of the hat matrix $H = X(X^TX)^{-1}X^T$. A point with residual 0.6 ps and leverage 0.8 has a leave-one-out residual of 3 ps: the fit bent to pass near it, so its small training residual hid a large prediction error. The sum of squared leave-one-out residuals is the PRESS statistic, a standard model-selection score.

## Model size against data

A polynomial surrogate of total degree $d$ in $k$ parameters has $\binom{k+d}{d}$ coefficients: 10 for a quadratic in 3 parameters, 91 for a quadratic in 12, 455 for a cubic in 12. A model with as many coefficients as simulations interpolates them exactly and validates terribly. A rule of thumb is several simulations per coefficient, plus regularization or sparse regression (ridge, lasso, least-angle) when the parameter count is large. The deeper risk is distribution shift: an optimizer drives the design toward regions where the surrogate is wrong in its favour. A validated surrogate is validated only on the region its data covered, and the loop that uses it must know where that region ends.

## Going deeper

Hastie, Tibshirani and Friedman, "The Elements of Statistical Learning" (2nd ed., 2009), chapters 3 and 7; Forrester, Sobester and Keane, "Engineering Design via Surrogate Modelling" (2008).
