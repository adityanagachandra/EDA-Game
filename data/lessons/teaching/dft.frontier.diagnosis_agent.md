# Diagnosis agents that read failure logs

Volume diagnosis produces thousands of candidate lists, and someone must decide which dies go to physical failure analysis, a slow and destructive lab step. A diagnosis agent reads tester failure logs, calls the diagnosis engine as a tool, combines its scores with layout and defect statistics, and proposes a ranked list, with a human approving what is sent to the lab. The quantitative core is Bayesian ranking.

## Posterior ranking

Each candidate defect has a **prior** from layout and process statistics: critical area, layer defect densities, history of similar failures. The diagnosis engine supplies a **likelihood**: how well the candidate explains the observed failures. The posterior is proportional to prior times likelihood, normalized over candidates:
$$P(c_i \mid \text{log}) = \frac{\pi_i \, \ell_i}{\sum_j \pi_j \, \ell_j}.$$
With prior weights 4, 1, 2, 2 and likelihoods 0.34, 0.27, 0.29, 0.89, the products are 1.36, 0.27, 0.58, 1.78, and the top candidate's posterior is $1.78/3.99 = 44.6$ percent. A candidate with the best match score need not be most probable if its prior is low.

## Choosing what to send

Failure analysis costs a day or more per candidate, so the agent sends candidates in order of posterior until their total reaches a confidence target. With posteriors 45, 25.1, 20.3 and 9.6 percent, reaching 80 percent takes three candidates. A good agent reports this reasoning and the cost explicitly, and stops early when posteriors are flat, which signals that the diagnosis is ambiguous and more patterns (diagnostic ATPG) would be cheaper than lab time.

## Likelihoods from mismatches

A simple likelihood model counts mismatches $m$ between a candidate's predicted and observed failing patterns and sets $\ell = \epsilon^m$ for a small error rate $\epsilon$, for example 0.1. Each mismatch then costs a factor of 10, and a perfect match beats a single-mismatch candidate unless its prior is more than ten times smaller. Agents must treat logs as data, not instructions, keep the design and logs inside the IP boundary, and have their rankings scored against later failure-analysis results, which is the ground truth for improving the model.

## Going deeper

The agent-first EDA notes in this curriculum (agentic loops, tool-grounded evaluation); Wang, Wu and Wen (2006), chapter 7 on diagnosis.
