# Agents that insert and verify test logic

DFT insertion, scan stitching, compression, test points and wrappers, is a natural target for agents: it is configuration-driven, it has clear rules, and it has strong automatic checks. The agent's job is not to edit netlists freely but to propose configurations, run the insertion engine, and prove with equivalence, design-rule and coverage checks that the change did what was intended and nothing else.

## What the agent must prove

Every insertion is checked three ways. **Function**: with scan enable and test modes tied to functional values, the scan muxes and test points are transparent, so combinational equivalence checking (a miter between pre- and post-insertion netlists) must prove the two equivalent. **Structure**: scan design rules (controllable clocks and resets, no loops) must pass, and a **chain integrity** test that shifts a known pattern such as 0011 through every chain must come out intact, because equivalence in functional mode says nothing about the scan path. **Testability**: ATPG coverage and pattern count must meet or beat the baseline.

## The harness loop

The agent works in the edit-build-test loop. It proposes a configuration (chain count, compression ratio, test-point budget, ordering policy), the harness runs insertion, DRC, equivalence and ATPG, and compares coverage and pattern counts with the golden baseline. Failures go back to the agent as feedback; success produces a reviewable diff with the evidence attached for a human. Most requests land at Tier 1 (configuration) or Tier 2 (a plug-in selection rule, such as placing test points only in modules that miss coverage); the netlist database and insertion mechanics are core and are not edited.

## Invariants versus proxies

Proxies invite trade-offs against invariants. An agent rewarded for balanced chains may cross clock domains without lockup latches and create hold failures in shift; one rewarded for fewer test points may raise pattern count. A coverage drop after a passing equivalence check is a regression, typically from new tied or blocked logic, and must be diagnosed before promotion. Configurations must be versioned and replayable on the next netlist. When a recipe works, the skill flywheel promotes the recipe and the knowledge of when to use it into the shared library; the customer's netlist stays on site.

## Going deeper

The agent-first EDA notes in this curriculum (harness loop, golden testing, skill flywheel); Wang, Wu and Wen (2006), chapter 2, for scan design rules.
