# Seed noise and how many runs a QoR claim needs

Placers, routers and synthesis recipes are heuristic, and their results move with the random seed and with small, irrelevant perturbations such as cell ordering. A change that improves wirelength by 0.3% may be real or may be noise. Deciding needs the standard error and a sample size, and an agent that claims an improvement from one run is making a claim that the arithmetic does not support.

## Standard error of a mean QoR

If a metric varies across seeds with standard deviation $\sigma$, the mean of $n$ independent runs has standard error $\sigma / \sqrt{n}$. A detailed router whose wirelength varies with $\sigma = 0.4\%$ gives a single-run measurement uncertain by 0.4%, while the mean of 16 seeds is uncertain by $0.4 / 4 = 0.1\%$. An improvement of 0.3% measured from one run each of baseline and candidate is less than one standard deviation of the difference, which is $0.4\sqrt{2} = 0.57\%$. It is not evidence.

## Runs needed to resolve a change

To detect a mean difference $\delta$ between baseline and candidate, each measured with $n$ independent runs, at roughly $z$ standard errors of separation, require
$$z\,\sigma\sqrt{2/n} \le \delta, \qquad n \ge 2\left(\frac{z\,\sigma}{\delta}\right)^2.$$
With $\sigma = 0.4\%$, $\delta = 0.3\%$ and $z = 2$, $n \ge 2 \times (0.8/0.3)^2 = 14.2$, so 15 runs of each. Halving the effect to be detected quadruples the runs. This is why QoR claims on small effects need many seeds, and why fast engines that make many runs affordable matter.

## Paired runs on the same seeds

When baseline and candidate run on the same designs with the same seeds, much of the variation is shared: a seed that gives an unlucky initial placement hurts both. Analyzing the per-seed differences, with their standard deviation $\sigma_d$, needs only
$$n \ge \left(\frac{z\,\sigma_d}{\delta}\right)^2$$
paired runs. If pairing reduces the spread of the difference to $\sigma_d = 0.25\%$, then $n \ge (0.5/0.3)^2 = 2.8$, so 3 paired runs suffice where 15 unpaired runs of each were needed. Pairing works only if the engine is deterministic for a given seed, which is one more reason to demand determinism.

## Going deeper

Any statistics text on paired and two-sample comparisons, for example Wasserman, "All of Statistics" (2004), chapter 10.
