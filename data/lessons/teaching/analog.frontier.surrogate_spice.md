# Fast surrogates versus SPICE in the loop

SPICE is accurate and slow; a learned or analytical surrogate predicts performances in milliseconds. An agentic or optimization loop uses the surrogate to search and SPICE to decide. The questions for a tool designer are quantitative: how wrong is the surrogate, how much time does it save, and how often does it send a failing design forward as passing.

## Surrogate error

A surrogate maps a sizing to predicted performances: gain, GBW, phase margin, power. Its error on one design is the absolute percentage error against SPICE, $|p - a|/a \times 100$. A GBW prediction of 52 MHz against SPICE's 48 MHz is 8.3% off. Errors are rarely uniform: surrogates trained on typical sizings are worst exactly where an optimizer pushes, at the edges of the design space and near region boundaries where devices leave saturation. Error must therefore be measured on held-out designs, ideally from the region the search is visiting now, and re-measured as the search moves. A surrogate's accuracy on its training distribution says little about its accuracy on the agent's proposals.

## The speed-up budget

Screening with a surrogate and verifying the best candidates with SPICE costs $n\,t_s + k\,c\,t_S$ instead of $n\,c\,t_S$, for $n$ candidates, $k$ verified, $c$ corners, surrogate time $t_s$ and SPICE time $t_S$. With 2000 candidates, 5 corners, 30 s per SPICE run, 0.01 s per surrogate call and the top 20 verified: all-SPICE costs 300,000 s; the mixed flow costs $20 + 3000 = 3020$ s, a speed-up of 99. The verification term dominates, so the useful lever is how few candidates need SPICE, which depends on how well the surrogate ranks them. This is the fast-accurate engine pattern: the fast engine prunes, the accurate engine decides, and the loop recurses into the accurate engine when the fast one is near a decision boundary.

## False accepts and guard bands

The dangerous error is a false accept: the surrogate says a design passes and SPICE says it fails. If the surrogate is the judge, those designs ship; if SPICE is the judge, they waste verification runs. A guard band reduces them: forward a candidate only when the prediction clears the spec by a margin, for example phase margin at least 63 degrees for a 60-degree spec. Counting false accepts on a validation set for several guard bands gives the trade-off curve: a wider band means fewer false accepts but more false rejects, good designs never verified. The guard band should track the surrogate's measured error in the region. However good the surrogate, final acceptance stays with SPICE, because the surrogate's error is unknown on exactly the designs it has never seen.

## Going deeper

The agent_first_eda lesson on fast and accurate engines in a modeling loop.
