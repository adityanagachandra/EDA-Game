# Debug agents and waveform root cause

When a regression fails, an engineer reproduces it, finds the first wrong value, and traces back to the cause. A debug agent does the same with simulation as its instrument: it reruns, queries waveforms, compares against a reference model, and tests hypotheses by experiment. The simulator makes debugging checkable, since every claim about the root cause can be confirmed by a rerun, and the agent's tool interface decides how efficiently it can ask questions of very large waveform databases.

## Inputs and the first divergence

The agent starts from the failing seed and command line, the scoreboard or assertion message with its time, and the waveform database (VCD, FST or a vendor format). Its first action is to rerun the seed and confirm the failure reproduces deterministically; if it does not, the problem is nondeterminism, a race or an environment issue, and the strategy changes. The scoreboard reports a symptom, often thousands of cycles after the cause: a corrupted pointer at cycle 6,000 may produce wrong data only at cycle 9,000. Root cause is the earliest divergence on the causal path, so the agent compares DUT state, or transactions, against the reference model backwards in time from the symptom.

## Bisection over time and cones

If state matches the reference at cycle 0 and mismatches at cycle $N$, and each rerun from a checkpoint can test whether state still matches at one chosen cycle, bisection finds the first mismatching cycle in $\lceil \log_2 N \rceil$ reruns in the worst case. For $N = 100{,}000$ that is 17 reruns instead of up to 100,000 linear checks. Once the cycle is known, the search moves into space: from the first wrong register, trace its fan-in cone backward to the drivers whose values were wrong at the previous cycle, and repeat until reaching logic whose inputs are right and output is wrong. That logic is the bug's location. Both searches are mechanical, which is what makes them suitable for an agent with good tools.

## Tools, hypotheses and review

Waveform databases are huge, so the agent needs a typed query interface rather than raw text: "values of these signals between these cycles", "first time this signal differs from the reference", "drivers of this net". Exposed as tools, for example over MCP, such queries fit context limits and avoid parsing errors. A proposed root cause is a hypothesis; it is confirmed by experiment: patch or force the suspected logic, rerun the same seed, and check that the mismatch disappears and no new one appears. A failure that vanishes when dumping is turned on points to a race. The output is a report with the reproducing seed, the divergence cycle and location, and the confirming rerun, plus the fix as a diff, because design edits are core changes that a human reviews.
