# The gap between proxy reward and signoff QoR

Agents train and search on fast proxies and are judged on signoff quality. The gap between the two is not a vague worry; it can be measured on paired data. Correlation says how well a proxy tracks signoff overall, regret says what choosing by the proxy costs, and shortlist recall says whether using the proxy to filter candidates before accurate evaluation loses the best ones.

## Correlation between proxy and signoff

For candidates evaluated by both the proxy and signoff, the Pearson correlation of their improvements measures linear agreement:
$$r = \frac{\sum (x_i - \bar x)(y_i - \bar y)}{\sqrt{\sum (x_i - \bar x)^2 \sum (y_i - \bar y)^2}}.$$
With proxy improvements $1, 2, 3, 4$ and signoff improvements $2, 1, 4, 3$, the deviations give $S_{xy} = 3$ and $S_{xx} = S_{yy} = 5$, so $r = 0.6$. A correlation of 0.6 is enough for coarse screening but not for picking a winner among close candidates. Correlations should be measured on the population the agent actually searches, since a proxy can track well across very different designs and poorly among the near-optimal candidates that matter.

## Regret of choosing by proxy

The cost of trusting the proxy is the regret: the signoff quality of the best candidate minus that of the candidate the proxy ranks first. If five candidates have proxy improvements $4.0, 3.1, 2.5, 1.2, 0.5$ and signoff improvements $1.0, 2.8, 2.0, 0.4, 0.6$, the proxy picks the first, which delivers 1.0 at signoff, while the second delivers 2.8. The regret is 1.8. A reward trained to maximize the proxy pushes towards exactly the outputs where proxy and signoff disagree most, so regret tends to grow as optimization pressure increases.

## Shortlisting with a proxy

The practical compromise is to shortlist by proxy and decide by signoff: evaluate the top $m$ proxy candidates accurately and choose the best of those. The question is whether the shortlist contains the true best. Recall of the top $k$ is the fraction of the true top $k$, ranked by signoff, that lands in the proxy's top $m$. If among eight candidates the true top two are proxy ranks 2 and 5 and the shortlist is the proxy top three, recall is 1/2. Growing $m$ raises recall at the cost of more accurate runs, which is the same budget trade as recursion in the fast modeling loop.

## Going deeper

Goodhart, "Problems of Monetary Management: The U.K. Experience" (1975), the origin of Goodhart's law; any statistics text on correlation and rank statistics.
