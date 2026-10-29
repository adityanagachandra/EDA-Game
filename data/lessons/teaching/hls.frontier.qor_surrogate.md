# Fast QoR surrogates for exploration

Every HLS run costs minutes; every implementation run costs far more. Design-space explorers therefore use surrogate models that predict latency, II and resources from a configuration in milliseconds, and spend real runs only on promising points. A surrogate is the fast engine in a fast-accurate loop, and it must be judged by how well it supports the decisions it drives, not only by its average error.

## Error of a single prediction

A surrogate maps a kernel and directive configuration to predicted QoR: from analytical models built on the schedule's bounds, such as $(N - 1)\,\mathrm{II} + L$ with ResMII and RecMII estimates, from learned models over features of the CDFG and directives, or from graph neural networks over the IR. For one configuration, the absolute percentage error is $|p - a| / a \times 100$, comparing prediction $p$ with the value $a$ from full synthesis. A prediction of 3520 cycles against an actual 4250 has an error of 17.2%. Analytical estimates are often exact for latency when the schedule is regular and badly off when the true binding constraint is something the model ignores, such as a memory port or a carried dependence through memory. Resource predictions are harder, because binding and logic optimization change them a lot.

## Mean absolute percentage error

Averaging the percentage errors over a validation set gives the mean absolute percentage error, $\mathrm{MAPE} = \frac{100}{n} \sum_i |p_i - a_i| / a_i$. For four designs with errors of 5.4%, 13.7%, 18.6% and 5.0%, MAPE is 10.7%. MAPE is scale-free, so a 50-cycle error on a 100-cycle design weighs more than on a 10,000-cycle one, which suits latency that spans orders of magnitude. The validation set must come from configurations the surrogate did not train on, ideally from kernels it has never seen, or the error estimate will be optimistic. A surrogate with a low MAPE on one kernel family can fail on another, which is why an explorer must periodically check it against real runs and fall back to the accurate engine when it drifts.

## Ranking quality matters more than error

An explorer uses a surrogate to choose which configurations to synthesize, so what matters is whether it orders candidates correctly. Count the discordant pairs: pairs of designs that the surrogate orders differently from full synthesis. Among five designs there are 10 pairs, and Kendall's $\tau$ is $(\text{concordant} - \text{discordant}) / 10$. A surrogate with a large constant bias has a high MAPE but can rank perfectly; one with small errors near the Pareto front can rank badly exactly where decisions are made. So explorers measure rank correlation near the front, synthesize the top candidates rather than trusting the single predicted best, and recurse into the accurate engine whenever the surrogate's ranking of candidates disagrees with the real results.

## Going deeper

Schafer and Wang, "High-level synthesis design space exploration: past, present, and future" (IEEE TCAD, 2020); the agent_first_eda lesson on fast and accurate engines in a modeling loop.
