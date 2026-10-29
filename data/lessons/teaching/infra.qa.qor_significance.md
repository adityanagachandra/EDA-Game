# Seed noise and QoR significance

Placers and routers use randomness: initial positions, move orders, tie-breaking seeds. Change the seed and wirelength moves by a fraction of a percent with no change to the code. A regression that reports a 0.3% improvement is meaningless unless that number is compared with the seed-to-seed noise. Running several seeds before and after a change, and testing the paired differences, separates real improvements from luck. It is ordinary statistics, applied to the most common question in engine development.

## Paired differences across seeds

Run the old and the new build with the same seeds $1, \ldots, n$ and record the metric for each. Pairing removes the variation that the seed causes in both builds alike: a seed that happens to give a good placement gives a good one before and after. Analyze the per-seed differences $d_i = \text{after}_i - \text{before}_i$.

Example: routed wirelength in mm for four seeds, before → after: 500 → 497, 512 → 510, 488 → 487 and 505 → 498. The differences are $-3, -2, -1, -7$, and their mean is $\bar d = -13/4 = -3.25$ mm. A negative mean suggests that the change shortens wirelength. Whether -3.25 is real depends on how much the differences scatter. Pairing is legitimate only if the seed has the same meaning in both builds; if the change consumes random numbers in a different order, the pairing weakens and more seeds are needed.

## Seed-to-seed spread

The sample standard deviation of the differences measures their scatter:
$$s = \sqrt{\frac{1}{n - 1}\sum_{i=1}^{n} (d_i - \bar d)^2}.$$
Dividing by $n - 1$ rather than $n$ makes $s^2$ an unbiased estimate of the variance, which matters for the handful of seeds a regression can afford.

For the example, deviations from -3.25 are $0.25, 1.25, 2.25, -3.75$, with squares $0.0625, 1.5625, 5.0625, 14.0625$ summing to 20.75. Then $s = \sqrt{20.75/3} \approx 2.63$ mm. A spread of 2.63 mm against a mean of 3.25 mm says that individual seeds disagree considerably about the size of the effect. Knowing the typical spread of each metric on each design also sets sensible tolerances for the nightly regression: a band narrower than the seed noise would fail on noise alone.

## The paired t statistic

The paired t statistic divides the mean difference by its standard error:
$$t = \frac{\bar d}{s / \sqrt{n}}.$$
For the example, $t = -3.25 / (2.63/\sqrt 4) = -3.25 / 1.315 \approx -2.47$. Under the null hypothesis of no effect, $t$ follows a t distribution with $n - 1 = 3$ degrees of freedom, whose two-sided 5% critical value is about 3.18. Since $|t| = 2.47 < 3.18$, four seeds are not enough to call the improvement real at that level, though they hint at one.

More seeds shrink the standard error as $1/\sqrt n$, so quadrupling the seeds halves it. The practical lesson for tool builders: report QoR changes with their uncertainty, choose the number of seeds from the noise of each metric, and remember that testing many metrics on many designs at once produces false alarms unless the threshold accounts for it.

## Going deeper

Montgomery, "Design and Analysis of Experiments", chapter 2 on paired comparisons; Kahng and Mantik, "On Mismatches between Incremental Optimizers and Instance Perturbations in Physical Design Tools" (ICCAD 2000), on noise in placement results.
