# Ranking metrics for design-space exploration

Many EDA uses of a predictor only need it to rank candidates: which flow configurations to run fully, which placements to route, which sizings to simulate. A model can have a large absolute error and still rank perfectly, or a small error and rank poorly. Rank correlation and top-k overlap measure what such uses depend on.

## Concordant pairs

For $n$ candidates there are $\binom{n}{2}$ pairs. A pair is concordant when the prediction and the truth order it the same way, and discordant otherwise (ties aside). If true TNS magnitudes for five configurations are 50, 60, 30, 20 and 40 ps and predictions are 80.4, 70.7, 40.1, 50.9 and 60.2, then of the 10 pairs, 8 are concordant: only the pairs (1, 2) and (3, 4) are reversed. Counting concordance ignores the scale of predictions entirely, so a model biased by a constant, or by any increasing transformation, loses nothing on this measure.

## Kendall's tau

Kendall's tau is $(C - D)/\binom{n}{2}$, the concordant minus discordant pairs over all pairs. In the example $\tau = (8 - 2)/10 = 0.6$. It ranges from -1 (reversed order) to 1 (identical order), and 0 means no association. For selecting configurations, tau answers the right question better than MAE or $R^2$: whether running the predicted-best candidates is likely to run the truly best ones. Spearman's rank correlation, the Pearson correlation of the ranks, is a common alternative with the same invariance to monotone transformations. Computing tau directly costs $O(n^2)$ pair comparisons, which is trivial for the tens of configurations an exploration loop ranks.

## Top-k overlap

When only the predicted best $k$ will be run, the quantity that matters is how many of the truly best $k$ they contain. With $k = 2$, if the true best two are configurations 1 and 4 and the predicted best two are 2 and 4, the overlap is 1: half the expensive runs go to a configuration that is not among the best. Tau weights every pair equally, including pairs among poor candidates that no one will run; top-k overlap focuses on the head of the ranking. A model judged for exploration should be judged on its ranking metrics on held-out designs, since ranking quality also suffers under distribution shift.

## Going deeper

Kendall, "A New Measure of Rank Correlation," Biometrika (1938); Li et al., "Hyperband" (JMLR 2018), on ranking candidates under budgets.
