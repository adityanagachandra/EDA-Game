# A placement rule as a Tier 2 plug-in

Customers often want placement behaviour no general placer offers: keep datapath bit-slices abutted, align regular arrays, cluster cells of one power domain near its switch. In an agent-first placer such a rule is a Tier 2 plug-in at a cost seam, which scores arrangements while the core keeps legality. "Abut regular cells" is the canonical example of this kind of plug-in.

## The seam and its context

A rule plug-in implements a narrow interface: given a batch of cells with their positions, widths, rows and tags, return a cost (and, for global placement, optionally a gradient). For "abut regular cells", the context must include a group tag per cell, so the rule knows which cells form the regular array, and row membership and x order, so it can measure gaps. The plug-in never writes positions; the core's optimiser adds its cost to the wirelength change of each candidate move and the core's legalizer still decides what is legal. Because the rule only expresses a preference, gaps remain legal, it is not a core change, and because it needs new logic, it is not a parameter.

## Computing the rule's cost

For each row, walk the members of group R in x order and add the distance from each member's right edge to the next member's left edge, whatever lies between, then multiply by the weight $\lambda$. With members at x = 0 (width 2), 5 (width 3) and 12 in one row, the gaps are 3 and 4, and with $\lambda = 2$ the cost is 14. The weight converts the rule's site units into the same scale as wirelength and is a Tier 1 parameter the user can tune. The cost should be computed incrementally, touching only the rows a move changes; recomputing every row for every candidate move makes a correct plug-in slow, and the seam should be called per batch, not per cell, so the core's inner loop stays fast.

## Proving the plug-in

A rule plug-in is accepted through golden and differential tests. Golden tests check its cost on hand-built rows with known answers, including edge cases such as one-member groups and groups split across rows. The differential test runs the placer with and without the plug-in on a regression suite: on the target design, group gaps should shrink; on every design, HPWL, overflow and timing must stay within tolerance, and legality must hold, which the core guarantees regardless of the plug-in's output. After promotion, the plug-in ships with a manifest describing the rule, its inputs and its tunable weight, so agents can select it for designs with regular datapaths. The customer's design data stays on-site; only the skill is shared.

## Going deeper

The agent-first EDA notes in this curriculum, sections on seams, extension points and golden testing; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4, for the detailed-placement moves the plug-in scores.
