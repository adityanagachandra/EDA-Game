# Probability of hitting a corner case

Random stimulus finds a bug only if it reaches the scenario that triggers it. Corner cases, such as a full FIFO during a flush or an address at the top of a range with a zero length, often have tiny per-test probabilities. Computing those probabilities tells a verification engineer whether a regression can be expected to hit a corner at all, how many runs that takes, and when to stop waiting and bias the stimulus.

## The probability of a corner per draw

When fields are independent and uniform, the probability of a corner is the product of the per-field probabilities. An unconstrained 8-bit address is at its maximum with probability $2^{-8}$; an unconstrained 4-bit length is 0 with probability $2^{-4}$. Both at once: $p = 2^{-12} = 1/4096$. The number of draws until the first hit is geometric with parameter $p$, so its mean is $1/p = 4096$ transactions. Correlated fields, constraints and sequential preconditions change the numbers: if the corner needs a FIFO to be full, $p$ is the probability that random traffic fills it, which may be astronomically small with a fast reader. Estimating $p$ is the first step in deciding whether random stimulus is enough.

## At least one hit in N tests

If each test independently hits the corner with probability $p$, the chance that $N$ tests all miss is $(1 - p)^N$, so the chance of at least one hit is $1 - (1 - p)^N$. With $p = 2^{-10}$ and $N = 1000$: $(1 - 1/1024)^{1000} \approx 0.3764$, so the probability of at least one hit is about 62.36%. A useful approximation for small $p$ is $1 - e^{-Np}$, here $1 - e^{-0.977} \approx 62.3\%$. Running exactly $1/p$ tests therefore hits the corner only about 63% of the time, not with certainty: the expected wait is not a guarantee.

## Tests needed for a target confidence

To hit the corner with probability at least $c$, solve $1 - (1 - p)^N \ge c$: $N \ge \ln(1 - c) / \ln(1 - p)$, rounded up. For $p = 2^{-10}$ and $c = 0.95$: $N \ge \ln 0.05 / \ln(1 - 1/1024) \approx 2.9957 / 0.000977 \approx 3066.1$, so 3067 tests, about three times the mean wait. For 99% the multiplier is about 4.6. These numbers set regression budgets, and they show why rare corners are not left to luck: if $p = 2^{-30}$, the required $N$ is in the billions, and the right move is to bias the stimulus (with `dist`, constraint modes or a directed sequence) until $p$ is reasonable.

## Going deeper

Ross, "A First Course in Probability", the chapter on discrete distributions (geometric); Wile, Goss and Roesner, "Comprehensive Functional Verification" (2005), the chapter on random stimulus generation.
