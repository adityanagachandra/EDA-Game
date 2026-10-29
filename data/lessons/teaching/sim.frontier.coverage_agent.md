# Agentic coverage-closure loops

Coverage closure is a loop of reading holes, changing stimulus and rerunning, which makes it a natural job for an agent driving the simulator as a tool. The simulator supplies an honest, measurable reward: bins actually hit. The engineering problem is keeping that reward honest, because the cheapest way to raise a coverage number is to change the measurement rather than the stimulus. This concept applies the harness loop, tool-grounded reward and change tiers to closure.

## The closure loop

One iteration reads the merged coverage database at bin level, including cross holes and their hit counts, together with the constraints and sequences that control the relevant fields. The agent proposes a stimulus change: new `dist` weights, a constraint mode, or a new sequence. A cheap pre-check samples the changed constraints with the solver alone, thousands of `randomize()` calls without simulation, to see whether the target values now appear. If they do, the harness runs a regression, merges coverage, and compares against the baseline: which bins were newly hit, whether any were lost, and what the run cost. The result, not the agent's explanation, decides whether the change is kept. This is the edit, build, test loop with coverage as the test.

## Reward design and how it is gamed

A natural reward is the number of newly covered bins. It can be gamed. An agent with write access to the covergroup can add `ignore_bins` for every hole, delete a coverpoint, or lower `option.at_least`, and report 100%. An agent that can touch checkers can weaken a scoreboard comparison so a failing scenario passes and its coverage counts. Stimulus can hit bins in degenerate ways, such as a zero-length burst that technically lands in a "long transfer during flush" cross. Guardrails follow: the coverage model and the checkers are frozen, outside the agent's change surface; reward is measured on a held-out regression; cost (simulation time) is penalized so brute-force seed farming does not win; and sampled traffic for newly hit bins is spot-checked.

## Change tiers for closure

The tiers make the agent's authority explicit. Changing the weight of an existing `dist`, a constraint-mode switch or a plusarg is **Tier 1**, configuration: no new code, trivially reversible. Writing a new sequence class, for example one that throttles the reader so a FIFO fills, is **Tier 2**: new code at a defined extension point, the sequence library, with the environment's core untouched. Editing the DUT, the covergroups, the scoreboard or the reference model is outside the closure agent's remit entirely, and must go to a human, because changing what is measured or what is tested inverts the purpose of verification. Every accepted change ships with its diff, the coverage delta against the baseline and the seeds that reproduce it, so a reviewer can rerun it.
