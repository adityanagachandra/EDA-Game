# PODEM

PODEM, path-oriented decision making, reformulated test generation so that decisions are made only on primary inputs. Every internal value is then implied by simulation, so assignments can never be inconsistent, and the search tree is bounded by the $2^n$ input space. That one change made ATPG dramatically faster on circuits with reconvergent logic, and the objective-backtrace-imply loop it introduced is still the skeleton of structural ATPG engines.

## Objectives

PODEM repeatedly asks what it needs next, an **objective**: a line and a desired value. If the fault is not yet activated, the objective is the fault site at the opposite of the stuck value. If it is activated, PODEM picks a gate on the D-frontier and sets one of its X inputs to the gate's non-controlling value, so the error can pass. For example, with the fault on b stuck-at-1 and all inputs X, the objective is b = 0. Once b carries $\overline{D}$, if $g_3 = $ NAND(b, $g_1$) is the chosen frontier gate and $g_1$ is X, the objective becomes $g_1 = 1$. An X-path check confirms some path of X values still links the frontier to an output; if none does, PODEM backtracks.

## Backtrace

An objective on an internal line is converted to a primary input assignment by **backtrace**. At each gate, compute the value needed on an input, inverting through NAND, NOR and NOT. If that input value is the gate's controlling value, one input suffices, so follow the **easiest** input (lowest SCOAP controllability for that value). If it is the non-controlling value, every input will be needed, so follow the **hardest** input first, to fail early. Repeat until a primary input is reached; that input and value are the decision. Backtrace makes no assignment itself; it only chooses where to decide.

## Implication and backtracking

After each input decision PODEM simulates the circuit in five-valued logic, which implies every internal value. If an output carries D or $\overline{D}$, the inputs are a test. If the fault site now holds the stuck value, or the D-frontier is empty, or no X-path remains, the current assignment cannot lead to a test: PODEM flips the most recent untried decision, and pops it when both values have failed. The decision tree has depth at most $n$, so an exhausted tree proves the fault redundant, while hitting a backtrack limit only aborts it.

## Going deeper

Goel, "An implicit enumeration algorithm to generate tests for combinational logic circuits," IEEE Transactions on Computers, 1981; Bushnell and Agrawal (2000), chapter 7.
