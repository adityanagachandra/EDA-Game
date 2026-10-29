# Randomized algorithms and random simulation

Randomized algorithms use coin flips to be faster or simpler than any known deterministic method, at the price of a bounded probability of error or a random running time. In Boolean reasoning the most important one is random simulation: apply random input vectors to two circuits and compare. It finds most bugs almost for free, but it can only show that circuits differ, never that they are equal. Quantifying what a clean simulation run proves is a computation every equivalence-checking flow depends on.

## The chance that random vectors miss a difference

Suppose $f$ and $g$ differ on $d$ of the $2^n$ input vectors, so a uniform random vector exposes the difference with probability $p = d / 2^n$. With $t$ independent vectors, all of them miss with probability $(1 - p)^t \approx e^{-pt}$.

For $n = 12$, $d = 128$ and $t = 48$: $p = 128 / 4096 = 1/32$, and $(31/32)^{48} \approx 0.218$. One run in five misses the bug. The danger is the small-$d$ case: a corner-case bug on a handful of vectors in a 64-input block has $p$ around $10^{-18}$, and no feasible simulation will find it. This is a one-sided error: a mismatch found is certain, a mismatch not found proves nothing.

## How many vectors are enough

To push the miss probability below $\varepsilon$, solve $(1 - p)^t \le \varepsilon$: $t \ge \ln \varepsilon / \ln(1 - p)$, approximately $\ln(1/\varepsilon) / p$. For $n = 14$, $d = 128$ and $\varepsilon = 0.1$, $p = 1/128$ and $t = \lceil \ln 0.1 / \ln(127/128) \rceil = 294$.

The required count scales as $1/p$, so halving the size of the difference doubles the vectors needed. A simulator can make $p$ larger by biasing vectors toward suspicious regions (as constrained-random verification does), but no fixed budget certifies equivalence. That is why SAT-sweeping engines use simulation only to propose candidate equivalences and then prove them with SAT.

## Expected time to the first counterexample

Draw vectors until one exposes the difference. The number drawn is geometric with success probability $p$, so its expectation is $1/p$. For $n = 10$ and $d = 64$, $p = 1/16$ and the expected number of vectors is 16.

Randomized algorithms come in two kinds. Monte Carlo algorithms have bounded running time and a probability of a wrong answer, like simulation-based equivalence with a fixed budget. Las Vegas algorithms are always right but have random running time, like a SAT solver with random restarts. Schoening's random walk for 3-SAT is a Monte Carlo algorithm that finds a satisfying assignment, if one exists, with high probability in $O((4/3)^n)$ time up to polynomial factors, better than the $2^n$ of brute force.

## Going deeper

Motwani and Raghavan, "Randomized Algorithms" (1995). Schoening, "A probabilistic algorithm for k-SAT and constraint satisfaction problems", FOCS (1999). Mishchenko, Chatterjee, Jiang and Brayton, "FRAIGs: a unifying representation for logic synthesis and verification", technical report, UC Berkeley (2005).
