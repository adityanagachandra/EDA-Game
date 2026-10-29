# Proving an agent's RTL refactor

When an agent refactors RTL, passing the existing tests is not proof that behavior is unchanged. What must be proven depends on the kind of change: renaming needs only combinational equivalence, re-encoding needs sequential equivalence, pipelining needs equivalence up to latency, and an intended bug fix needs proof that behavior moved only where intended. This concept turns golden and differential testing into concrete RTL obligations.

## Golden and differential tests

A golden test is a fixed stimulus with recorded known-good outputs, compared after every change: any behavioral change, intended or not, becomes visible. A differential simulation drives the old and new versions with the same stimulus and compares outputs cycle by cycle; it needs no expected values, because the old design is the reference.

Both only cover the stimulus applied. A difference on an input sequence never simulated survives. That is why the reference must be the last approved version, not the agent's previous attempt, so errors cannot accumulate across iterations, and why simulation is paired with formal methods where the change allows it. For a pure rename or reformat, where registers correspond one to one, a combinational equivalence check with registers matched by name or position is cheap and complete.

## Choosing the equivalence method

The method follows from what the refactor did to state. If registers still correspond, compare the combinational logic between them; open-source flows do this with Yosys, reading both designs and proving equivalence with equiv_make, equiv_simple and equiv_induct, or with a SAT-based miter. If the agent re-encoded a state machine from binary to one-hot, the registers no longer correspond, and only sequential equivalence, which reasons over input sequences from reset, applies.

If the agent pipelined a combinational multiplier into two stages, outputs now appear one cycle later, so equivalence holds only against the original delayed by one cycle, and the check must say so. Counterexamples that appear only when an input is x usually mean the two versions treat don't-cares differently; constraints on the inputs decide whether that matters.

## Intended changes and quality

Not every change should preserve behavior. For an intended fix, such as an off-by-one in a counter's wrap value, the differential proof must show that outputs differ only in the cycles where the counter wraps and match everywhere else: the analogue of a differential timing report showing that only the intended paths moved. A refactor of a module instantiated in 40 places with different parameters must be checked for each distinct parameterization actually used, since each is separate hardware.

Correct is not the same as acceptable. A refactor that passes equivalence but grows synthesized area by 12 percent is a quality regression, and the harness should report it beside the equivalence pass. Every engine change in an agent-first flow ships with both kinds of check: function against a reference and quality against a baseline.

## Going deeper

The Yosys manual, equiv_make, equiv_simple, equiv_induct and miter; Kuehlmann and van Eijk, "Combinational and Sequential Equivalence Checking" in "Logic Synthesis and Verification" (2002).
