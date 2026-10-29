# Branch penalties and prediction

Branches are about one instruction in five, and every mispredicted branch flushes the work fetched behind it. Branch predictors guess the outcome early, and the simplest dynamic ones are tiny finite-state machines per branch. Their accuracy, multiplied by the pipeline's flush penalty, is a large term in CPI, and predictor state machines are a neat place where FSM design, microarchitecture and performance arithmetic meet.

## The CPI cost of branches

If a fraction $f_b$ of instructions are branches, a fraction $m$ of those are mispredicted, and each misprediction costs $P$ cycles of flushed work, then

$$CPI = CPI_{base} + f_b \cdot m \cdot P.$$

With $CPI_{base} = 1$, $f_b = 0.2$, $m = 0.1$ and $P = 3$: $CPI = 1.06$. Deep pipelines resolve branches late, so $P$ can be 15 or more cycles in high-frequency cores, and the same 10% misprediction rate then costs $0.2 \cdot 0.1 \cdot 15 = 0.3$ cycles per instruction. That is why predictor accuracy improved from roughly 90% to well above 95% as pipelines deepened.

## The two-bit saturating counter

A 2-bit predictor is a four-state counter per branch: 0 and 1 predict not taken, 2 and 3 predict taken. A taken outcome increments (saturating at 3), not taken decrements (saturating at 0). Two wrong outcomes in a row are needed to flip a strong prediction.

Start at state 1 (weakly not taken) and feed outcomes T, T, N, T, T, T, N, T. The predictions are N, T, T, T, T, T, T, T, and the states after each outcome are 2, 3, 2, 3, 3, 3, 2, 3. Mispredictions occur on the first T (predicted N) and the two N outcomes: 3 of 8. The counter is a Moore machine whose output is its prediction.

## One bit versus two bits on loops

A 1-bit predictor simply repeats the last outcome. On a loop branch taken $k - 1$ times then not taken once, repeated, it mispredicts twice per loop execution after the first: once at the exit (predicted taken) and once at the next entry (it remembered not taken). The 2-bit counter mispredicts only at the exit, since one not-taken outcome moves it from 3 to 2, still predicting taken.

For a loop of 4 iterations executed 3 times (outcome pattern TTTN repeated), starting from not taken, a 1-bit predictor mispredicts 6 times (the first T, then each of the first two exits and the re-entry after it, then the final exit: 1 + 2 + 2 + 1), and a 2-bit counter starting at state 1 mispredicts 4 times. Modern predictors (gshare, TAGE) index counters by global history to capture correlations, but the saturating counter remains their basic cell.

## Going deeper

Smith, "A Study of Branch Prediction Strategies" (ISCA, 1981); McFarling, "Combining Branch Predictors" (DEC WRL Technical Note TN-36, 1993); Seznec and Michaud, "A Case for (Partially) Tagged Geometric History Length Branch Prediction" (JILP, 2006).
