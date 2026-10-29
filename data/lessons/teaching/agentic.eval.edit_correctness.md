# Skill-edit correctness beyond compiling

An edit that compiles may still be wrong. For a skill-editing harness the measure that matters is skill-edit correctness: the fraction of change requests for which the model's edit builds, passes the golden tests, shows a clean differential, stays within its tier and does what was asked. Measuring it as a funnel shows where each model fails and why models with similar compile rates differ sharply.

## The correctness funnel

Write the correctness rate as a product of conditional pass rates:
$$P(\text{correct}) = P(\text{build}) \times P(\text{golden} \mid \text{build}) \times P(\text{differential} \mid \text{golden}).$$
If 90% of edits build, 70% of those pass goldens and 80% of those have a clean differential, the correctness rate is $0.9 \times 0.7 \times 0.8 = 0.504$. Each stage removes a different kind of failure: syntax and type errors, then broken behavior on known cases, then collateral damage on reference designs. A model that fails mostly at the differential stage is making edits that work but touch too much.

## Comparing models on correctness

Model A builds 95% of its edits, but only 60% pass goldens and 75% of those are differential-clean: correctness $0.95 \times 0.6 \times 0.75 = 0.4275$. Model B builds 85%, then 75% and 85%: correctness $0.85 \times 0.75 \times 0.85 = 0.5419$. A wins on compile rate by 10 points and loses on correctness by 11.4 points. This is why the go or no-go number for a harness model is skill-edit correctness on real engines, and why generic code-generation benchmarks, which mostly measure the first stage on fresh code, under-separate models on this task.

## Scope and minimality

A correct edit also stays in scope. A Tier 2 request answered by an edit that also touches core files is not correct even if every test passes, because it changed what the request did not authorize. Benchmarks therefore check the diff: which files and functions changed, whether they lie inside the extension point, and how large the change is. Of 40 edits that pass all tests, if 6 touched files outside their tier, the in-scope correctness is $34/40 = 85\%$ of the passing edits. Minimal diffs are also easier for a human to review, which is the final gate.

## Going deeper

The agent-first EDA notes in this curriculum, sections on the harness loop and on evaluation.
