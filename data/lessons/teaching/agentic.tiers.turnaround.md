# Turnaround and risk by change tier

The tiers differ in cost as well as in risk. A configuration change can be applied and checked in minutes; a plug-in needs a build, golden tests and review; a core edit needs full regression and an engine owner's time. The mix of tiers in the request stream therefore decides how fast a deployment responds to its users, and the arithmetic shows why moving a seam is worth an engineer's week.

## Expected turnaround of a request mix

If a fraction $f_j$ of requests land in tier $j$ and a tier-$j$ request takes $t_j$ hours from request to approved change, the expected turnaround is
$$E[T] = f_1 t_1 + f_2 t_2 + f_3 t_3.$$
With 70% Tier 1 at 0.5 h, 25% Tier 2 at 4 h and 5% Tier 3 at 60 h (a full regression and a review queue), $E[T] = 0.35 + 1 + 3 = 4.35$ h. The median request is fast, but the mean is dominated by the rare core edits.

## Where the time goes

The share of total change time spent on tier $j$ is $f_j t_j / E[T]$. In the example, Tier 3 accounts for $3 / 4.35 = 69\%$ of all the time spent on changes although it is only 5% of requests, and Tier 1 accounts for 8%. Risk follows the same pattern: a Tier 3 edit changes behavior for every design and every customer that uses the engine, while a Tier 1 change affects one configuration. Measuring this split on a real request log is how an engine team finds where its effort goes.

## The payoff of moving a seam

Suppose analysis of the Tier 3 requests shows that 80% of them were all one kind, say new derating rules that needed context the derate seam did not carry. One reviewed core edit extends the seam's context, and from then on those requests are Tier 2. The fraction moved is $0.05 \times 0.8 = 0.04$, and the new mix is 70% at 0.5 h, 29% at 4 h and 1% at 60 h, so $E[T] = 0.35 + 1.16 + 0.6 = 2.11$ h, less than half the old turnaround. This is the quantitative form of the rule that a seam is in the right place when core edits for that kind of change essentially never happen.

## Going deeper

The agent-first EDA notes in this curriculum, section on change tiers.
