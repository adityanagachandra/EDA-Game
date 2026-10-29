# Detecting a fault with a test vector

Fault simulation answers one question many times: does this input vector detect this fault? Everything else, coverage, fault dropping, compaction, diagnosis, is built from that predicate. Computing it by hand on small netlists is how a tool designer internalizes the two conditions of detection and why only about half of all faults are even candidates for any given vector.

## Activation and propagation

A vector detects line $x$ stuck-at-$v$ when two things happen. **Activation**: the fault-free circuit drives $x$ to $\overline{v}$, so the faulty circuit differs at $x$. **Propagation**: that difference reaches at least one primary output, which requires every gate along some path from $x$ to have its other inputs at non-controlling values (1 for AND and NAND, 0 for OR and NOR). Procedurally: simulate the fault-free circuit, simulate again with $x$ forced to $v$, and compare the outputs. If any output differs, the vector detects the fault. For example, with $g = $ AND(a, b) feeding $z = $ OR(g, c), the vector a = 1, b = 1, c = 0 detects $g$ stuck-at-0, while c = 1 would block it at the OR.

## Counting detected faults

For one vector, every line has a definite fault-free value, so only the fault stuck at the opposite value can be activated: at most half of the $2L$ faults are candidates. For each candidate, check propagation. A useful shortcut is that the set of lines whose change propagates to an output is the set of **sensitive** lines of the vector; the detected faults are exactly the sensitive lines stuck at the complement of their value. Critical path tracing exploits this to find all detected faults of a vector in one backward pass for fanout-free regions, instead of one simulation per fault.

## Counting detecting vectors

Turning the question around, the number of vectors that detect a fault measures how hard it is. For a circuit with $n$ inputs, exhaustively simulating all $2^n$ vectors gives the count exactly; dividing by $2^n$ gives the fault's random-pattern detection probability. A fault detected by no vector is **undetectable** (or redundant): the logic around it is such that it can never change an output, often because of redundant gates or tied inputs. Such faults are excluded from test coverage. A fault detected by a single vector out of thousands is random-pattern resistant and is the kind of fault ATPG and BIST test points are aimed at.

## Going deeper

Bushnell and Agrawal (2000), chapter 5; Abramovici, Menon and Miller, "Critical path tracing," IEEE Design and Test, 1984.
