# Clean results as golden and differential tests

For physical verification, the golden test is simple to state: the reference designs must stay DRC-clean and LVS-clean, and the regression designs must produce exactly the expected markers. Differential runs compare marker sets before and after a change. Counting what was introduced, what was fixed and what merely moved is the core of judging an agent's fix or an engine change.

## Introduced markers

A marker is a rule name plus a location. With exact matching, a new-run marker is introduced if no baseline marker has the same rule and position. If the baseline has five markers and the new run has six, of which three match baseline markers exactly, then three were introduced. An accepted fix must introduce none, whatever it did to the old ones: a fix that clears its target marker but adds one elsewhere is a regression. Engines key markers by rule and cell coordinates so that hierarchical results compare across runs.

## Fixed markers

A baseline marker is fixed if it no longer appears in the new run. A marker of the same rule at a different location is a new violation, not the old one moved, so it counts on both sides: one fixed, one introduced. The pair (fixed, introduced) summarizes a change better than the total count, which can stay flat while the layout gets worse. For engine changes rather than layout changes, the expected answer on a regression suite is zero fixed and zero introduced, unless the change was meant to alter results.

## Matching with tolerance

Reruns can shift a marker slightly, through different tiling or sliver handling, without any real change. Differential tools therefore match markers of the same rule within a tolerance, for example a Manhattan distance of at most 3 nm. A marker moved by 2 nm still matches; one moved by 5 nm is unmatched and counts as introduced. Too small a tolerance turns jitter into false regressions, too large one hides real moves. The tolerance is a Tier 1 parameter of the differential check, chosen from the engine's known numerical behaviour.
