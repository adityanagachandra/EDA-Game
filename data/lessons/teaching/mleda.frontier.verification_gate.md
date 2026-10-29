# Verification as the gate on learned outputs

An agent-first flow uses learned models everywhere: predictors, placement policies, RTL and script generators. The architectural rule that keeps it trustworthy is that every learned output that is kept passes a gate run by an engine with known guarantees. Designing the gate means choosing the right check for each output, placing it in the harness loop, and knowing what it cannot see.

## Every learned output has a checker

Each kind of learned output has its own gate. A learned macro placement is gated by legalization and a routed evaluation of timing and congestion. An RTL rewrite is gated by sequential equivalence against the original, plus synthesis to confirm the claimed benefit. A predictor's claim that a path is safe is gated by incremental STA on that path. A generated flow script is gated by running it in a sandbox and checking the tool's state and reports. In each case the gate is an engine whose answer does not depend on the model being right.

## Placing the gate in the harness

In the harness loop the gate sits after the proposal and before promotion: golden and differential checks run on every proposal, and failures go back to the model as feedback. Not every output needs the full gate. When a fast mode ranks thousands of candidates, rejected candidates may rest on the fast mode alone, since an error there costs only a missed opportunity; every candidate that is accepted or reported must pass the accurate gate. When the full check is too slow, cheaper incremental or differential checks gate inner-loop steps and full checks gate promotion. Recording each verdict makes decisions auditable and yields grounded training labels.

## What the gate cannot see

A gate is only as strong as its check. A directed testbench passes incorrect designs that agree on the tested inputs. A gate that compares proxy costs checks the same proxy the model optimized, so it cannot catch outputs that exploit the proxy's blind spots. A formal pass on a generated assertion can hide vacuity. A fast timer used as a gate is acceptable for intermediate decisions only if signoff STA gates the final result. The design rule follows: the gate must use a check independent of the model and of the objective the model was trained on, and its weaknesses must be known.

## Going deeper

Kuehlmann and Krohm, "Equivalence Checking Using Cuts and Heaps" (DAC 1997); Amodei et al., "Concrete Problems in AI Safety" (2016), on reward hacking and proxy objectives.
