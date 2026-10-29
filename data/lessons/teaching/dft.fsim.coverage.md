# Fault simulation and coverage

Coverage is the number a test program is judged by: what fraction of the modeled faults the pattern set detects. There are several coverage metrics that differ only in the denominator, and confusing them is a classic source of disagreement between design teams and test engineers. A tool designer must report them precisely, because customers set sign-off targets on them and because their definitions are easy to game.

## Fault classes after ATPG

After ATPG and fault simulation, each fault in the (usually collapsed) list falls in one class. **Detected** faults have a pattern that detects them. **Untestable** faults are proven undetectable: redundant logic, inputs tied to constants, or blocked paths mean no vector can ever detect them. **Aborted** faults hit the ATPG's backtrack or time limit, so the tool neither found a test nor proved untestability. **Not detected** covers the rest, for example faults ATPG never targeted. The proof of untestability matters: an aborted fault might be testable, an untestable one is certainly not.

## Fault coverage and test coverage

**Fault coverage** is detected over all faults, $FC = D/T$. **Test coverage** removes the untestable faults from the denominator, $TC = D/(T - U)$, on the grounds that no test can detect them. Test coverage is therefore higher. With $T = 25{,}959$, $D = 24{,}974$ and $U = 155$, $FC = 96.21$ percent and $TC = 24{,}974/25{,}804 = 96.78$ percent. Because $U$ appears in the denominator, a tool that proves more faults untestable raises TC without detecting anything new, which is legitimate only if the proofs are real; this is exactly the kind of metric an automated optimizer could exploit.

## ATPG efficiency

**ATPG efficiency** measures how well the tool resolved the fault list, $(D + U)/T$: every fault either detected or proven untestable. The gap to 100 percent is the aborted and untargeted faults, which a higher backtrack limit, a SAT-based engine or test points may close. Coverage also depends on the fault list: coverage on a collapsed list differs slightly from coverage on the uncollapsed list because equivalence classes have different sizes, so reports state which list they use. Coverage is always relative to a model; the defect level of shipped parts depends on how well that model represents real defects.

## Going deeper

Bushnell and Agrawal (2000), chapters 5 and 7; Wang, Wu and Wen (2006), chapter 3.
