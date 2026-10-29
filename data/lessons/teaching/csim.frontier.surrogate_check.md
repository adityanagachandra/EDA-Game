# Verifying a surrogate against SPICE

ML surrogates promise SPICE-like answers in microseconds: cell delays, amplifier gains, waveform features. Before an agent may use one in a loop, the surrogate must be verified against SPICE golden results with checkable acceptance rules. The metrics differ in what they reveal: a worst-case error bounds decisions, a failure count maps where the surrogate cannot be trusted, and an RMS error summarizes typical accuracy.

## Worst-case relative error

For each golden point with SPICE value $g$ and surrogate value $s$, the relative error is $|s - g|/|g|$. The worst case over the golden set is the number a decision can rely on. For pairs (238, 228), (87, 95) and (232, 244) ps, the errors are 4.2%, 9.2% and 5.2%, so the worst is 9.2%, set by the shortest delay, where a fixed absolute error is a larger fraction. Worst-case error is only as good as the golden set's coverage: points must span the input space the agent will query, including corners and extreme loads, or the bound says nothing about the region that matters.

## Counting tolerance failures

An acceptance rule such as "at most 3% relative error per point" turns verification into a count. Each failing point marks a region where the surrogate must not be used and the agent must fall back to SPICE, or where more training data is needed. For pairs (164, 166), (58, 70), (101, 107), (273, 263), (71, 66) and (193, 184) ps, the errors are 1.2%, 20.7%, 5.9%, 3.7%, 7.0% and 4.7%, so five of six fail. That surrogate is unusable at this tolerance, whatever its average. Failure locations matter more than the count: clusters at small delays or extreme inputs reveal systematic weaknesses rather than noise.

## RMS error over the golden set

The RMS error $\sqrt{\tfrac{1}{n}\sum (s_i - g_i)^2}$ summarizes typical accuracy in the metric's own unit. For differences of $-5, -1, 1, -10, -12$ and 4 ps, the mean square is $287/6 = 47.8$ ps$^2$ and the RMS error is 6.92 ps. RMS is what training minimizes, which is why a surrogate can look good on it while failing the worst-case or per-point rules: a few large errors hide among many small ones. A verification gate for an agent's surrogate should therefore combine metrics, and it should be rerun as a differential check whenever the surrogate or the SPICE models it imitates change.

## Going deeper

Hastie, Tibshirani and Friedman, "The Elements of Statistical Learning" (2nd ed., Springer, 2009), on model assessment; the ngspice manual, on .measure for generating golden metrics.
