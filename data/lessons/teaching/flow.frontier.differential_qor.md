# Differential QoR checks on a flow change

When an agent changes a flow, by a parameter, a plug-in or a script, someone must decide whether the change did what it claimed and nothing else. A differential QoR check runs the flow on reference designs before and after the change and compares results endpoint by endpoint and metric by metric. It is the evidence that lets a human approve an automated change in minutes rather than re-deriving it.

## Comparing two runs

A differential check pairs every endpoint, net or region in the baseline run with its counterpart in the new run and computes the change. Small numerical noise is expected, so the check uses a tolerance: an endpoint whose slack changed by more than the tolerance has moved. With baseline slacks $12, -5, 30, 8, -2$ ps, new slacks $12, 4, 29, 8, -14$ ps and a tolerance of 2 ps, two endpoints moved: the second improved by 9 ps and the fifth worsened by 12 ps. The tolerance should come from measured run-to-run noise on the reference designs, not from a guess, or the check will either flag noise or miss real changes.

## Summary metrics and their limits

Summary numbers are computed alongside. Total negative slack sums the negative slacks: -7 ps in the baseline and -14 ps after, a worsening of 7 ps. Worst negative slack moved from -5 to -14 ps. Summaries alone can hide trouble: a change can leave WNS unchanged while moving hundreds of endpoints, or improve TNS while breaking one critical path. That is why the per-endpoint differential is the primary evidence and summaries are secondary. Summary numbers are still useful as a first screen: if WNS, TNS and the count of failing endpoints are all unchanged within noise, the per-endpoint differential is usually short to review.

## Intended versus unintended movement

A change request names what should move: a set of paths, a block, a layer. The check then reports two sets: intended endpoints that moved as requested, and endpoints outside the intended set that moved at all. A change to derating on high-threshold data paths should move only paths through those cells. If 3 of the moved endpoints are outside the intended set, the change has side effects that must be explained before approval. The same idea applies beyond timing: area per block, congestion per region and DRC counts per layer each get their own differential with their own tolerance.

## Going deeper

The agent-first EDA notes in this curriculum, sections on golden and differential testing and on the derating seam.
