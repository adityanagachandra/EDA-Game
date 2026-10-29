# Calibrating a fast mode against the accurate engine

A fast engine is trustworthy only to the extent its error against the accurate engine is known. Calibration measures that error on paired runs: the same designs or objects evaluated by both engines. From the pairs come a bias, a spread, a correction and a measure of whether the fast engine at least ranks candidates correctly, which is often all an agent needs.

## Bias and spread of the error

For paired values, fast $x_i$ and accurate $y_i$, the error is $e_i = x_i - y_i$. Its mean is the bias and its sample standard deviation is the spread. Suppose a fast timer reports slacks of $-20, -5, 10, 25$ ps on four endpoints and the signoff timer reports $-28, -11, 6, 21$ ps. The errors are $8, 6, 4, 4$, so the bias is 5.5 ps (the fast timer is optimistic) and the sample standard deviation is $\sqrt{11/3} = 1.91$ ps. Subtracting the bias removes the systematic error; the spread is what remains, and it sets how close to zero a corrected fast slack can be before it is untrustworthy.

## A linear correction

When the error grows with the value, a linear fit $y \approx a + b x$ by least squares does better than a constant bias: $b = S_{xy} / S_{xx}$ and $a = \bar{y} - b \bar{x}$, with $S_{xy} = \sum (x_i - \bar x)(y_i - \bar y)$. For the pairs above, $\bar x = 2.5$, $\bar y = -3$, $S_{xx} = 1125$ and $S_{xy} = 1230$, so $b = 1.093$ and $a = -5.73$. A fast slack of 0 ps then predicts a signoff slack of -5.73 ps. The fit is only valid inside the range of calibration data; extrapolating to a design far outside it is exactly the regime where the agent should recurse instead.

## Ranking agreement

An agent comparing candidates often needs only the right order, not the right values. For $n$ candidates there are $n(n-1)/2$ pairs; a pair is concordant if the fast and accurate engines order it the same way. The fraction of concordant pairs measures ranking fidelity (Kendall's tau is $2 \times$ that fraction minus 1, when there are no ties). Five placements give 10 pairs; if 8 are concordant, the fast engine ranks correctly 80% of the time. A fast mode with a large bias but perfect concordance is excellent for selection and poor for signoff; one with small bias but poor concordance is the reverse.

## Going deeper

Kendall, "A New Measure of Rank Correlation" (Biometrika, 1938); any regression text on least squares and prediction intervals.
