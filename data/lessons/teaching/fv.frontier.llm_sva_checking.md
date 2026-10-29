# Checking machine-written SVA with the prover

Language models can turn specifications into SystemVerilog assertions quickly, but a generated assertion can be malformed, vacuous, trivially true or simply wrong. The formal engine is the gate that decides which generated properties are kept. A good harness runs a staged series of deterministic checks, from parsing to mutation testing, and feeds every failure back to the model, so that the human only reviews properties that have already earned their place.

## Syntax and elaboration gates

The cheapest checks come first. Parse each generated property, elaborate it against the RTL, and reject anything with syntax errors, unknown signals, width mismatches or constructs the formal tool ignores or rewrites; warnings count as failures, because a construct the tool silently changes no longer means what the model intended. These checks are fast, deterministic and precise, so their messages make excellent feedback for the model's next attempt, in the edit-check loop of an agent harness. Generated properties live in a separate bind file or checker module, not in the RTL, so the design under test is untouched and the change sits at the smallest change tier, easy to review and revert. Only properties that elaborate cleanly move on to the prover, which keeps expensive solver time for candidates that are at least well formed.

## Vacuity and covers

A proved assertion is not yet a useful one. An implication whose antecedent never fires passes vacuously, so the harness covers the antecedent of every generated implication and discards or flags properties whose cover is unreachable. Generated assumptions need more scrutiny than generated assertions, because assumptions are trusted: a wrong one can make every assertion pass by removing behaviors. So the harness keeps a fixed, human-reviewed set of covers for normal operation, such as a completed transfer or a full FIFO, and rejects any generated assumption under which one of them becomes unreachable. Generated assumptions about a block's inputs should also be proved as assertions on the block that drives them. Properties that pass, are non-vacuous and leave the reference covers reachable go on to the last filter, which asks whether they constrain anything at all.

## Does the assertion have teeth

A property can be non-vacuous and still useless, for example a register compared with itself one cycle later when it never changes. Mutation testing measures strength: inject small faults into the RTL (flip a comparison, drop an enable, swap two mux inputs), re-run the prover, and count how many mutants each property kills, that is, makes fail. A property set that passes on the real design and kills most mutants constrains real behavior; surviving mutants point to behavior no property covers, unless the mutation was functionally equivalent. When an expert reference property exists, the harness can check equivalence formally: build both properties as monitors on the same design and prove their verdicts always agree, a miter of two monitors. Throughout, the model proposes and the prover disposes: verdicts come from sound, checkable engines, never from the model's own confidence.

## Going deeper

Claessen, "A Coverage Analysis for Safety Property Lists" (FMCAD 2007); Seligman, Schubert and Kumar, "Formal Verification: An Essential Toolkit for Modern VLSI Design" (2nd ed., 2023), on vacuity and coverage.
