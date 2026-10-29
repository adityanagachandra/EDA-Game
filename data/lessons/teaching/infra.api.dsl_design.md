# Domain-specific languages for EDA

EDA is full of small languages: SDC for constraints, Liberty for cell libraries, UPF for power intent, rule decks for DRC and LVS, synthesis scripts, and flow configuration files. A tool designer must decide, for each new kind of input, whether it is an embedded language hosted in Tcl or Python, or an external declarative format with its own parser, and how it will be validated. Those choices decide what users can express and what tools can analyze.

## Embedded and external DSLs

An embedded DSL adds domain commands to a host language. SDC is the classic case: create_clock and set_input_delay are commands run by a Tcl interpreter inside the timer, so SDC inherits Tcl variables, loops and procedures. UPF and many flow scripts follow the same pattern. An external DSL has its own grammar and parser. Liberty is one: groups and attributes describing cells, with no variables or loops. Declarative configuration in YAML or JSON is another. Embedding is cheap to build and powerful to use, because the host language provides parsing, expressions and control flow for free. External declarative formats are easier to analyze: a tool can validate them against a schema, compare two versions, or generate them from a program without running any code. Yosys scripts sit between, as a command language that names passes in a registry and orders them.

## Validation and semantic checks

A DSL is only as good as its error reporting. Syntax errors are found by the parser, but most real mistakes are semantic: a clock with a negative period, a delay on a port that does not exist, a misspelled configuration key. Semantic checks belong where the context is known, when the command executes or the key is read, and they should name the file and line. Silent defaults are the worst failure: a flow configuration with PLACE_DENSTY misspelled that quietly uses the default density produces wrong results with no error. Declaring a schema, the known keys with their types and ranges, and rejecting unknown keys turns such typos into immediate failures. The same discipline makes a DSL safe to accept from automated sources: an agent proposing flow parameters through a validated schema can only touch the intended knobs.

## Analyzing and trusting DSL programs

The power of an embedded DSL limits what can be known without running it. An SDC file that uses a Tcl loop to generate 10,000 false-path commands has those commands only after evaluation, so a static checker must execute the Tcl, or a restricted subset of it, to see them. A declarative rule language, by contrast, exposes its structure, letting an engine reorder rules, share work among them and run them in parallel. Trust is the other issue. A full Tcl interpreter can execute shell commands, open files and sockets; when input may come from an agent or another team, the tool needs a safe interpreter or a restricted command set. Finally, the meaning of a configuration key can change between tool versions. Versioning the schema and recording it with every run lets a regression explain why the same input gave different results.

## Going deeper

Martin Fowler, "Domain-Specific Languages" (2010); the Synopsys Design Constraints format documentation; the Liberty reference manual's description of its syntax.
