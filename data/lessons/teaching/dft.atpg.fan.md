# The FAN algorithm

FAN, fanout-oriented test generation, kept PODEM's idea of deciding only where decisions are safe, but noticed that most decisions are wasted inside fanout-free regions where justification can never fail. By stopping backtrace at head lines, implying mandatory values early and backtracing many objectives at once, FAN cut backtracks substantially. Its techniques, especially aggressive implication and learning, foreshadow how SAT solvers are used in modern ATPG.

## Head lines and bound lines

A line is **bound** if it is reachable from some fanout stem, and **free** otherwise. A **head line** is a free line that feeds a bound line. The cone behind a head line is fanout-free, so any value on it can later be justified without conflict: a tree of gates can always produce any value its function allows. FAN therefore stops backtrace at head lines rather than primary inputs, makes decisions there, and justifies the head lines last, in a final conflict-free pass. This removes all branching inside fanout-free regions.

## Unique sensitization and implication

FAN implies as much as possible before deciding. **Immediate implication** propagates mandatory values forward and backward after every assignment, finding conflicts early. **Unique sensitization** looks at the D-frontier: if every path from it to the outputs passes through one gate, that gate's side inputs must be non-controlling, so they are set immediately without a decision. These steps turn would-be decisions into forced values.

## Multiple backtrace

PODEM backtraces one objective at a time, which can pick a stem value that a later objective contradicts. FAN's **multiple backtrace** traces all current objectives together, counting at each fanout stem how many objectives request 0 and how many request 1. The stem is assigned the majority value, which satisfies more objectives and backtracks less; a stem with three requests for 1 and one for 0 gets 1. Production ATPG tools combine such structural engines, which are fast on the many easy faults, with SAT-based engines for the hard and redundant ones; static and dynamic implication learning, as in SOCRATES, bridges the two.

## Going deeper

Fujiwara and Shimono, "On the acceleration of test generation algorithms," IEEE Transactions on Computers, 1983; Schulz, Trischler and Sarfert, "SOCRATES," IEEE Transactions on CAD, 1988.
