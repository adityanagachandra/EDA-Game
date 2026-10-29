# Design-rule manuals and rule decks

The design-rule manual (DRM) is the foundry's human-readable definition of every rule, with figures and conditions. The rule deck, or runset, is the executable form that a DRC engine runs. Modern nodes have thousands of rules, most of them conditional, and the gap between the manual's words and the deck's code is where many tool bugs and design surprises live.

## From simple numbers to conditional rules

Early processes needed a few dozen rules of the form "width at least $w$, space at least $s$". Advanced nodes add rules whose values depend on context: spacing that grows with parallel run length or wire width, end-of-line spacing for line tips facing other shapes, enclosure that depends on the direction of the line, minimum area and minimum enclosed area, colour-dependent spacing for multiple-patterned layers, antenna ratios limiting metal area connected to a gate during plasma processing, and density windows for CMP. Each condition encodes a process mechanism such as line-end pull-back, charging damage or polish dishing, so reading the mechanism is the fastest way to understand the rule.

## How a deck computes a check

A rule deck is a program in a geometric language. It derives layers by boolean and sizing operations (gate equals poly AND active), selects edges by orientation and context, measures distances between selected edges, and emits violation markers. An extension check for gate poly past active, for example, intersects the layers, takes the poly edges outside active, and measures them. Engines run these operations on polygon databases with spatial indices and sweeps, often hierarchically to reuse results for repeated cells. When the manual and the deck disagree, signoff follows the deck, and the discrepancy is reported to the foundry rather than edited locally.

## Recommended rules and restricted layouts

Beyond mandatory rules, many manuals list recommended rules: wider spaces, double vias, longer extensions. They are not required for signoff but raise yield, so DFM flows score and apply them where area and timing allow. Some hotspots depend on two-dimensional configurations that no width or space rule captures, so decks add pattern-matching rules from libraries of known-bad shapes. Gridded, unidirectional layers replace many of these rules with a few simple ones, at the cost of routing flexibility. For routers the lesson is direct: any rule the router does not model becomes a violation that a later repair loop must fix.

## Going deeper

The KLayout DRC scripting documentation and the Magic VLSI technology-file documentation, as open examples of rule-deck languages; Lienig and Scheible, "Fundamentals of Layout Design for Electronic Circuits" (2020).
