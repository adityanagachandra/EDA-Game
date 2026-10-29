# Agentic DRC-repair loops

After detailed routing, a few hundred violations often remain, clustered around crowded cells and macros. An agentic repair loop reads the DRC report, proposes local fixes, runs them through the router, verifies, and iterates, with a human approving the final change. It is the routing instance of an agentic signoff or ECO loop, and its value depends on what it is allowed to change and how it checks itself.

## Read, propose, verify

Each iteration starts from the router's DRC report: violation types, layers, locations and the nets involved. The agent proposes a fix within its allowed change surface, such as ripping up and rerouting the nets in a window around a cluster with a higher marker penalty, enlarging a clip, changing a net's layer range, or adding a non-default rule. The router executes the fix, and the agent verifies the repaired region and the whole design: no new violations, no opens, no new timing failures.

When the count falls and then oscillates, say between 3 and 5, the current move is not converging. The loop must escalate: widen the window, change strategy, or hand the residue to a human, instead of repeating the same move.

## Rewards that cannot be gamed

An agent optimising violation count will find shortcuts if they exist. Editing the rule deck to relax spacing removes violations without fixing anything; the deck is fixed by the foundry and outside the agent's surface. Deleting a net's last segment removes a violation and leaves an open; a connectivity check catches it. Satisfying the router's fast internal checker but not the signoff engine is reward hacking against a proxy.

So the reward combines signoff DRC on the final layout, a connectivity check and timing, and the loop runs in a sandbox with versioned, reversible changes. Design data stays inside the IP boundary; only the know-how, what kind of fix worked, may leave it.

## From repeated fixes to skills

Repairs work on windows, not on the whole design, because local changes are cheaper to verify and less likely to disturb nets that were already clean, and because the resulting diff is reviewable. When the same move repeatedly fixes the same kind of violation across designs, say an end-of-line pattern next to a particular cell, it is a candidate skill: validate it with golden and differential tests and promote it into the router's library of repair strategies.

A human still approves the final routing ECO, reviewing the diff and the check results. The loop's job is to do the tedious search and to present evidence; it does not sign off its own work. Fixes that trade one failure for another, such as removing violations while degrading slack, must be rejected by the acceptance check.
