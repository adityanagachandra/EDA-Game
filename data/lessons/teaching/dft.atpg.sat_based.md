# SAT-based ATPG

SAT-based ATPG turns test generation into Boolean satisfiability: build a formula that is satisfiable exactly when some input vector makes the fault-free and faulty circuits differ at an output, and hand it to a CDCL solver. It reuses the encodings and miters of formal equivalence checking, excels on hard and redundant faults where structural search aborts, and gives a proof when a fault is untestable.

## Encoding the circuit

Each line becomes a Boolean variable and each gate a few clauses by the Tseitin encoding. For $z = $ AND(a, b): $(\overline{z} \vee a)$, $(\overline{z} \vee b)$, $(z \vee \overline{a} \vee \overline{b})$, three clauses; in general a $k$-input AND, NAND, OR or NOR gives $k + 1$ clauses, a NOT or buffer 2, and a 2-input XOR 4. The formula size is linear in the netlist. A circuit with gates AND(2), XOR(2), NOT, NOT, NAND(2) has $3 + 4 + 2 + 2 + 3 = 14$ clauses.

## The ATPG miter

The fault is modeled with a second, faulty copy of the circuit, but only the part the fault can affect: the transitive fan-out of the fault site. Lines outside that cone share variables with the fault-free copy. The faulty copy's site variable is fixed to the stuck value by a unit clause. For each primary output in the cone, an XOR compares the good and faulty versions; an OR of those XORs, asserted true by a final unit clause, requires that at least one output differ. A satisfying assignment restricted to the primary inputs is a test. If the solver proves the formula unsatisfiable, the fault is **redundant**. Activation is implicit: the outputs cannot differ unless the site does.

## Why it works well

Encoding only the fan-out cone keeps formulas small, and adding clauses that say "the site differs" and "a path of differences exists" (the D-chain idea) helps the solver prune. CDCL's clause learning discovers the same kinds of implications FAN and SOCRATES compute, but automatically and incrementally. Incremental SAT reuses the fault-free encoding across many faults, so per-fault cost is small. In production, SAT engines resolve faults a structural engine aborted, prove redundancy, and serve test compaction by adding constraints that a new test also detects other target faults.

## Going deeper

Larrabee, "Test pattern generation using Boolean satisfiability," IEEE Transactions on CAD, 1992; Drechsler, Eggersgluss, Fey and Tille, "Test Pattern Generation using Boolean Proof Engines" (2009).
