# Multi-agent planners and shared design state

A full design flow is too broad for one agent. The natural structure mirrors a design team: specialized agents for RTL, verification, implementation and signoff, coordinated by a planner. Making that work depends less on the agents' cleverness than on how they share state, hand off work and avoid stepping on each other, which is a systems problem familiar from distributed computing.

## Specialized agents and a planner

Each specialist has its own tools, skills and context: a verification agent knows the testbench, coverage and formal tools; an implementation agent knows the placer, router and their skills; a signoff agent knows timing, IR and DRC reports. A planner decomposes a goal ("close timing on block X at 1.2 GHz without increasing area by more than 2%") into tasks, assigns them, tracks dependencies and decides what to do with results. Specialization keeps each agent's context small and its tool set narrow, which makes it more reliable and limits what any one agent can break.

## Shared state in the design database

Agents should share state through the design database, not by passing designs in messages. Each agent reads and writes versioned design states by handle: the implementation agent produces placement state 41 from synthesis state 17; the signoff agent analyzes state 41 and writes a report linked to it. Messages carry handles, task descriptions and summaries. Versioning gives every result a precise lineage, so a timing report can never be confused with one from an older placement, and it lets the planner roll back to a known-good state. Concurrent edits to the same state need locking or a merge step, as in any database.

## Coordination failures

Multi-agent systems fail in characteristic ways. **Stale reads**: the signoff agent reports violations on a placement that the implementation agent has already replaced. **Conflicting edits**: a timing fix upsizes a driver while a power fix downsizes it. **Duplicated work**: two agents rerun the same extraction. **Cascading errors**: one agent's wrong summary misleads every agent downstream. The remedies are the ones distributed systems use: version every state and check versions on write, give each object an owning agent per phase, cache results by input version, and have the planner verify claims by reading engine output rather than trusting agents' summaries.

## Going deeper

Kleppmann, "Designing Data-Intensive Applications" (2017), on versioning and concurrency control; the OpenDB documentation on the design database.
