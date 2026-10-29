# Golden tests for rule-deck changes

Rule decks are software, and agents will increasingly edit them: adding a rule, tightening a value, speeding up a check. A deck change is safe only if golden tests pin each rule's intended behaviour and a differential run on reference layouts shows that only the intended results moved. This concept turns the golden and differential testing pattern into the concrete test suite a DRC deck needs.

## Golden clips for each rule

A golden test for a rule is a set of small layouts that must pass and others that must fail, each with the expected violation markers. Boundaries matter most: a 28 nm spacing rule needs clips at 27, 28 and 29 nm to catch off-by-one and inclusive-versus-exclusive errors. A conditional rule needs clips that cross each condition threshold, width and run length, separately and together, because implementations disagree at the corners. The suite lives in version control beside the deck, so every deck version has its known-good outputs, and it is kept for every rule even after the deck stabilizes, because refactors and engine upgrades change behaviour silently.

## Differential runs on reference designs

A differential run executes the old and new deck on a fixed set of reference layouts and compares violation reports rule by rule. For an intended change to one end-of-line rule, the proof is that only that rule's counts changed and only where the new condition applies. Movement in unrelated rules signals a side effect, often through a shared derived layer that many rules use, and blocks promotion. The reference set should include several real layouts for production impact and synthetic stress layouts that exercise many rules densely. Running two DRC engines on the same rule and clips is a differential test too.

## What passing does not prove

Tests prove only what they cover. If an agent adds a new condition and all golden tests pass but none probes that condition, nothing about the new behaviour has been shown; the change must add clips that pin it. If the same agent writes the change and the only clips that judge it, a shared misunderstanding passes both, so clips should come from an independent author or from the manual by a separate process. A performance rewrite that leaves reference reports identical can still break boundary cases that production layouts rarely hit, so the boundary clips must run too.

## Going deeper

The KLayout DRC reference documentation, for a rule-deck language that test clips exercise; the agent-first EDA anchor on golden and differential testing.
