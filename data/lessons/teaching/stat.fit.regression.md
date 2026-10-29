# Regression for model fitting

Variation-aware EDA runs on fitted models: delay against load and slew, leakage against threshold voltage, a performance metric against process parameters. Each is a regression on simulation data. The least-squares machinery is familiar; what matters here is how characterization data is laid out, what the fitted coefficients mean physically, and where the fit can be trusted.

## Fitting a characterization line

Over a moderate range, cell delay is close to linear in output load: $d = b_0 + b_1 C$, where $b_0$ is the intrinsic delay and $b_1$ the effective drive resistance. Least squares picks
$$b_1 = \frac{S_{xy}}{S_{xx}} = \frac{\sum (x_i - \bar{x})(y_i - \bar{y})}{\sum (x_i - \bar{x})^2}, \qquad b_0 = \bar{y} - b_1\bar{x}.$$
For loads 2, 4, 6, 8 fF with delays 30, 34, 39, 42 ps, $\bar{x} = 5$, $\bar{y} = 36.25$, $S_{xy} = 41$ and $S_{xx} = 20$, so $b_1 = 2.05$ ps/fF and $b_0 = 26.0$ ps. A slope of 2.05 ps/fF corresponds to an effective resistance of 2.05 kilohm, a number a designer can sanity-check against the cell's drive strength.

## Prediction and extrapolation

The fitted line predicts $\hat{y} = b_0 + b_1 x$. Inside the characterized range this interpolation is reliable to about the residual scatter. Outside it, the prediction rests on the assumption that the physics stays linear, which delay does not: at very large loads the output slew degrades and delay curves upward, and at very small loads intrinsic effects dominate. Library characterization therefore samples load and slew on a grid that spans the expected range, and a timer that falls off the table edge is extrapolating. A tool should flag extrapolation rather than silently return the line's value, because the error is unbounded in principle.

## How good is the fit

The coefficient of determination $R^2 = S_{xy}^2/(S_{xx}S_{yy})$ is the fraction of response variance the line explains. In the example $S_{yy} = 84.75$, so $R^2 = 1681/1695 = 0.992$. A high $R^2$ does not prove the model form is right: a curved relationship sampled over a narrow range can show $R^2 > 0.99$ and still extrapolate badly. Residual plots against the input reveal curvature that $R^2$ hides. For surrogate models used inside optimization loops, the right measure is error on held-out data, because a flexible model can drive training $R^2$ to 1 while predicting new points poorly. That is the subject of surrogate validation.

## Going deeper

Montgomery, Peck and Vining, "Introduction to Linear Regression Analysis" (5th ed., 2012), chapters 2 and 4.
