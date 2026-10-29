# Legality rules: sites, rails and regions

At advanced nodes, legal means much more than "in a row and not overlapping". Power rail alignment, multi-row cells, fence and guide regions, blockages and neighbour-dependent spacing all constrain where cells may go. These rules live in the placer's core: every downstream engine assumes their output holds, so a legalizer that misses one corrupts the whole flow.

## Rows, sites and rail alignment

A cell must start on a site boundary of a row and use that row's orientation. Rows alternate N and FS so that adjacent rows share rails, so a cell moved between rows of different parity must be flipped. Multi-row-height cells, such as double-height flip-flops, span two rows and can only start in rows where their bottom rail matches: a double-height cell expecting VSS at the bottom cannot start in a row whose bottom rail is VDD. Some libraries add well or implant rules that force spacing between cells of different types, and tap cells must be present at regular intervals. All of these are discrete constraints on the slot assignment, checked exactly.

## Regions and blockages

A fence region is exclusive: its member cells must be placed inside it and no other cells may enter. A guide region only attracts its members, without excluding others. Hard placement blockages forbid all cells, for example under a sensitive analog block; partial blockages cap the allowed density, leaving room for routing. FIXED instances in DEF, such as macros and pre-placed cells, occupy their sites permanently. A legalizer represents all of these by splitting rows into usable segments per region, so a row becomes a set of intervals, each with an owner. Global placement must also see them, otherwise it spreads fenced cells outside their fence and legalization moves them far.

## Neighbour rules and why legality is core

Some rules depend on the neighbour: edge-type spacing requires a gap between particular cell edges, and pin-access rules can make two specific cells unroutable when abutted. These make legality a property of pairs of adjacent cells, not just of each cell's slot, so the legalizer must check the left and right neighbours whenever it places a cell. Because rules like these grow with each technology, a placer is designed with legality as a small, heavily tested core: a slot model and a rule checker that every move must pass. Heuristics that choose positions can change freely, but the guarantee cannot. A plug-in that prefers some arrangement may score it, but only the core decides what is legal.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4; the ISPD 2015 and ICCAD 2017 placement contest descriptions, on fence regions and advanced-node rules.
