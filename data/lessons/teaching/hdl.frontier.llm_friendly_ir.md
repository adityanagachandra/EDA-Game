# LLM-friendly hardware IRs

An agent that edits hardware can work on Verilog text, on a netlist IR or on a structured compiler IR. The choice matters: some representations make edits local, explicit and checkable, while others hide meaning in implicit rules and distant context. This concept asks what makes an IR friendly to language models, and how tools should expose it to them.

## Properties that make edits safe

Four properties make an IR easier for a model to edit correctly. A textual form that round-trips exactly, so the model reads and writes the same thing the compiler holds. Explicit types and widths on every value, so an edit can be checked without reconstructing context from declarations far away. SSA form, so each value has one definition and the effect of a change can be traced through explicit uses. And a strict parser and verifier that report the offending operation and the expected type, turning a malformed edit into immediate, actionable feedback.

Stable, source-derived names complete the list. They let the agent relate an IR value to the RTL signal it came from and to the timing or lint report that mentions it, which is what makes acting on a report possible.

## Comparing hardware representations

Verilog is what humans write, but its implicit rules (expression sizing, signedness, blocking order, preprocessor macros) are exactly where edits silently go wrong. A structural IR such as CIRCT's hw and comb dialects has none of them: each operation has one meaning, widths are explicit and signless integers remove signedness ambiguity. MLIR adds declarative operation definitions with verifiers and a pattern-rewriting framework, giving a structured way to express transformations.

Lower levels serve algorithms, not editors. An and-inverter graph is ideal for logic optimization but is bit-level and anonymous, so the designer's intent and word structure are gone. RTLIL sits between: word-level cells with source attributes, readable and round-trippable, close to what synthesis sees. For a 200,000-line design, the representation matters less than retrieval: a query tool that returns the relevant operations and their neighborhood beats pasting the whole design into a prompt.

## Exposing IR edits as tools

Free-text patches to IR can be malformed in ways the model cannot see. A typed tool interface, for example over MCP, offers operations such as replace operand, insert operation or rename value, validated by the tool before they touch the IR. This is a seam: the model chooses edits, the tool guarantees their well-formedness.

Well-formed is not correct. An edit that passes the verifier proves only structural validity, not that the design still computes the intended function, so equivalence checks or tests still follow. The easiest proposals to verify are small: a local rewrite pattern with a stated precondition can be proven by an equivalence tool on small instances or tested exhaustively, while a rewritten whole design cannot be checked as cheaply. Explicit dataflow, where every dependence is a visible operand rather than shared mutable state, keeps an edit's reach computable from its uses.

## Going deeper

Lattner et al., "MLIR: Scaling Compiler Infrastructure for Domain Specific Computation" (CGO 2021); the CIRCT rationale documents; the Model Context Protocol specification.
