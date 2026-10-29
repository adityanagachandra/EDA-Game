# Rewrite rule sets as skills

An equality-saturation engine splits naturally into a fixed core and a variable surface. The core, hashconsing, union-find, rebuild and the extraction interface, guarantees that every represented term is equivalent. Rule sets, schedulers and cost models vary per user. Packaging the variable parts as skills with manifests lets an agent pick and extend them safely.

## Core versus change surface

The immutable core of an e-graph engine is the machinery that makes it correct and fast: hashconsing of e-nodes, union-find over e-classes, congruence rebuild, and the extraction algorithm interface. Rule sets (for adders, multipliers, shifts, a customer's MAC unit), rule schedulers, node and time budgets, and cost models are the change surface. A manifest for a datapath rewriting skill lists the rule sets it supports, the cost models it can extract with, its input and output formats (for example word-level RTLIL or another IR) and its tunable budgets. Raising a node budget is Tier 1, adding a rule set is Tier 2, editing union-find is Tier 3.

## Golden tests for rules

A new rule is promoted only with evidence. Soundness first: both sides must be equal at every bit width and signedness the rule matches, which for small widths can be checked exhaustively and in general with an SMT or bit-vector proof. For example, $x \cdot 2^k \to x \ll k$ holds for two's-complement operands of fixed width, but rules involving division or right shifts need sign conditions. Then equivalence: designs extracted with the rule must pass equivalence checking against their originals. Then a differential QoR report over a benchmark suite shows which designs moved and by how much. Versioning rule sets separately from the engine makes each change traceable and reversible.

## Seams and the flywheel

Two seams carry most customer-specific behavior. The cost seam: the extractor asks, for each e-node, what it costs, given operator, widths and context, so a customer's power or area model plugs in without touching the core. The scheduling seam: a policy decides which rules run when, so explosive rules such as associativity do not exhaust the budget before useful rules fire. A rule set that proves valuable for one design moves through the skill flywheel, iterated as a plug-in, validated by golden and differential tests, then promoted into the shared library with guidance on when to select it. The customer's design stays behind the IP boundary.

## Going deeper

Willsey et al., "egg: fast and extensible equality saturation", POPL (2021), on rule scheduling and extraction; the agent-first engine model in this curriculum's agentic domain.
