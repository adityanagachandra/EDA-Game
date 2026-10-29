# LLMs for assertion generation

SystemVerilog assertions turn design intent into properties that simulation and formal tools check. Writing them is skilled work, and LLMs can draft them from specifications. Assertions are short and checkable, which suits generate-and-verify loops, but a generated assertion can be wrong, vacuous or a restatement of the code, and each failure needs its own check.

## From intent to SVA

The sentence "every request is granted within 3 cycles" becomes `req |-> ##[1:3] gnt`: an implication triggered by each request, with the grant required 1 to 3 cycles later. A safety property such as "never push when full" becomes an assertion of `!(push && full)` at each clock, with reset disabling it. The model should be given the specification and the interface, not only the RTL: properties derived from the implementation tend to restate it, bugs included. Assertions, assumptions and covers have different roles, checking the design, constraining its environment and showing reachability, and a generator must use the right one.

## Checking a generated assertion

Three checks apply to every generated property. Run it against the design with a formal tool or simulation: a failure is either a design bug or a misstated property, and the specification decides which. Check non-vacuity: an implication whose antecedent never fires is trivially true, so a cover on the antecedent must be reachable. Measure detection power: inject mutations into the RTL and count how many mutants some assertion catches. Generated assumptions need extra care, because an over-constrained assumption can exclude exactly the inputs where the design fails and make every proof meaningless.

## Failure modes

Off-by-one timing is the classic error: `a |=> b` checks b in the cycle after a, while `a |-> b` checks the same cycle, and both compile. Assertions that mirror an RTL expression pass on every implementation that shares the bug. Metrics that count compiled or passing assertions reward these failures; metrics that count injected bugs caught, non-vacuous passes and specification clauses covered do not. In a verification flow the generator proposes, the formal or simulation engine checks, and failing or vacuous properties are triaged before they join the regression. A model must never edit the design until its own properties pass.

## Going deeper

Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015); IEEE Std 1800, the SystemVerilog standard, chapter 16.
