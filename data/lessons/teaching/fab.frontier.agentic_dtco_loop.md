# Agent-driven DTCO loops and their rewards

DTCO is a search over coupled design and technology knobs, scored by expensive flows. That makes it a natural target for agents: propose a candidate, pick fast or full engines to evaluate it, read the results, choose the next one. The engineering challenge is the reward. A score built from a cheap proxy will be gamed, and a technology choice made on a gamed score is very expensive to undo.

## The agent proposes, the engines compute

A candidate is a vector of knobs: cell height, fins or sheets per device, CPP, metal pitches, rail scheme and rule values. The agent's job is selection: which candidate next, which benchmark blocks, fast or full flow. The engines do the computation: cell generation, characterization, synthesis, placement, routing, extraction and timing. Fast flows rank many candidates cheaply; full flows confirm the few that matter. A proposal such as a 5-track library cannot be scored until the loop has generated its cells, characterized them and routed benchmark blocks, because cell height alone does not say whether blocks route.

## Designing a tool-grounded reward

The reward should be the real objective under hard constraints: routed block PPA at a fixed timing target, with DRC-clean routing and routability as constraints, plus cost per good die including mask count. Leaving cost out invites options that win PPA with expensive process steps. Scoring on one benchmark block invites overfitting to that block, so the reward averages several diverse blocks. A candidate whose fast flow ran with DRC disabled is unverified and earns nothing until a DRC-clean full flow confirms it. Rule values may move only if they are declared DTCO variables with ranges the process team approved.

## Reward hacking and hand-off

Typical exploits: very small cells that score well because the fast placer ignores pin access; layouts that sit just outside the litho model's validity range; rule values nudged until DRC passes. Each is closed by adding the missing check to the reward and by auditing the top candidates with full flows. The final output goes to the human DTCO team as a ranked set with PPA, cost, the engines used and the audits run, because technology direction is a human decision. A policy learned on one process encodes that process's trade-offs, so moving it to a new process is a hypothesis to revalidate, not a result.

## Going deeper

The agent-first EDA anchors on tool-grounded reinforcement learning and fast and accurate engines; Clark et al., "ASAP7: A 7-nm finFET predictive process design kit" (2016), as an open basis for DTCO experiments.
