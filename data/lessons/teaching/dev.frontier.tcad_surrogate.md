# ML surrogates for TCAD in DTCO loops

Design-technology co-optimization explores many device options, but each 3D TCAD point is slow. A machine-learned surrogate trained on TCAD runs can predict device characteristics in milliseconds, which turns exploration into a fast inner loop. The surrogate is a fast engine; TCAD remains the accurate one. Deciding when to trust the surrogate and when to recurse to TCAD is the core of using it well.

## Surrogates as the fast engine

A TCAD surrogate learns a map from device design inputs, such as gate length, fin or sheet dimensions, spacer thickness and doping, to outputs such as I-V curves on a fixed bias grid, threshold, swing, DIBL, capacitances or compact-model parameters. Training data comes from TCAD runs over a designed set of inputs.

The value is speed. If a TCAD point takes minutes to hours and a surrogate takes milliseconds, a DTCO loop can sweep thousands of options, run optimizers and study variability, then re-simulate only the promising candidates in TCAD before decisions are made. This is the fast-versus-accurate engine pattern: the fast engine explores, the accurate engine confirms.

Output representation matters for integration. A surrogate that emits compact-model parameters or I-V curves on the grid that extraction already consumes drops into existing flows; one that emits raw mesh values or a pass-fail flag does not.

## When to recurse to TCAD

Optimizers seek the corners of a model. If a DTCO optimizer drives the surrogate to a geometry outside its training range and reports a record drive current, that result is untrusted extrapolation. The loop must re-simulate the point in TCAD before believing it, and the confirmed point can be added to the training set.

A surrogate should therefore report uncertainty, for example from an ensemble of models, and flag inputs far from its training data. High uncertainty is the trigger to recurse; this is also how active learning chooses where to spend new TCAD runs.

A surrogate is only as current as its training data. When the TCAD deck is recalibrated to new silicon, the surrogate is stale and must be checked against new runs and retrained; versioning the surrogate with the deck it learned keeps this visible.

Validation should match the decision: error on held-out TCAD runs in the region the decision depends on, and agreement with TCAD on the ranking of candidates, since a decision depends on ordering.

## Physics, rewards and sampling

A surrogate that sometimes predicts current falling as gate voltage rises violates physics and breaks downstream solvers and extraction. Constraints such as monotonicity in $V_{GS}$, positivity and smoothness should be built into the surrogate's form or training loss, or the surrogate should predict quantities from which a physical model is constructed.

When a reinforcement-learning agent or optimizer uses the surrogate as its reward, it will exploit surrogate errors: reward hacking. The reward rises while true TCAD performance does not. Periodic TCAD scoring of the best candidates keeps the reward honest.

The surrogate inherits every error of the TCAD deck, so its accuracy is bounded by how well that deck is calibrated to silicon. It approximates TCAD, not reality.

With a fixed TCAD budget, training data should first cover the input space with a space-filling design, to avoid blind spots, and then concentrate near decision boundaries or optima through active learning, where accuracy pays.

## Going deeper

The repository's docs/eda/agent-first-concepts.md, sections 8 and 9; Forrester, Sóbester and Keane, "Engineering Design via Surrogate Modelling" (2008).
