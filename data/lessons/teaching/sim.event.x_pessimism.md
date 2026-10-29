# X-pessimism and X-optimism

An X in simulation should mean "this could be 0 or 1". Simulators do not track that exactly: they apply local truth tables, which lose correlations between signals. The result can be **X-pessimism**, an X where real hardware would produce a known value, or **X-optimism**, a known value where hardware could produce either. Both mislead verification: pessimism produces false failures and wasted debug, optimism hides real bugs. Designers of simulators and X-propagation modes work on exactly this trade-off.

## X-pessimism at gate level

Build a 2:1 mux from gates: $y = (s \,\&\, a) \mid (\lnot s \,\&\, b)$. Let $a = b = 1$ and $s = X$. Hardware outputs 1 whatever $s$ really is. The simulator computes $s \,\&\, 1 = X$, $\lnot X = X$, $X \,\&\, 1 = X$, and $X \mid X = X$: it reports X. The truth tables do not know that the two X terms are complements of one another, so the correlation is lost. This is X-pessimism. It is common in gate-level simulation after synthesis, especially around reset and clock gating, and it makes engineers chase phantom X values. The same structure with $a = b = 0$ gives $0 \mid 0 = 0$, which is accurate.

## X-optimism in RTL control flow

RTL is often optimistic instead. In `if (s) y = a; else y = b;` an X condition is treated as false, so the else branch runs and $y = b$, a known value, even when $a \ne b$ and the hardware could output either. A `case (s)` compares with `===`, so X matches only an explicit X item; with a default the default branch runs, and without one $y$ keeps its previous value. The conditional operator is more careful: `s ? a : b` with $s = X$ merges the two data values bit by bit, passing bits where $a$ and $b$ agree and producing X where they differ. So for $a = 1$, $b = 0$, `s ? a : b` gives X (accurate), while the if statement gives 0 (optimistic).

## Consequences and mitigations

Optimism is the more dangerous failure: an uninitialized flop feeding an if condition can make RTL simulation pass while the gate-level netlist or the silicon misbehaves, a classic simulation-synthesis mismatch. Pessimism wastes time and can push teams to add resets everywhere. Commercial simulators offer X-propagation modes that make RTL constructs more pessimistic (for example merging both branches of an if with an X condition), to expose optimism bugs at RTL. Two-state simulators such as Verilator have no X at all, so they cannot show either effect and instead randomize initial values to shake out dependence on uninitialized state. Formal tools avoid the problem by reasoning about both values of an unknown symbolically.

## Going deeper

Turpin, "The Dangers of Living with an X (bugs hidden in your Verilog)" (SNUG 2003); IEEE Std 1800-2017, clause 12 (procedural if and case semantics).
