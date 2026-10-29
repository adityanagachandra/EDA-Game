# Agentic signoff and ECO loops

Late in a design, closing the last timing, noise, IR and DRC violations is a loop: read reports, propose small changes, apply them incrementally, re-analyze, repeat. That loop is a natural job for an agent, because it is repetitive, report-driven and bounded by engines that check every step. It also shows clearly where a human must stay in the loop: approving the final engineering change order.

## What the loop reads

The agent reads signoff reports as structured data: setup and hold violations by endpoint and corner from the timer, crosstalk delta delays and glitches from signal-integrity analysis, IR drop maps and the cells above a drop threshold from power analysis, and DRC and antenna violations from physical verification. It groups violations by cause, for example a cluster of setup violations through one congested region, or hold violations on paths from a single clock branch, because fixes target causes. It also reads what constrains fixes: available spare cells, free placement sites, routing resources and the area and power budget.

## Proposing and verifying fixes

Fixes are the standard ECO moves: upsize or downsize a cell, swap Vt, insert or remove a buffer, insert hold delay cells, reroute a net with wider spacing or on another layer, add a shielding wire. The agent proposes a batch, applies it in the implementation engines with legal placement and routing, and reruns incremental analysis. Each fix must be checked across all corners and modes and across domains: upsizing a driver to fix setup can create a hold violation elsewhere, raise IR drop in its region or worsen crosstalk on a neighbor. The loop accepts a batch only if the total violation count and severity fall and no new class of violation appears.

## The human-approved ECO

The final change list goes to a human with its evidence: before and after reports for every corner, the list of changed instances and nets, and the checks that ran. A human approves it because an ECO at this stage can affect tape-out and because some decisions, such as whether to waive a marginal violation or accept a small area increase, are engineering judgments with consequences outside the reports. The agent's role is to make that review fast: few, explained changes with complete evidence, rather than a large unexplained diff.

## Going deeper

The OpenROAD documentation on its repair_timing and repair_design commands; Kahng et al., "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure.
