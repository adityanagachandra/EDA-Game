# Deciding where to recurse to the accurate engine

The fast loop produces estimates; the accurate engines produce answers. Recursing everywhere wastes the fast loop's speed, and never recursing trusts numbers that may be wrong. With a calibrated bias and spread, the decision becomes a calculation: how likely is the fast verdict to be wrong for this object, and where does an accurate run buy the most certainty?

## A calibrated error band

If the fast engine's error has bias $b$ and standard deviation $\sigma$, approximately normal, then for a fast slack $m$ the accurate slack is approximately normal with mean $m - b$ and standard deviation $\sigma$. The z-score of the corrected value, $z = (m - b)/\sigma$, says how many standard deviations it lies from the pass or fail threshold at zero. The probability that the true slack is negative is $\Phi(-z)$. With $b = 4$ ps and $\sigma = 6$ ps, a fast slack of 13 ps gives $z = 1.5$ and a 6.7% chance of an actual violation; a fast slack of 25 ps gives $z = 3.5$ and a 0.02% chance.

## Counting candidates that need recursion

A simple policy recurses on every object with $|z| < k$: those whose corrected estimate is within $k$ standard deviations of the threshold, where the fast verdict could be wrong. With $b = 4$, $\sigma = 6$, $k = 2$ and fast slacks of $-30, -10, 2, 9, 15, 40$ ps, the z-scores are $-5.67, -2.33, -0.33, 0.83, 1.83, 6$, so three endpoints (2, 9 and 15 ps) need the accurate engine and the rest can be decided from the fast run. Raising $k$ to 3 adds the -10 ps endpoint. The choice of $k$ trades accurate-engine time against the risk of a missed violation.

## Spending a recursion budget

When the accurate engine can run on only $R$ objects, run it on the $R$ with the smallest $|z|$, the most uncertain ones, and accept the fast verdict elsewhere. The expected number of wrong verdicts among the objects left to the fast engine is $\sum \Phi(-|z_i|)$. In the example with $R = 2$, the two most uncertain endpoints ($z = -0.33$ and 0.83) go to the accurate engine; the remaining expected errors are about $\Phi(-1.83) + \Phi(-2.33) + \cdots \approx 0.034 + 0.010 = 0.04$. Variance that grows during the loop raises $\sigma$, shrinks every $|z|$ and should trigger more recursion, or a recalibration.

## Going deeper

The agent-first EDA notes in this curriculum, section on fast and accurate engines; any statistics text on the standard normal distribution.
