# Defects, faults, errors and failures

Manufacturing test exists because fabrication is imperfect: particles, missing vias, shorts and drifted devices make some dies bad. Test engineering separates the physical cause from the logical model a tool can compute with. A tool designer needs this split because every downstream engine, from fault simulation to ATPG to diagnosis, reasons about a finite list of modeled faults, while the quality the customer sees depends on physical defects that the model only approximates.

## The causal chain

A **defect** is a physical imperfection in one die: a particle bridging two metal lines, a void in a via, a gate-oxide pinhole. A **fault** is the abstraction of that defect at a level a tool can reason about, such as "net n is stuck at 0" or "nets x and y are wired-AND". When the circuit is driven so that the faulty line takes a value different from the fault-free one, the fault is *activated* and produces an **error**, a wrong internal value. If the error travels to an observed output, the die shows a **failure**. Each link can break: a defect can produce no fault at the logic level, a fault may never be activated, and an error may be masked before it reaches an output. A test is therefore a set of input vectors chosen to activate each modeled fault and make its error visible.

## Why tools need fault models

Physical defects are unbounded in kind, size and location, so they cannot be enumerated. A fault model replaces them with a finite list that is linear or nearly linear in circuit size. The single stuck-at model gives exactly two faults per line; for a netlist of one million lines that is two million faults, which fault simulation can grade and ATPG can target. The price is that coverage is always measured *against the model*. A test set with 99 percent stuck-at coverage can still miss resistive opens that only add delay, bridges whose behaviour depends on drive strength, and transistor defects inside a cell. That gap is why production flows add transition, path-delay, bridging and cell-aware models. Parametric loss, where the die is functionally correct but too slow or too leaky because of process variation, is not a localized defect at all and is screened by at-speed and parametric tests.

## Yield, quality and reliability

Three numbers describe a test program. **Yield** is the fraction of manufactured dies that are good; test cannot raise it, only measure it. **Quality**, usually quoted as the **defect level** in defective parts per million (DPPM), is the fraction of shipped parts that are bad, and it depends on how many defective dies the test lets escape. **Reliability** concerns parts that are good at time zero but fail early in the field from latent defects such as weak oxide; burn-in and stress tests accelerate those failures before shipment. A systematic failure, where every die fails the same vector the same way, points to a design bug, a test-program error or a systematic process problem rather than to random defects, and it is handled by a different team.

## Going deeper

Bushnell and Agrawal, "Essentials of Electronic Testing for Digital, Memory and Mixed-Signal VLSI Circuits" (2000), chapters 1 to 4; Wang, Wu and Wen, "VLSI Test Principles and Architectures" (2006), chapter 1.
