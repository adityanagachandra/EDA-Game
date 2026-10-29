# Rare-event estimation

Memory bitcells, sense amplifiers, flip-flops in safety-critical logic and high-volume analog blocks all need failure rates far below what plain Monte Carlo can measure. Rare-event estimation is the family of methods that make such rates computable. A tool designer must know the main methods, what each assumes, and how each one fails silently, because a confidently wrong failure rate is worse than none.

## The cost of the naive approach

Plain Monte Carlo needs about $z^2/(p\,\varepsilon^2)$ samples for relative error $\varepsilon$, roughly 400 observed failures for 10 percent accuracy at 95 percent confidence. At 5 sigma, $p = 2.9 \times 10^{-7}$ and that is $1.3 \times 10^9$ simulations. The variance per sample is the problem: almost every sample lands where the circuit passes and contributes nothing. Every rare-event method attacks this by making failure samples common and then accounting exactly for the distortion. The figure of merit is the number of circuit simulations needed for a given relative error, and it is always quoted against the naive count.

## The main methods

**Importance sampling** draws from a proposal centred near the failure region, often at the most probable failure point or a mixture of several, and reweights. **Statistical blockade** trains a cheap classifier to block samples that are clearly far from the tail, simulates only the rest, and fits a generalized Pareto distribution to the exceedances beyond a threshold. **Scaled-sigma sampling** inflates all parameter sigmas so failures become frequent, measures the failure rate at several scale factors, and extrapolates a fitted model back to scale 1. **Subset simulation** writes the rare probability as a product of moderate conditional probabilities of nested events and estimates each with Markov chain sampling. Each trades the naive method's generality for an assumption about the shape of the failure region or the tail.

## How the estimates fail

Each method has a characteristic silent failure. A single-shift importance sampler misses a second failure region and underestimates, while its own variance estimate looks healthy because that variance is computed from the samples drawn. In high dimension, the weights degenerate and one sample carries the whole answer; the effective sample size $(\sum w)^2/\sum w^2$ exposes it. A blockade classifier that mislabels true tail samples discards them for ever. Scaled-sigma extrapolation is only as good as its assumed functional form. The defence is engineering, not cleverness: golden tests on limit states with known answers, multiple-region diagnostics, effective-sample-size checks, and periodic cross-checks against brute force at a moderate sigma where brute force is affordable.

## Going deeper

Singhee and Rutenbar, "Statistical blockade: very fast statistical simulation and modeling of rare circuit events and its application to memory design," IEEE TCAD, 2009; Au and Beck, "Estimation of small failure probabilities in high dimensions by subset simulation," Probabilistic Engineering Mechanics, 2001.
