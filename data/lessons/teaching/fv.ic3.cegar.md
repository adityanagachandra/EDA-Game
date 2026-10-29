# Abstraction refinement (CEGAR)

Counterexample-guided abstraction refinement (CEGAR) proves properties of a design by model checking a smaller abstract model and refining it only where a spurious counterexample shows it is too coarse. Localization abstraction, which frees latches the property does not seem to need, is its most common form in hardware. The loop is a template that reappears in cut-point equivalence checking, in IC3 variants and in software model checkers.

## The CEGAR loop

CEGAR repeats four steps. Abstract: build a model that over-approximates the design, so that every concrete behavior has an abstract counterpart. Check: model check the property on the abstraction. If it holds, it holds on the design, because no concrete path can reach a bad state when no abstract path can; stop. If it fails, concretize: try to follow the abstract counterexample on the real design, typically with a BMC query constrained to the abstract trace. If the concrete trace exists, it is a real bug; stop. Otherwise the counterexample is spurious: refine the abstraction so that it no longer allows this trace, and repeat. On a finite design the loop terminates, because each refinement adds precision and the worst case is the concrete design itself.

## Existential abstraction

Hiding variables gives an existential abstraction. Each concrete state maps to the abstract state that keeps only the visible variables; an abstract state is initial if some concrete initial state maps to it; and there is an abstract transition from $a$ to $b$ whenever some concrete transition goes from a state mapping to $a$ to one mapping to $b$, whether or not that concrete source is reachable. Example: states $x_2x_1x_0$, hide $x_0$, so abstract states are $x_2x_1$. Concrete transitions 000 → 001, 001 → 010, 010 → 010, 011 → 110 and 110 → 110, initial state 000. Concretely, 000, 001 and 010 are reachable, projecting to 00 and 01. Abstractly, 00 → 00 and 00 → 01, and from 01 the edge 011 → 110 gives 01 → 11, although 011 is unreachable. So 3 abstract states are reachable, one more than the projection.

## Spurious counterexamples and refinement

Extra abstract states and edges are where spurious counterexamples come from. In localization abstraction, a hidden latch becomes a free input that can change arbitrarily each cycle, so an abstract trace may need it to jump from 0 to 1 with no cause. Refinement makes the responsible latches concrete again. Proof-based refinement runs BMC on the concrete design along the counterexample's length; when it is UNSAT, the unsat core names the latches and constraints the infeasibility depends on, and exactly those are added back. Predicate abstraction, common in software model checking, tracks the truth of chosen predicates such as $x < y$ and refines by adding predicates, often obtained by Craig interpolation. Over-approximation is essential: an abstraction that drops behaviors could miss bugs, so it could never prove safety. Abstraction pays off when the property depends on a small part of a large design.

## Going deeper

Clarke, Grumberg, Jha, Lu and Veith, "Counterexample-Guided Abstraction Refinement" (CAV 2000); McMillan and Amla, "Automatic Abstraction without Counterexamples" (TACAS 2003).
