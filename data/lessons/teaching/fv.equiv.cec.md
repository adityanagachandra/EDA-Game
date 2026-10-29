# Combinational equivalence checking

Combinational equivalence checking (CEC) proves that two netlists with matching state elements compute the same next-state and output functions. It runs after every synthesis, optimization and ECO step in a production flow, so it must be fast, fully automatic and trustworthy. Its engine is a pipeline: match the designs, build miters, simplify aggressively, and leave as little as possible for the SAT solver.

## The CEC flow

CEC needs a correspondence between the two designs: primary inputs and outputs by name, and registers matched one to one, usually by name and then by structural or functional matching for renamed ones. With registers matched, each register's output becomes a shared free input (a pseudo-input) and its next-state function a pseudo-output, so the sequential designs reduce to a set of combinational comparisons, one per output or register input. Each comparison is a miter over that output's cone of logic. A satisfiable miter gives a distinguishing vector over inputs and register values; an unsatisfiable one proves that cone equivalent. If synthesis retimed or re-encoded the registers so that no one-to-one correspondence exists, the problem is no longer combinational and needs sequential equivalence checking.

## Simplify before you solve

Modern CEC engines are organized so that the SAT solver sees as little as possible. First, both designs are merged into one AIG with structural hashing, which proves structurally identical logic equivalent for free; when synthesis made only local changes, many outputs already meet at the same node. Second, random simulation computes a signature for every node; nodes with different signatures are certainly different. Third, SAT sweeping proves candidate internal equivalences bottom-up and merges each proved pair, often with cut points, so that later queries are small. Counterexamples to internal candidates are not design bugs; they are added to the simulation patterns to split other wrong candidates. Only after this does the engine solve the output miters, which are by then usually trivial. Synthesized netlists have many internal equivalences with their source because synthesis mostly makes local changes, which is what makes this pipeline work.

## Engines and results

BDDs and SAT both decide CEC. BDDs are canonical, so once both outputs are built, comparison is free, but building them fails on large datapaths and multipliers. SAT with sweeping scales to large netlists and is the default today. In ABC the command cec runs the full simulation, sweeping and SAT pipeline; Yosys offers equivalence passes (the equiv family) that check designs it has transformed itself. A result must be read carefully. UNSAT on every miter is a proof. A counterexample may still reflect an under-specified environment: if it drives an input combination the system never applies, the designs may be equivalent under a missing constraint, and the fix is an assumption, not an RTL change. A timeout on one cone is not a proof, though the other cones remain proved.

## Going deeper

Kuehlmann et al., "Robust Boolean Reasoning for Equivalence Checking and Functional Property Verification" (IEEE TCAD, 2002); Mishchenko et al., "Improvements to Combinational Equivalence Checking" (ICCAD 2006).
