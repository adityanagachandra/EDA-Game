# Fast versus accurate timing estimation

FPGA flows estimate timing at every stage: from LUT depth during mapping, from a distance-based delay model during placement, and exactly only after routing. An agent running a continuous loop uses the fast estimates to decide what to try and recurses to the accurate timer where the estimate cannot be trusted. Measuring the estimator's error and choosing the guard band are what make that loop safe.

## Measuring the estimator

Compare the fast estimate with post-route static timing on the same paths. The mean absolute error, $\frac{1}{n}\sum |\hat s_i - s_i|$, summarizes how far off the estimate is. With slack estimates and routed values of (190, 190), (-50, -130), (460, 400), (380, 390), (450, 370) and (330, 390) ps, the absolute errors sum to 290 ps, so the MAE is 48.3 ps. The sign of the error matters as much as its size: placement-stage models ignore congestion, so they tend to be optimistic, predicting more slack than routing delivers. Measuring the error per design class, and refreshing it as the tools change, is part of maintaining the fast engine.

## Guard bands decide when to recurse

An agent trusts the fast estimate only when its slack exceeds a guard band $g$ and sends every other path to the accurate timer. With $g = 60$ ps and estimates of 150, -90, -60, 330, -40, -150 and 190 ps, four paths fall below 60 and go to the accurate engine. A larger guard band sends more paths, costing run time, while a smaller one risks trusting wrong estimates. A principled choice ties $g$ to the measured error distribution, for example a high percentile of the optimistic error. Placement and routing tools already work this way: they recompute criticalities from an approximate model frequently and from exact timing occasionally.

## False passes are the dangerous error

Two kinds of error are not symmetric. A pessimistic error, where the estimate predicts a violation that routing does not show, costs only extra work. A false pass, with estimated slack of 0 or more but negative routed slack, lets a violation through. With estimates and routed values of (-30, 50), (150, 100), (370, 380), (-20, -80), (20, -50) and (220, 280) ps, only the fifth path is a false pass. An agent's fast-accurate loop must be designed around the false-pass rate, not the average error: it is acceptable for the estimator to be noisy if the guard band makes false passes rare, and every promoted result is still signed off by the accurate timer.

## Going deeper

Marquardt, Betz and Rose, FPGA (2000), for placement-stage delay estimation.
