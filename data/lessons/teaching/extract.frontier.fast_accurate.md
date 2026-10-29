# Fast and accurate extraction in an agent loop

An agent that optimizes a layout cannot wait for signoff extraction after every move. It runs a fast extraction mode in its inner loop and recurses to the accurate engine only where the fast answer cannot be trusted. Deciding where that is requires an error bound per net, a rule that turns bounds and slack into a re-extraction list, and an accounting of what the loop saves.

## Error bounds on fast results

A fast mode is characterized by its error against the accurate engine, for example within 10% on net capacitance across a validation set. If a net's interconnect-driven delay scales linearly with its capacitance, its delay error bound is that fraction of its delay: a 247 ps net carries a 24.7 ps bound. The bound is the honest uncertainty of the fast result, and an agent should carry it with every fast-mode slack. Bounds are better per pattern class or per net type than global: a fast mode may be within 2% on long parallel runs and 15% on dense via regions. An observation hook that reports the bound with the parasitics is what makes the fast mode usable by an agent at all.

## Which nets to recurse on

A net's fast-mode slack is ambiguous when its magnitude is within its error bound, because the true slack could have either sign. Those nets go to accurate extraction; clearly passing or clearly failing nets keep their fast values. With a 4% error and six nets of slack and delay $(-3, 30)$, $(30, 210)$, $(-10, 290)$, $(-12, 20)$, $(0, 70)$ and $(-1, 110)$ ps, the bounds are 1.2, 8.4, 11.6, 0.8, 2.8 and 4.4 ps, and the ambiguous nets are the third, fifth and sixth. This is the recursion rule of the fast-accurate loop. It is conservative in the right direction: a net is accepted on fast data only when even the worst fast-mode error could not change its verdict.

## What the loop saves

Cost decides how tight the error bound must be. With 200000 nets, fast extraction at 0.2 ms per net and accurate at 10 ms, all-accurate takes 2000 s. Fast on everything plus accurate on 5% takes $40 + 100 = 140$ s, a 14.3 times speed-up. The accurate fraction dominates: halving it to 2.5% gives $40 + 50 = 90$ s. A tighter fast mode shrinks the ambiguous fraction, so investment in fast-mode accuracy, for example better calibration, pays directly in loop speed. The risk is the reverse: an error bound that is too optimistic lets wrong verdicts through, and the agent closes paths that reopen at signoff. Final signoff always uses the accurate engine on the whole design.

## Going deeper

The OpenROAD documentation on its estimated versus extracted parasitics modes.
