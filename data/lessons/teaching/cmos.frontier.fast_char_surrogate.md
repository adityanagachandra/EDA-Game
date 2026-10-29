# Fast characterization surrogates

Full characterization of a new corner or a new cell variant costs millions of SPICE transients. A fast surrogate, a regression or learned model that predicts tables from a few simulations and from related corners, can produce candidate tables in a fraction of the time. Used inside a fast-versus-accurate loop, it decides where SPICE must still run. Used carelessly, it writes confident wrong numbers into a library.

## What the fast engine trades away

A surrogate replaces simulation at most grid points with prediction: it interpolates between characterized voltages or temperatures, fits smooth delay surfaces from a sparse set of points, or uses a learned model of delay against device parameters, slew and load. It trades guaranteed accuracy for speed. That is acceptable for exploration, for example ranking library variants or estimating a new corner early, but signoff numbers must be traceable to accurate simulation within the agreed tolerance. A surrogate is only safe in an agent loop if it reports a calibrated uncertainty for each prediction, so the loop can decide which values to trust.

## Recursing to the accurate engine

The loop uses the uncertainty directly: accept a prediction whose standard deviation is below a threshold, and send every other point to SPICE. If a table's predictions carry standard deviations of 1.2, 3.5, 0.8, 4.1 and 2.0 ps and the threshold is 3 ps, two points are re-simulated and three accepted. The accurate results then refine the surrogate for the next table. The value of the surrogate is measured end to end: total simulations, including the re-simulations it triggers, needed to reach tables that pass the golden check. A surrogate that sends most points back to SPICE saves little, and one that never does is suspicious.

## Where surrogates fail

Surrogates fail most when extrapolating: a model trained on 0.7 to 0.9 V corners asked for 0.55 V meets near-threshold behaviour it has never seen. A process revision that shifts thresholds changes the training distribution, so the model must be recalibrated with fresh SPICE data. Constraint tables are harder than delay tables, because setup and hold come from searches over a sharp pass-fail boundary, where small model errors move the answer. Transition errors matter as much as delay errors because slew propagates. Most dangerous is miscalibration: small claimed uncertainty with large real error. Spot checks at table corners, high-curvature regions and low-confidence points are the defence, and a failed spot check must widen the recursion to SPICE.

## Going deeper

The agent-first EDA material in this curriculum (fast and accurate engines in a modeling loop); Rasmussen and Williams, "Gaussian Processes for Machine Learning" (2006) for calibrated predictive uncertainty.
