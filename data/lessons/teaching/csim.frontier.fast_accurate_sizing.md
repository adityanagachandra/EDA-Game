# Fast and accurate simulation in a sizing loop

Analog sizing loops evaluate thousands of candidate designs. An agent cannot afford accurate SPICE at every corner for every candidate, so it runs a fast simulator, a surrogate or a fast SPICE mode, to screen, and recurses into accurate SPICE only where the decision depends on it. Budgeting that loop, and knowing exactly which candidates the fast engine's error forces back to SPICE, turns the fast-and-accurate idea into numbers.

## Budgeting a two-level loop

A two-level loop screens $N$ candidates with a fast engine at $t_f$ per run and verifies the best $k$ in accurate SPICE at $C$ corners, $t_a$ per run. The total is $N t_f + k C t_a$. With 1000 candidates at 2 s, then the best 10 at 15 corners at 3 min, screening takes 33.3 min and verification 450 min, 483 min in all. Verification dominates, which is typical: corners multiply accurate cost. The budget tells an agent where to spend effort: a faster screen saves little here, while pruning corners per metric or verifying fewer candidates saves a lot. All runs are independent, so the wall time also divides by the available parallel licenses.

## Speedup over all-accurate

Running every candidate at every corner in accurate SPICE costs $N C t_a$. In the example, $1000 \times 15 \times 3 = 45000$ min against 483 min, a speedup of 93. The speedup is real only if the fast screen ranks candidates well enough that the true best design is among the $k$ verified. That requires the fast engine to be validated against SPICE golden results over the region the search explores. A surrogate trained on a different region can rank confidently and wrongly, and a reward computed from the fast engine alone invites an optimizer to exploit its errors, which is reward hacking against a proxy.

## Which candidates must be re-verified

A validated fast engine comes with an error bound $\pm e$. With a specification of gain at least $S$, a fast value above $S + e$ passes regardless of the error, a value below $S - e$ fails regardless, and a value within $[S - e, S + e]$ could go either way and must be re-simulated accurately. With $S = 70$ dB, $e = 1.5$ dB and fast gains 70.3, 67.5, 67.9, 66.9, 71.4, 68 and 73.9 dB, the band is 68.5 to 71.5 dB and two candidates fall in it. Tightening the fast engine's error bound shrinks the band and the accurate workload, which is the measurable payoff of better surrogates.

## Going deeper

Gielen and Rutenbar, "Computer-aided design of analog and mixed-signal integrated circuits" (Proceedings of the IEEE, 2000); Rutenbar, Gielen and Roychowdhury, "Hierarchical modeling, optimization, and synthesis for system-level analog and RF designs" (Proceedings of the IEEE, 2007).
