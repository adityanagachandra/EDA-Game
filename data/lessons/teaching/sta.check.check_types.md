# Choosing the right timing check

A timer applies different checks at different pins: setup and hold at flop data pins, recovery and removal at asynchronous controls, clock-gating checks at gates that combine a clock with an enable, minimum pulse width at clock pins, and data-to-data checks where two signals must stay aligned. A tool designer must know which check each structure needs, because a missing check is a silent hole in signoff.

## Data and asynchronous checks

Setup guards against data that arrives too late for the next capturing edge; hold guards against new data that arrives too early and corrupts the current capture. Both are checked at the data pins of flip-flops and latches against the clock pin, using Liberty check arcs. For asynchronous set and reset, the deasserting edge is checked against the clock with recovery, the setup-like check, and removal, the hold-like check. Assertion is not checked because it acts without the clock. For a latch, setup is checked against the closing edge with borrowing allowed after the opening edge. Each check type has its own path group in reports, so a design can pass setup while failing recovery.

## Clock-pin and gating checks

When a gate such as an AND combines a clock with an enable, the enable must change only while the clock is in its controlling state, low for an AND, or the gated clock glitches. Timers infer clock-gating setup and hold checks on the enable pin, or apply them from integrated clock-gating cells in the library. Clock pins also need a minimum pulse width: unequal rise and fall delays along a long clock tree, or a divider, can shrink the high or low phase below what the flop needs to latch reliably. Liberty gives the limit, and the timer checks the propagated pulse. A data-to-data check, set with set_data_check, times one data pin against another when two signals must arrive within a window and neither is a clock.

## Matching checks to failures

The right fix depends on which check failed. A hold failure needs delay added on the short path, sized at the fast corner and limited by setup slack at the slow corner; changing the period does nothing. A path between opposite clock edges gets a partial cycle for setup and an easy hold check. An input from an unrelated clock domain has no defined launch-capture relationship at all: no setup or hold check is meaningful, so the crossing is declared asynchronous and handled by a synchronizer. When reading a report, first identify the check type and its capture edge; many apparent violations come from checks applied with the wrong edge relationship.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8; the Liberty reference manual on timing_type values.
