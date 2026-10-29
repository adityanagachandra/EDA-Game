# Robust and stochastic optimization

Designs must work across process corners, voltages, temperatures and operating modes, and their parameters vary randomly. Robust optimization protects the worst case over a set of scenarios; stochastic optimization minimizes the expected cost under a probability model. Multi-corner multi-mode timing closure is robust optimization, yield-aware sizing and statistical timing are stochastic, and the difference between the two answers is a measurable price paid for the guarantee.

## Min-max over corners

A robust formulation picks the decision whose worst-case cost over the scenarios is smallest: $\min_d \max_s c(d, s)$. With three sizing options evaluated at SS, TT and FF: option 1 costs 6, 11, 18; option 2 costs 2, 11, 8; option 3 costs 14, 10, 9. The worst cases are 18, 11 and 14, so the robust choice is option 2 with worst-case cost 11.

The probabilities of the corners play no role. This is the right model when every scenario must be met, as in signoff, where a path must close at every corner and mode. In continuous problems the min-max is written in epigraph form: minimize t subject to $t \ge c(x, s)$ for every scenario, which keeps an LP or convex problem in its class, with one constraint per scenario.

## Minimizing expected cost

A stochastic formulation minimizes $\mathbb{E}[c(d, S)] = \sum_s p_s c(d, s)$. With probabilities 0.1, 0.8, 0.1 and options 1: (6, 7, 2), 2: (9, 11, 14), 3: (20, 10, 12), 4: (18, 7, 3), the expected costs are $0.6 + 5.6 + 0.2 = 6.4$, $0.9 + 8.8 + 1.4 = 11.1$, $2 + 8 + 1.2 = 11.2$ and $1.8 + 5.6 + 0.3 = 7.7$. The best is option 1 at 6.4.

With continuous random parameters the expectation becomes an integral, approximated by sampling: sample average approximation optimizes the average over N Monte Carlo samples. Chance constraints, such as requiring timing to be met with probability at least 0.99 (a yield target), are the stochastic counterpart of worst-case constraints; under Gaussian models they become deterministic constraints on mean plus a multiple of the standard deviation, which is the form statistical timing and POCV margins take.

## The price of robustness

The robust choice is usually not the expected-cost optimum. The difference between the robust choice's expected cost and the minimum expected cost is the price of robustness, what the worst-case guarantee costs on average. With probabilities 0.25, 0.5, 0.25 and options 1: (19, 16, 6), 2: (7, 4, 11), 3: (18, 10, 4): worst cases 19, 11, 18 make option 2 robust; expected costs 14.25, 6.5, 10.5 make option 2 best on average too, so the price is 0.

When the two differ, the designer must decide whether rare scenarios justify the average loss. Between the extremes sit distributionally robust optimization (worst case over a set of plausible distributions) and conditional value-at-risk, which averages the worst fraction of outcomes. Over-margining every corner, as flat worst-case derates do, is the EDA instance of paying a large price of robustness, which is why statistical and location-aware margins were introduced.

## Going deeper

Ben-Tal, El Ghaoui and Nemirovski, "Robust Optimization" (2009); Shapiro, Dentcheva and Ruszczynski, "Lectures on Stochastic Programming" (2nd ed., 2014); Bertsimas and Sim, "The price of robustness" (Operations Research, 2004).
