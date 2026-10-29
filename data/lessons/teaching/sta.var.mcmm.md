# Multi-corner multi-mode analysis

A chip has several modes, functional, test, low-power, each with its own clocks and constraints, and must work at several PVT and extraction corners. Each mode-corner pair is a scenario. Multi-corner multi-mode analysis times all relevant scenarios, merges their results per endpoint, and lets optimization see every scenario at once, so a fix for one does not quietly break another.

## Counting scenarios

A scenario pairs one mode's SDC with one corner's libraries and parasitics. With $m$ modes and $c$ corners there are $m c$ candidate scenarios, often too many to run. Teams prune pairs that can never dominate: a test mode that runs at a low scan frequency need not be checked for setup at every corner. With 4 modes, 9 corners and 11 pruned pairs, $4 \times 9 - 11 = 25$ scenarios remain. Pruning is a judgment that must be justified, because a pruned scenario is never checked. Some flows also mark scenarios as setup-only or hold-only, which reduces work further without dropping checks that matter.

## Merging an endpoint's slack

An endpoint passes only if it passes in every scenario, so its merged slack is the minimum over scenarios. If endpoint E1 has slacks 26, 36 and 6 ps in three scenarios, its merged slack is 6 ps, set by the third scenario. The merged view shows which scenario limits each endpoint, and that is what an optimizer must improve: speeding up E1's path helps only if it raises the slack in the limiting scenario without lowering another below it. Timers keep per-scenario values so that a fix can be evaluated in all of them; MCMM-aware optimizers evaluate each candidate change across the active scenarios before accepting it.

## Merged TNS

The design's merged TNS sums each endpoint's merged slack, counting only the negative ones. With merged slacks of -38, -30 and -22 ps across three endpoints, merged TNS is -90 ps. Adding the per-scenario TNS values instead would count an endpoint once for every scenario in which it fails, overstating the problem and rewarding fixes that help only an already-dominated scenario. Merged WNS is the minimum merged slack. These merged metrics are what closure dashboards track. A rise in merged TNS after a change, even with per-scenario improvements, usually means the fix shifted the violation to another scenario.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on timing closure; the OpenROAD documentation on multi-corner analysis.
