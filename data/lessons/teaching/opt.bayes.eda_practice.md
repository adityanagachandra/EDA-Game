# Bayesian optimization for analog sizing and flow tuning

Textbook Bayesian optimization assumes a few continuous inputs, noise-free evaluations and one at a time. EDA tuning breaks every assumption: analog specs are constraints, flows have dozens of knobs and categorical modes, runs are noisy across seeds, licenses allow parallel jobs, cheap early-stopped runs are available, and results must generalize across designs. Each break has a standard remedy, and knowing them decides whether BO helps or wastes the compute budget.

## When BO fits an EDA tuning problem

BO pays off when each evaluation is expensive, gives no gradient, and the budget is tens to hundreds of runs: sizing an amplifier with SPICE in the loop, or tuning a place-and-route flow. It is the wrong tool for placing cells directly, where millions of continuous variables call for gradient methods on an analytical model.

Specifications are constraints, not objectives. An analog design must keep phase margin at least 60 degrees while minimizing power; constrained BO fits a separate surrogate for each spec and multiplies expected improvement by the modelled probability that every spec holds. Dimensionality limits GP surrogates: they work best up to roughly a few tens of inputs, so 40 flow knobs call for screening, embeddings or local trust regions. Categorical knobs, such as which CTS strategy or placer mode, need one-hot encodings or categorical kernels, never integer labels treated as real numbers, which would invent an ordering.

## Batches, noise and fidelity

Eight licenses mean eight runs per round. Batch acquisition proposes several diverse points at once, for example by fantasizing outcomes for pending runs before choosing the next, rather than eight copies of the current best or eight random points.

Tool runs are noisy: two seeds of the same configuration can differ by a few percent in QoR. The GP must include a noise term, or it interpolates seed luck and over-exploits lucky runs. Data from earlier designs in the same technology can warm-start or share a surrogate across tasks (transfer or multi-task BO). Multi-fidelity BO exploits cheap proxies: a placement stopped at 30% of its runtime, a fast timer, or a smaller design, spending most samples at low fidelity and confirming promising points at full fidelity. Each of these keeps the expensive full runs for decisions that need them.

## Judging BO results

A tuner that reports a configuration 1.5% better than default on its training design has proven little until the gain holds on held-out designs and other seeds. Overfitting one design and one seed is the commonest failure. Compared with random search at a budget of 50 runs over 10 knobs, BO usually finds better QoR, but neither gives a global guarantee, and the comparison should be run, not assumed.

BO optimizes exactly what it is told. If the objective is worst slack and DRC legality is not modelled, the tuner may trade legality for slack and return a configuration that fails signoff. The formulation must include every requirement as an objective term or a constraint. Finally, report the best evaluated configuration with its seeds and tool versions, not the surrogate's predicted optimum, which may never have been run.

## Going deeper

Lyu et al., "An efficient Bayesian optimization approach for automated optimization of analog circuits" (IEEE TCAS-I, 2018); Snoek, Larochelle and Adams, "Practical Bayesian optimization of machine learning algorithms" (NeurIPS, 2012), for batch and cost-aware acquisition.
