# Engines as callable tools for agents

For forty years EDA tools were built for a human at a console: a GUI, a Tcl prompt, a log to read. An agent-first tool is built for a program that calls it. The placer, router, timer and extractor become engines with typed entry points and structured results, and the human moves from driving every step to reviewing what the agents did. For a tool designer, this changes what an engine must expose, what it must return and where its state may live.

## From sessions to calls

A traditional flow is a session: the user loads a design, sets variables, runs commands and reads the console, and the tool keeps a large amount of implicit state between commands. An agent cannot reliably drive that. It needs entry points with explicit inputs (a design handle, a constraint set, a parameter record, a seed) and explicit outputs, so the same call with the same inputs gives the same answer. The scripting layers that already exist, Tcl in OpenROAD and OpenSTA, Python bindings to OpenDB, are the starting point, but they are usually thin wrappers over a stateful session. The agent-first version makes each operation a function of its arguments: run global placement on this design state with these weights, return the new state and its metrics. The GUI does not disappear; it becomes a viewer of the same design database.

## What a call must return

An engine called by an agent must return results a program can act on without parsing prose. A placement call returns HPWL, overflow, density and runtime as named numbers, a status (converged, diverged, timed out), the identifiers of the design state it produced, and provenance: engine version, parameters and seed. Failures are structured too: which check failed, on which object, with what values. Compare a timer that prints a 3,000-line report with one that returns the worst 20 endpoints as records of endpoint, slack, path group and the arcs on the path. The second lets the agent decide its next step in milliseconds; the first forces it to scrape text, which breaks whenever the report format changes.

## Where the compute stays

Making engines callable does not move the algorithms into the model. The language model makes cheap decisions: which engine to call, with which skill and parameters, and whether the result is good enough. The compiled engine does the expensive, exact work: solving the placement, propagating arrival times, legalizing cells. An agent that tried to compute slack itself, or to legalize a placement by reasoning in text, would be slow and wrong. The design principle is that every number the agent relies on comes from an engine, and the agent's job is orchestration and judgment. That is also why engine speed still matters as much as ever: an agent loop calls the engine hundreds of times.

## Going deeper

The OpenROAD documentation on its Tcl and Python interfaces; the OpenSTA command reference.
