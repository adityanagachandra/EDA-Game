# Fast versus accurate QoR in the loop

An agent exploring synthesis recipes cannot run full mapping, placement and timing for every candidate. It scores most candidates with fast proxies, AIG node count, levels or a quick mapping, and spends accurate runs only where decisions are close. Knowing how proxies disagree with signoff metrics, and when to recurse to the accurate engine, is the core skill of a design-space exploration loop.

## Measuring proxy agreement

Compare the ranking a proxy gives with the ranking of the accurate metric. For four recipes, count the pairs ordered oppositely, the discordant pairs. With AIG node counts 963, 1046, 1115, 914 and mapped areas 480, 402, 433, 453, the pairs (1, 2), (1, 3), (2, 4) and (3, 4) are discordant: 4 of 6. Kendall's tau is (concordant minus discordant) over the number of pairs, here $(2 - 4)/6 = -1/3$. A proxy that disagrees this often should not choose the winner. Rank agreement over all pairs can be high while the top few are still misordered, which is exactly where the decision is made.

## Pruning with fast engines

A typical loop scores hundreds of candidates with fast metrics: AIG nodes for area, levels or a unit-delay LUT mapping for depth. It keeps a shortlist and runs full mapping and static timing only on that shortlist. The engine must expose cheap observation hooks after every pass, so the agent can read node count, level count and a quick mapped estimate without exporting the design. Learned surrogate estimators trained on mapped results can sharpen the fast scores, but an optimizer drives designs out of the training distribution, where the surrogate is least reliable, so periodic accurate checks are essential.

## When to recurse to accuracy

Escalate to the accurate engine when candidates differ by less than the proxy's error, when proxy and accurate rankings have disagreed on similar designs, or when the fast model omits something that matters here. Unit-delay depth ignores fanout load and wire delay, so a recipe that looks fast can be slow after real timing. Metric variance across seeds sets a noise floor: differences below it should not drive decisions. Reward functions built only on fast proxies invite reward hacking, so a loop's final choice always comes from accurate signoff-relevant metrics.

## Going deeper

The fast and accurate engine loop in this curriculum's agentic domain; Kendall, "A new measure of rank correlation", Biometrika (1938).
