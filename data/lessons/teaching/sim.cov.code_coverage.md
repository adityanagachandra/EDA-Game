# Code coverage

Code coverage measures which parts of the RTL source a regression exercised. The simulator instruments statements, branches, toggles and state machines and counts hits. It is automatic, needs no model of intent, and is the first metric any verification plan tracks. A tool designer cares about the precise definitions, because "branch" and "toggle" can be counted several ways, and about the metric's blind spot: it shows what never ran, not whether what ran was checked.

## Statement coverage

Statement (or line) coverage is the fraction of executable statements that ran at least once. The simulator inserts a counter per statement during compilation and reports executed over total. When combining several modules, pool the counts rather than averaging percentages: if module A has 200 statements with 180 executed (90%) and module B has 50 with 25 executed (50%), the design's statement coverage is $205/250 = 82.0\%$, not the 70% you would get by averaging the two percentages. Unexecuted statements are the cheapest bug-finding lead in verification: code that never ran has never been tested at all, and it often marks a whole feature that the stimulus cannot reach.

## Branch coverage

Branch coverage counts the outcomes of each decision. An if statement with an else has two outcomes. An if / else-if chain with $m$ conditions and no final else has $m + 1$ outcomes, because falling through all conditions is a branch even though no code is written for it. A case statement has one outcome per item, plus the default if it has one. Example: a chain with 2 conditions (3 outcomes) and a case with 4 items plus a default (5 outcomes) has 8 branch outcomes; if the tests reached 2 and 4 of them, branch coverage is $6/8 = 75.0\%$. Branch coverage is stricter than statement coverage: the implicit else of a chain can be missed while every written statement executed.

## Toggle and FSM coverage, and the limits

Toggle coverage records, for each bit of each signal, whether it made a 0 to 1 transition and a 1 to 0 transition. Counting the two directions as separate items, a 4-bit signal has 8 items. If it takes the values 0000, 0110, 0100 on successive cycles, bits 1 and 2 rose, bit 1 fell, and bits 0 and 3 never moved: 3 of 8 items, 37.5%. Toggle coverage is cheap and catches stuck control bits and unused bus lines. FSM coverage reports which states and which transitions of an extracted state machine occurred. All code coverage shares one limit: 100% means every structure was exercised, not that the outputs were checked or that the specified scenarios happened. A missing feature has no code to cover.

## Going deeper

Piziali, "Functional Verification Coverage Measurement and Analysis" (2004), chapters on code coverage; IEEE Std 1800-2017, clause 40 (code coverage control and API).
