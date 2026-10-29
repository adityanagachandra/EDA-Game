# Surrogates and active learning

Labels for an EDA surrogate cost tool runs, so which points to label matters as much as how many. Active learning lets the model choose: it labels the candidates it is least sure about, where a label teaches it most. Uncertainty from an ensemble or a Gaussian process drives the choice, and retraining overhead decides whether the saving is real.

## Ensemble disagreement

An ensemble of models trained on different subsets or seeds gives several predictions for one candidate. Their spread estimates the model's uncertainty there. If four surrogates predict a sizing's phase margin as 60, 64, 58 and 62 degrees, the mean is 61 and the population variance is $(1 + 9 + 9 + 1)/4 = 5$ square degrees. Disagreement is high where training data is sparse or the function changes fast, and low where the models have seen similar points. It measures the model's own uncertainty, not noise in the simulator, and it can be low and wrong together if every member shares the same bias.

## Uncertainty sampling

The simplest active-learning rule, uncertainty sampling, sends the unlabeled candidate with the largest predictive variance to the expensive tool. Among five candidate paths, if three ensemble members agree closely on four of them and spread over 11 ps on the fifth, the fifth goes to the timer next, whatever its predicted slack. Pure uncertainty sampling explores; it does not care whether a point is relevant to the design goal. Variants weight uncertainty by relevance, for example preferring uncertain points near zero slack or near a specification boundary, and batch versions pick several diverse points per round so the tool runs in parallel.

## Counting the saving

Suppose random sampling needs 2000 labeled runs to reach a target error and active learning needs 600. At 4 machine-hours per label, labels saved are worth $1400 \times 4 = 5600$ hours. Active learning retrains after each batch; with batches of 50 and 2 hours per retrain it runs $\lceil 600/50 \rceil = 12$ rounds at 24 hours, so the net saving is 5576 hours. Small batches choose better points but retrain more often and serialize the tool runs. The labels it collects are also biased toward hard regions, so a model evaluated on them looks worse than it is, and a separate random test set is needed to measure accuracy.

## Going deeper

Settles, "Active Learning Literature Survey" (University of Wisconsin-Madison, 2009); Rasmussen and Williams, "Gaussian Processes for Machine Learning" (2006), chapter 2.
