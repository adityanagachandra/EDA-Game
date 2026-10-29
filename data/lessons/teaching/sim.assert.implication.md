# Implication and sequence checking

Most concurrent assertions are implications: whenever an antecedent sequence matches, a consequent must follow. Evaluating them on a trace by hand is the best way to understand what a simulator's assertion engine does each cycle, how attempts overlap, and why a property can pass without checking anything. Vacuity is a central concern for both simulation and formal verification, because a vacuous pass looks exactly like a real one in a summary report.

## Overlapping and non-overlapping implication

`a |-> b` is **overlapping** implication: if $a$ holds at a clock tick, $b$ must hold at the same tick. `a |=> b` is **non-overlapping**: $b$ must hold at the next tick, and it means the same as `a |-> ##1 b`. More generally `a |-> ##k b` checks $b$ exactly $k$ ticks after $a$, and `a |-> ##[1:3] b` accepts $b$ at any of the next three ticks. A new attempt of the property starts at every clock tick, so attempts overlap: an antecedent at cycle 3 and another at cycle 4 each create their own obligation. Operand values are the sampled values at each tick, that is, the values just before each clock edge.

## Counting failures on a trace

Evaluate each attempt separately. Take `a |-> ##1 b` over cycles 1 to 8 with $a$ = 1 0 1 1 0 0 1 0 and $b$ = 0 1 0 1 0 1 1 0. Attempts starting in cycles 1 to 7 complete inside the trace. Antecedents hold at cycles 1, 3, 4 and 7. The consequent is checked one cycle later: $b$ at cycle 2 is 1 (pass), at cycle 4 is 1 (pass), at cycle 5 is 0 (fail), at cycle 8 is 0 (fail). So 2 attempts fail. With `a |-> b` instead, the checks are at cycles 1, 3, 4 and 7: $b$ = 0, 0, 1, 1, so the attempts at cycles 1 and 3 fail. A simulator reports each failure at the cycle where the consequent is found false, which may be several cycles after the antecedent.

## Vacuous success

An attempt whose antecedent does not match succeeds **vacuously**: the implication is true because nothing was required. In the trace above, the attempts at cycles 2, 5 and 6 are vacuous, the two at cycles 1 and 3 are real passes and two fail. A property that only ever passes vacuously has checked nothing: `req |-> ##1 ack` on a testbench that never asserts req reports success on every cycle. Simulators can count vacuous and non-vacuous passes separately, and good practice adds `cover property (a)` for each antecedent so coverage proves the check was exercised. In formal verification, vacuity checking asks whether the antecedent can ever occur at all under the constraints, catching over-constrained environments.

## Going deeper

IEEE Std 1800-2017, clause 16.12 (implication); Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015); Beer, Ben-David, Eisner and Rodeh, "Efficient detection of vacuity in temporal model checking", Formal Methods in System Design (2001).
