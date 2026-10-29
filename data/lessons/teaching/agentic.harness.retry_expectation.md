# Expected fix iterations in a harness loop

A harness gives a model several attempts at a change, each one informed by the failure of the last. How many attempts should it allow, and how long will a typical request take? With a per-attempt success probability, the questions become small exercises with the geometric distribution, and the answers set the harness's try budget and its expected turnaround.

## Success within a try budget

If each attempt succeeds independently with probability $p$, the chance that at least one of $k$ attempts succeeds is
$$P(\text{success within } k) = 1 - (1 - p)^k.$$
With $p = 0.35$ and a budget of 4 attempts, this is $1 - 0.65^4 = 1 - 0.1785 = 0.8215$. To reach 95% success the budget must satisfy $(1-p)^k \le 0.05$, so $k \ge \ln 0.05 / \ln 0.65 = 6.95$, that is 7 attempts. Doubling the budget does not double the success rate; each extra attempt adds less than the one before.

## Expected number of attempts

The number of attempts until the first success is geometric with mean $1/p$. With a budget of $k$, the harness stops at the first success or after $k$ attempts, so the expected number of attempts is
$$E[\min(G, k)] = \sum_{i=0}^{k-1} (1-p)^i = \frac{1 - (1-p)^k}{p}.$$
With $p = 0.4$ and $k = 5$: $(1 - 0.6^5)/0.4 = (1 - 0.07776)/0.4 = 2.3056$ attempts. If each attempt takes 90 s of model, build and test time, a request costs about 207 s on average. Failed requests use the whole budget, so the budget caps the worst case.

## When attempts are not independent

Real attempts are not independent coin flips. Good feedback raises the success probability on later attempts: a model that sees the exact compiler error is likelier to fix it. Then the chance of success within three attempts is $1 - (1-p_1)(1-p_2)(1-p_3)$; with $p_1 = 0.3$, $p_2 = 0.5$, $p_3 = 0.6$ it is $1 - 0.7 \times 0.5 \times 0.4 = 0.86$. The opposite also happens. When a model misunderstands the request, every attempt fails for the same reason, and extra attempts buy nothing. A harness should watch for repeated identical failures and stop early, because the independence assumption has broken.

## Going deeper

Any probability text on the geometric distribution, for example Ross, "A First Course in Probability".
