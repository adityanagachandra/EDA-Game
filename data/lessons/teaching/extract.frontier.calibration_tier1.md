# Calibration as a Tier 1 change

Calibration is the extractor's most frequent change, and it should never require a recompile. If correction factors per pattern class live in a configuration table, a new process revision, a new fill style or new golden data becomes a Tier 1 change: fit the coefficients, validate them, promote the table. An agent can run that whole loop, provided the validation is strict enough that a bad fit cannot slip through.

## Fitting the table entry

For one pattern class, collect samples of fast-mode values $x_i$ and golden 3D-solver values $g_i$. The least-squares scale factor is $k = \sum x_i g_i / \sum x_i^2$. For four samples with $\sum x_i g_i = 124.54$ and $\sum x_i^2 = 130.82$, $k = 0.952$: fast mode reads 5% high for this class. The new $k$ is written into the configuration table for that class. Nothing in the core changes: the fracturing, the lookup and the netlist assembly are the same code, and the only artifact that changes is data. That is what makes the change Tier 1, and what lets it be versioned, reviewed as a data diff and rolled back by restoring the previous table.

## Applying the calibrated table

After the update, fast mode multiplies each extracted value of the class by its factor. With $k = 1.0215$ and an uncalibrated value of 4.7 fF, the reported value is 4.801 fF. The factor applies per class, so a net that crosses several classes gets each segment corrected by its own factor before summing. Agents and flows must record which table version produced each parasitic file, since two runs with different tables are not comparable. This matters for differential checks: a calibration change moves every net of the class, and a differential report must attribute that movement to the table change, not to an ECO run at the same time. Calibration updates and design changes should therefore never be mixed in one step.

## Validation before promotion

A fitted factor must be validated on held-out nets that were not used in the fit. Scale each held-out fast value by the new $k$, compare with golden, and count nets outside the tolerance. If six held-out nets give four beyond a 2% tolerance after scaling by 1.0068, the fit has not captured the class: the scatter is not a uniform bias, and the right response is a new or finer pattern, a Tier 2 change, not a further tweak of the factor. Promotion requires the count to be zero, or within a stated budget. This gate is what keeps a Tier 1 change safe in an agent's hands: the agent can propose and fit freely, but only a validated table reaches signoff.

## Going deeper

Kao, Lo, Basel and Singh, "Parasitic extraction: current state of the art and future trends", Proceedings of the IEEE (2001), on calibration of pattern-based extraction.
