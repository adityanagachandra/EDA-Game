# Timing prediction

Before routing, the placer and optimizer see only estimated wire parasitics, so their slack differs from signoff slack after routing. Pre-route timing prediction learns the post-route slack, or the net and arc delays, from placement-stage features. It is a regression problem, but the metric that matters is not the average error: it is whether the prediction gets the critical endpoints right.

## Mean absolute error

Given actual post-route slacks $y_i$ and predictions $\hat y_i$ at $n$ endpoints, the mean absolute error is $\frac{1}{n}\sum |y_i - \hat y_i|$. For actual slacks -10, 5 and 20 ps and predictions -4, -2 and 26 ps, the errors are 6, 7 and 6 ps and the MAE is 6.33 ps. MAE is in the units of the label and is robust to a few large errors, but it weights every endpoint equally. In a design with 100,000 endpoints most are far from critical, so a low MAE can coexist with large errors on the few hundred endpoints that decide closure.

## The coefficient of determination

$R^2 = 1 - \sum (y_i - \hat y_i)^2 / \sum (y_i - \bar y)^2$ compares the residual error with the spread of the labels. In the example, $\bar y = 5$, the total sum of squares is $225 + 0 + 225 = 450$ and the residual sum is $36 + 49 + 36 = 121$, so $R^2 = 1 - 121/450 = 0.731$. $R^2$ rewards a model for following a wide spread of slacks, which is easy when slacks range over hundreds of picoseconds. A model with $R^2 = 0.95$ over all endpoints can still be poor in the narrow band near zero slack, which is the band an optimizer acts on.

## Getting the violation status right

An optimizer's decision at an endpoint depends on the sign of slack: fix it or leave it. Counting endpoints where the predicted sign differs from the actual sign measures the errors that change behavior. In the example, endpoint 2 has actual slack 5 ps and predicted -2 ps, so one endpoint flips: the optimizer would waste effort on a passing path. The opposite flip, predicting a pass on a failing path, is worse, since nothing will fix it until signoff. Timing GNNs are therefore judged on critical-path ranking and violation classification as well as regression error, and their predictions guide optimization, while signoff STA remains the judge.

## Going deeper

Guo et al., "A Timing Engine Inspired Graph Neural Network Model for Pre-Routing Slack Prediction" (DAC 2022); Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), for the signoff analysis being predicted.
