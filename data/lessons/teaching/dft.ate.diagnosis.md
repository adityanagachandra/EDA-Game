# Fault diagnosis and yield learning

When a die fails, the tester log says which patterns failed at which outputs. Diagnosis turns that log into candidate defect locations, and volume diagnosis across thousands of failing dies turns those candidates into the yield-learning signal that tells the fab which layer and mechanism to fix. It is the test flow's feedback loop to manufacturing, and an inference problem a tool designer can attack with dictionaries, cone analysis and statistics.

## Dictionary matching

**Cause-effect** diagnosis simulates candidate faults and compares their predicted behaviour with the observed one. A **pass/fail dictionary** stores, for each fault, the set of patterns it fails. Score each candidate by the mismatch between prediction and observation, the size of the symmetric difference of the two pattern sets. If a die fails patterns {P3, P4, P5, P6, P9} and candidates predict {P1, P2, P3, P4, P5, P7, P9}, {P3, P4, P5, P6, P9} and {P3, P5, P6, P7, P9}, their scores are 4, 0 and 2, and the second matches perfectly. A full-response dictionary also records which outputs fail, which resolves more candidates at a cost of $F \times P \times O$ bits instead of $F \times P$; production tools simulate candidates on demand instead of storing dictionaries.

## Effect-cause analysis

**Effect-cause** diagnosis reasons backward from the failures. Under a single-fault assumption, the defect must lie in the fan-in cone of every failing output on a failing pattern, so intersecting those cones bounds the candidates. Outputs that passed cannot simply be subtracted, because a fault in their cone need not reach them on that pattern. Intersections over many patterns, followed by simulation of the survivors, typically localize a defect to a few nets, which layout then maps to physical locations for failure analysis. Multiple defects, bridges and cell-internal faults require scoring methods that tolerate partial matches.

## Yield learning

**Volume diagnosis** aggregates candidates over many failing dies, attributes them to layers and mechanisms using layout (a candidate net that runs mostly on metal 2 next to a neighbour suggests a metal-2 short), and builds a Pareto of yield loss. With Poisson yield and mechanism contributions $\lambda_i$ to the expected killer defects per die, removing mechanism $j$ raises yield from $e^{-\sum \lambda_i}$ to $e^{-(\sum \lambda_i - \lambda_j)}$. If poly shorts contribute 0.26 of a total 0.53, eliminating them raises yield from 58.9 to 76.3 percent, which tells the fab where its effort pays most.

## Going deeper

Abramovici, Breuer and Friedman (1990), chapter 12; Wang, Wu and Wen (2006), chapter 7; Tang et al. and other volume-diagnosis papers at the International Test Conference.
